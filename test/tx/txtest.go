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

package main

import (
	//"bytes"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	//"sync"

	jww "github.com/spf13/jwalterweatherman"
	"github.com/fractalplatform/fractal/asset"
	"github.com/fractalplatform/fractal/common"
	"github.com/fractalplatform/fractal/crypto"
	"github.com/fractalplatform/fractal/params"
	tc "github.com/fractalplatform/fractal/test/common"
	"github.com/fractalplatform/fractal/types"
	"github.com/fractalplatform/fractal/utils/rlp"
	"github.com/fractalplatform/fractal/accountmanager"
	"time"
)

var (
	minerprikey, _   = crypto.HexToECDSA("289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032")
	minerpubkey      = common.HexToPubKey("0x047db227d7094ce215c3a0f57e1bcc732551fe351f94249471934567e0f5dc1bf795962b8cccb87a2eb56b29fbe37d614e2f4c3c45b789ae4f1f51f4cb21972ffd")
	newPrivateKey, _ = crypto.HexToECDSA("8ee847ae5974a13ce9df66083e453ea1e0f7995379ed027a98e827aa8b6bc211")
	gaslimit         = uint64(2000000)
	minername        = common.Name("ftsystemio")
	toname           = common.Name("testtest11")
	issueAmount      = new(big.Int).Mul(big.NewInt(10), big.NewInt(1e18))
	inCreateAmount   = big.NewInt(100000000)
	indexstr         = "abcdefghijklmnopqrstuvwxyz0123456789"
	basefrom         = "newnamefrom%s"
	baseto           = "newnameto%s"
	testbase         = "testtest"
	testname1        = ""
)

type GenAction struct {
	*types.Action
	PrivateKey *ecdsa.PrivateKey
}

func init() {
	jww.SetLogThreshold(jww.LevelTrace)
	jww.SetStdoutThreshold(jww.LevelInfo)
}

func GeneragePubKey() (common.PubKey, *ecdsa.PrivateKey) {
	prikey, _ := crypto.GenerateKey()
	return common.BytesToPubKey(crypto.FromECDSAPub(&prikey.PublicKey)), prikey
}

func createAccount(accountName common.Name, founder common.Name, from, newname common.Name, nonce uint64, publickey common.PubKey, prikey *ecdsa.PrivateKey) {
	account := &accountmanager.AccountAction{
		AccountName: accountName,
		Founder:     founder,
		ChargeRatio: 80,
		PublicKey:   publickey,
	}
	payload, err := rlp.EncodeToBytes(account)
	if err != nil {
		panic("rlp payload err")
	}
	to := newname
	gc := newGeAction(types.CreateAccount, from, to, nonce, 1, gaslimit, nil, payload, prikey)
	var gcs []*GenAction
	gcs = append(gcs, gc)
	sendTxTest(gcs)
}

func updateAccount(from, founder common.Name, newpubkey common.PubKey, nonce uint64, privatekey *ecdsa.PrivateKey) {
	account := &accountmanager.AccountAction{
		AccountName: from,
		Founder:     founder,
		ChargeRatio: 80,
		PublicKey:   newpubkey,
	}
	payload, err := rlp.EncodeToBytes(account)
	if err != nil {
		panic("rlp payload err")
	}
	gc := newGeAction(types.UpdateAccount, from, "", nonce, 1, gaslimit, nil, payload, privatekey)
	var gcs []*GenAction
	gcs = append(gcs, gc)
	sendTxTest(gcs)
}

func issueAsset(atype types.ActionType, assetid uint64, from, Owner, founder common.Name, amount *big.Int, assetname string, nonce uint64, prikey *ecdsa.PrivateKey) {
	ast := &asset.AssetObject{
		AssetId:    assetid,
		AssetName:  assetname,
		Symbol:     fmt.Sprintf("Symbol%d", nonce),
		Amount:     amount,
		Decimals:   2,
		Founder:    founder,
		AddIssue:   nil,
		Owner:      Owner,
		UpperLimit: big.NewInt(100000000000000000),
	}
	payload, err := rlp.EncodeToBytes(ast)
	if err != nil {
		panic("rlp payload err")
	}
	gc := newGeAction(atype, from, "", nonce, 1, gaslimit, nil, payload, prikey)
	var gcs []*GenAction
	gcs = append(gcs, gc)
	sendTxTest(gcs)
}

