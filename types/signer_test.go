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

package types

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/fractalplatform/fractal/common"
	"github.com/fractalplatform/fractal/crypto"
)

func TestSigning(t *testing.T) {
	key, _ := crypto.GenerateKey()
	exp := crypto.FromECDSAPub(&key.PublicKey)
	signer := NewSigner(big.NewInt(18))
	if err := SignAction(testTx.GetActions()[0], testTx, signer, key); err != nil {
		t.Fatal(err)
	}

	pubkey, err := Recover(signer, testTx.GetActions()[0], testTx)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(pubkey.Bytes(), exp) != 0 {
		t.Errorf("exected from and address to be equal. Got %x want %x", pubkey, exp)
	}
}

func TestSigningMultiKey(t *testing.T) {
	keyPair := make([]*KeyPair, 0)
	pubs := make([]common.PubKey, 0)
	for i := 0; i < 4; i++ {
		key, _ := crypto.GenerateKey()
		exp := crypto.FromECDSAPub(&key.PublicKey)
		keyPair = append(keyPair, &KeyPair{priv: key, index: []uint64{uint64(i)}})
		pubs = append(pubs, common.BytesToPubKey(exp))
	}

	signer := NewSigner(big.NewInt(18))
	if err := SignActionWithMultiKey(testTx.GetActions()[0], testTx, signer, keyPair); err != nil {
		t.Fatal(err)
	}

	pubkeys, err := RecoverMultiKey(signer, testTx.GetActions()[0], testTx)
	if err != nil {
		t.Fatal(err)
	}

	for i, pubkey := range pubkeys {
		if pubkey.Compare(pubs[i]) != 0 {
			t.Errorf("exected from and address to be equal. Got %x want %x", pubkey, pubs[i])
		}
	}
}

func TestChainID(t *testing.T) {
	key, _ := crypto.GenerateKey()

	signer := NewSigner(big.NewInt(1))
	if err := SignAction(testAction, testTx, signer, key); err != nil {
		t.Fatal(err)
	}

	if testTx.GetActions()[0].ChainID().Cmp(signer.chainID) != 0 {
		t.Error("expected chainId to be", signer.chainID, "got", testTx.GetActions()[0].ChainID())
	}
}
