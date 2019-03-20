// Copyright 2018 The Fractal Team Authors
// This file is part of the fractal project.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package txpool

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"testing"
	"time"

	am "github.com/fractalplatform/fractal/accountmanager"
	"github.com/fractalplatform/fractal/asset"

	"github.com/fractalplatform/fractal/common"
	"github.com/fractalplatform/fractal/crypto"
	"github.com/fractalplatform/fractal/event"
	"github.com/fractalplatform/fractal/params"
	"github.com/fractalplatform/fractal/state"
	"github.com/fractalplatform/fractal/types"
	"github.com/fractalplatform/fractal/utils/fdb"
)

// testTxPoolConfig is a transaction pool configuration without stateful disk
// sideeffects used during testing.
var testTxPoolConfig Config

func init() {
	testTxPoolConfig = Config{
		Journal:      "",
		Rejournal:    time.Hour,
		PriceLimit:   1,
		PriceBump:    10,
		AccountSlots: 16,
		GlobalSlots:  4096,
		AccountQueue: 64,
		GlobalQueue:  1024,
		Lifetime:     3 * time.Hour,
		GasAssetID:   uint64(1),
	}
}

type testBlockChain struct {
	statedb       *state.StateDB
	gasLimit      uint64
	chainHeadFeed *event.Feed
}

func (bc *testBlockChain) CurrentBlock() *types.Block {
	return types.NewBlock(&types.Header{
		GasLimit: bc.gasLimit,
	}, nil, nil)
}

func (bc *testBlockChain) GetBlock(hash common.Hash, number uint64) *types.Block {
	return bc.CurrentBlock()
}

func (bc *testBlockChain) StateAt(common.Hash) (*state.StateDB, error) {
	return bc.statedb, nil
}

func (bc *testBlockChain) GetForkID(statedb *state.StateDB) (uint64, uint64, error) {
	return 0, 0, nil
}

func (bc *testBlockChain) State() (*state.StateDB, error) {
	return bc.statedb, nil
}

func transaction(nonce uint64, from, to common.Name, gaslimit uint64, key *ecdsa.PrivateKey) *types.Transaction {
	return pricedTransaction(nonce, from, to, gaslimit, big.NewInt(1), key)
}

func pricedTransaction(nonce uint64, from, to common.Name, gaslimit uint64, gasprice *big.Int, key *ecdsa.PrivateKey) *types.Transaction {
	tx := newTx(gasprice, newAction(nonce, from, to, big.NewInt(100), gaslimit, nil))
	if err := types.SignAction(tx.GetActions()[0], tx, types.NewSigner(params.DefaultChainconfig.ChainID), key); err != nil {
		panic(err)
	}
	return tx
}

func generateAccount(t *testing.T, name common.Name, managers ...*am.AccountManager) *ecdsa.PrivateKey {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	pubkeyBytes := crypto.FromECDSAPub(&key.PublicKey)
	for _, m := range managers {
		if err := m.CreateAccount(name, common.Name(""), 80, common.BytesToPubKey(pubkeyBytes)); err != nil {
			t.Fatal(err)
		}
	}
	return key
}

func setupTxPool(assetOwner common.Name) (*TxPool, *am.AccountManager) {

	statedb, _ := state.New(common.Hash{}, state.NewDatabase(fdb.NewMemDatabase()))
	asset := asset.NewAsset(statedb)
	asset.IssueAsset("ft", "zz", new(big.Int).SetUint64(params.Fractal), 10, assetOwner, assetOwner, big.NewInt(1000000))
	blockchain := &testBlockChain{statedb, 1000000, new(event.Feed)}
	manager, _ := am.NewAccountManager(statedb)
	return New(testTxPoolConfig, params.DefaultChainconfig, blockchain), manager
}

// validateTxPoolInternals checks various consistency invariants within the pool.
func validateTxPoolInternals(pool *TxPool) error {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	// Ensure the total transaction set is consistent with pending + queued
	pending, queued := pool.stats()
	if total := pool.all.Count(); total != pending+queued {
		return fmt.Errorf("total transaction count %d != %d pending + %d queued", total, pending, queued)
	}
	if priced := pool.priced.items.Len() - pool.priced.stales; priced != pending+queued {
		return fmt.Errorf("total priced transaction count %d != %d pending + %d queued", priced, pending, queued)
	}
	// Ensure the next nonce to assign is the correct one

	for addr, list := range pool.pending {
		// Find the last transaction
		var last uint64
		for nonce := range list.txs.items {
			if last < nonce {
				last = nonce
			}
		}

		nonce, err := pool.pendingAccountManager.GetNonce(addr)
		if err != nil {
			return err
		}
		if nonce != last+1 {
			return fmt.Errorf("pending nonce mismatch: have %v, want %v", nonce, last+1)
		}
	}
	return nil
}

// validateEvents checks that the correct number of transaction addition events
// were fired on the pool's event event.
func validateEvents(events chan *event.Event, count int) error {
	var received []*types.Transaction
	for len(received) < count {
		select {
		case ev := <-events:
			if ev.Typecode == event.TxEv {
				received = append(received, ev.Data.([]*types.Transaction)...)
			}
		case <-time.After(time.Second):
			return fmt.Errorf("event #%v not fired", received)
		}
	}

	if len(received) > count {
		return fmt.Errorf("more than %d events fired1: %v", count, received[count:])
	}
	select {
	case ev := <-events:
		return fmt.Errorf("more than %d events fired2: %v", count, ev.Typecode)
	case <-time.After(50 * time.Millisecond):
		// This branch should be "default", but it's a data race between goroutines,
		// reading the event channel and pushing into it, so better wait a bit ensuring
		// really nothing gets injected.
	}
	return nil
}

func deriveSender(tx *types.Transaction) ([]common.PubKey, error) {
	var pubkeys []common.PubKey
	for _, a := range tx.GetActions() {
		pubkey, err := types.Recover(types.NewSigner(params.DefaultChainconfig.ChainID), a, tx)
		if err != nil {
			return nil, err
		}
		pubkeys = append(pubkeys, pubkey)
	}
	return pubkeys, nil
}

type testChain struct {
	*testBlockChain

	name    common.Name
	trigger *bool
}

// testChain.State() is used multiple times to reset the pending state.
// when simulate is true it will create a state that indicates
// that tx0 and tx1 are included in the chain.
func (c *testChain) State() (*state.StateDB, error) {
	// delay "state change" by one. The tx pool fetches the
	// state multiple times and by delaying it a bit we simulate
	// a state change between those fetches.
	stdb := c.statedb
	if *c.trigger {
		c.statedb, _ = state.New(common.Hash{}, state.NewDatabase(fdb.NewMemDatabase()))
		am, err := am.NewAccountManager(c.statedb)
		if err != nil {
			return nil, err
		}

		// simulate that the new head block included tx0 and tx1
		if err := am.SetNonce(c.name, 2); err != nil {
			return nil, err
		}
		if err := am.AddAccountBalanceByID(c.name, uint64(0), new(big.Int).SetUint64(params.Fractal)); err != nil {
			return nil, err
		}
		*c.trigger = false
	}
	return stdb, nil
}

func newAction(nonce uint64, from, to common.Name, amount *big.Int, gasLimit uint64, data []byte) *types.Action {
	return types.NewAction(types.Transfer, from, to, nonce, uint64(1), gasLimit, amount, data)
}

func newTx(gasPrice *big.Int, action ...*types.Action) *types.Transaction {
	return types.NewTransaction(uint64(1), gasPrice, action...)
}