func increaseAsset(from, to common.Name, assetid uint64, nonce uint64, prikey *ecdsa.PrivateKey) {
	ast := &accountmanager.IncAsset{
		AssetId: assetid,
		To:      to,
		Amount:  inCreateAmount,
	}
	payload, err := rlp.EncodeToBytes(ast)
	if err != nil {
		panic("rlp payload err")
	}
	gc := newGeAction(types.IncreaseAsset, from, "", nonce, 1, gaslimit, nil, payload, prikey)
	var gcs []*GenAction
	gcs = append(gcs, gc)
	sendTxTest(gcs)
}

func transfer(from, to common.Name, amount *big.Int, nonce uint64, prikey *ecdsa.PrivateKey) {
	gc := newGeAction(types.Transfer, from, to, nonce, 1, gaslimit, amount, nil, prikey)
	var gcs []*GenAction
	gcs = append(gcs, gc)
	sendTxTest(gcs)
}

func newGeAction(at types.ActionType, from, to common.Name, nonce uint64, assetid uint64, gaslimit uint64, amount *big.Int, payload []byte, prikey *ecdsa.PrivateKey) *GenAction {
	action := types.NewAction(at, from, to, nonce, assetid, gaslimit, amount, payload)

	return &GenAction{
		Action:     action,
		PrivateKey: prikey,
	}
}

func sendTxTest(gcs []*GenAction) {
	signer := types.NewSigner(params.DefaultChainconfig.ChainID)
	var actions []*types.Action
	for _, v := range gcs {
		actions = append(actions, v.Action)
	}
	tx := types.NewTransaction(uint64(1), big.NewInt(1), actions...)
	for _, v := range gcs {
		err := types.SignAction(v.Action, tx, signer, v.PrivateKey)
		if err != nil {
			panic(fmt.Sprintf("SignAction err %v", err))
		}

	}
	rawtx, err := rlp.EncodeToBytes(tx)
	if err != nil {
		jww.ERROR.Fatalln(err)
	}
	hash, err := tc.SendRawTx(rawtx)
	if err != nil {
		panic(err)
	}
	jww.INFO.Printf("hash: %x", hash)
}

var (
	pub1 = "0x0468cba7890aae10f3dde57d269cf7c4ba14cc0efc2afee86791b0a22b794820febdb2e5c6c56878a308e7f62ad2d75739de40313a72975c993dd76a5301a03d12"
	pri1 = "357a2cbdd91686dcbe2c612e9bed85d4415f62446440839466bf7b2f1ab135b7"

	pub2 = "0x04fa0b2a9b2d0542bf2912c4c6500ba64a26652e302370ed5645b1c32df50fbe7a5f12da0b278638e1df6753a7c6ac09e68cb748cfe6d45102114f52e95e9ed652"
	pri2 = "340cde826336f1adb8673ec945819d073af00cffb5c174542e35ff346445e213"

	pubkey1    = common.HexToPubKey(pub1)
	prikey1, _ = crypto.HexToECDSA(pri1)

	pubkey2    = common.HexToPubKey(pub2)
	prikey2, _ = crypto.HexToECDSA(pri2)
)

func main() {
	nonce, _ := tc.GetNonce(minername)
	//pub, pri := GeneragePubKey()
	createAccount(toname, "", minername, toname, nonce, pubkey1, minerprikey)
	nonce++
	transfer(minername, toname, issueAmount, nonce, minerprikey)
	nonce++

	toname1 := common.Name("testtest12")
	//pub1, _ := GeneragePubKey()
	createAccount(toname1, "", minername, toname1, nonce, pubkey2, minerprikey)
	nonce++

	transfer(minername, toname1, issueAmount, nonce, minerprikey)
	nonce++

	time.Sleep(time.Duration(3) * time.Second)

	t1nonce, _ := tc.GetNonce(toname)
	//newpub2, _ := GeneragePubKey()
	updateAccount(toname, toname1, pubkey2, t1nonce, prikey1)

	time.Sleep(time.Duration(3) * time.Second)

	t2nonce, _ := tc.GetNonce(toname1)
	issueAsset(types.IssueAsset, 0, toname1, toname1, toname1, big.NewInt(10000000000000), "testnewasset", t2nonce, prikey2)

	time.Sleep(time.Duration(3) * time.Second)

	t2nonce ++
	increaseAsset(toname1, toname1, 2, t2nonce, prikey2)

	time.Sleep(time.Duration(3) * time.Second)

	t2nonce ++
	issueAsset(types.DestroyAsset, 2, toname1, "", "", big.NewInt(100000), "testnewasset", t2nonce, prikey2)

	time.Sleep(time.Duration(3) * time.Second)

	t2nonce++
	issueAsset(types.SetAssetOwner, 2, toname1, toname, "", big.NewInt(100000), "testnewasset", t2nonce, prikey2)
}
