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
	"time"

	//"os"
	"strconv"
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
)

var (
	minerprikey, _   = crypto.HexToECDSA("289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032")
	minerpubkey      = common.HexToPubKey("0x047db227d7094ce215c3a0f57e1bcc732551fe351f94249471934567e0f5dc1bf795962b8cccb87a2eb56b29fbe37d614e2f4c3c45b789ae4f1f51f4cb21972ffd")
	newPrivateKey, _ = crypto.HexToECDSA("8ee847ae5974a13ce9df66083e453ea1e0f7995379ed027a98e827aa8b6bc211")
	gaslimit         = uint64(2000000)
	minername        = common.Name("ftsystemio")
	toname           = common.Name("testtest14")
	issueAmount      = new(big.Int).Mul(big.NewInt(1000000), big.NewInt(1e18))
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

func createAccount(accountName common.Name, founder common.Name, from, newname common.Name, nonce uint64, prikey *ecdsa.PrivateKey, pubkey common.PubKey) {
	account := &accountmanager.AccountAction{
		AccountName: accountName,
		Founder:     founder,
		ChargeRatio: 80,
		PublicKey:   pubkey,
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

func updateAccount(accountName common.Name, founder common.Name, from, newname common.Name, nonce uint64, prikey *ecdsa.PrivateKey, pubkey common.PubKey) {
	account := &accountmanager.AccountAction{
		AccountName: accountName,
		Founder:     founder,
		ChargeRatio: 80,
		PublicKey:   pubkey,
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

func deleteAccount(from common.Name, nonce uint64, prikey *ecdsa.PrivateKey) {
	gc := newGeAction(types.DeleteAccount, from, "", nonce, 1, gaslimit, nil, nil, prikey)
	var gcs []*GenAction
	gcs = append(gcs, gc)
	sendTxTest(gcs)
}


func issueAsset(from, Owner common.Name, amount *big.Int, assetname string, nonce uint64, prikey *ecdsa.PrivateKey) {
	ast := &asset.AssetObject{
		AssetName: assetname,
		Symbol:    fmt.Sprintf("Symbol%d", nonce),
		Amount:    amount,
		Decimals:  2,
		Owner:     Owner,
	}
	payload, err := rlp.EncodeToBytes(ast)
	if err != nil {
		panic("rlp payload err")
	}
	gc := newGeAction(types.IssueAsset, from, "", nonce, 1, gaslimit, nil, payload, prikey)
	var gcs []*GenAction
	gcs = append(gcs, gc)
	sendTxTest(gcs)
}

func increaseAsset(from common.Name, assetid uint64, assetname string, nonce uint64, prikey *ecdsa.PrivateKey) {
	ast := &accountmanager.IncAsset{
		AssetId:   assetid,
		//AssetName: assetname,
		//Symbol:    fmt.Sprintf("Symbol%d", nonce),
		Amount:    inCreateAmount,
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

func setAssetOwner(from, newowner common.Name, assetid uint64, nonce uint64, prikey *ecdsa.PrivateKey) {
	ast := &asset.AssetObject{
		AssetId: assetid,
		Owner:   newowner,
	}
	payload, err := rlp.EncodeToBytes(ast)
	if err != nil {
		panic("rlp payload err")
	}
	gc := newGeAction(types.SetAssetOwner, from, "", nonce, 1, gaslimit, nil, payload, prikey)
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
	//nonce := GetNonce(sendaddr, "latest")
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
	//pubkey, err := types.Recover(signer, tx.GetActions()[0], tx)
	// if err != nil {
	// 	jww.ERROR.Fatalln(err)
	// }
	//fmt.Println("===>",pubkey)

	rawtx, err := rlp.EncodeToBytes(tx)
	//fmt.Sprintf("rawtx:", rawtx)
	if err != nil {
		jww.ERROR.Fatalln(err)
	}
	//jww.INFO.Println(rawtx)
	_, _ = tc.SendRawTx(rawtx)
}

func main() {
	nonce, _ := tc.GetNonce(minername)
	for i := 1; i <= 4; i++ {
		itost := strconv.FormatUint(uint64(i), 10)
		testname1 = testbase + itost
		toname = common.Name(testname1)
		jww.INFO.Println(i)
		//jww.INFO.Println("nonce:",nonce)

		createAccount(toname, "", minername, toname, nonce, minerprikey, minerpubkey)
		nonce++
		transfer(minername, toname, issueAmount, nonce, minerprikey)
		time.Sleep(10 * time.Millisecond)
		//time.Sleep(time.Duration(timecount)*time.Millisecond)
		//pubkey, _ := GeneragePubKey()
		//updateAccount(toname, uint64(2), pubkey, minerprikey)
		//createAccount(minername,toname,nonce,minerprikey,)
	}
	nonce++
	//updateAccount()
	//issueAsset()
	//extraData := "73756f7918dd840ba20c58c907fc7c1e4cbc5fc87ebe8ce76b273175057da0fe88854faf5cafec14091c67ab327b2fb60a561b90b7e1b454cdab0ee48af32e8d95136c7d01"
	//str,err := hex.DecodeString(extraData)
	//jww.INFO.Println("extraData:", string(str[0:len(str)-65]), len(str)-65, err)
}
func main1() {
	nonce, _ := tc.GetNonce(minername)
	//jww.INFO.Println(hexutil.Encode(crypto.FromECDSAPub(&minerprikey.PublicKey)))
	//toname = common.Name("")
	transfer(minername, toname, issueAmount, nonce, minerprikey)
	//issueAsset(minername,minername, toname, issueAmount, issueAmount,issueAmount,"as4", nonce, minerprikey)
	//updateAccount(minername, nonce, minerpubkey, minerprikey)

}

func test() {
	rs, _ := tc.GetAccountByName("testtest1")
	for i := 0; i < len(rs.Balances); i++ {
		fmt.Println(rs.Balances[i].AssetID, rs.Balances[i].Balance.Uint64())
	}
}
