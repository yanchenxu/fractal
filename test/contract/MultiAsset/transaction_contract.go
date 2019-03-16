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
	"bytes"
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"
	"time"
	"crypto/ecdsa"

	jww "github.com/spf13/jwalterweatherman"

	"github.com/fractalplatform/fractal/common"
	"github.com/fractalplatform/fractal/crypto"
	"github.com/fractalplatform/fractal/types"
	"github.com/fractalplatform/fractal/utils/abi"
	"github.com/fractalplatform/fractal/utils/rlp"
	tc "github.com/fractalplatform/fractal/test/common"

)

var (
	abifile         = "./MultiAsset.abi"
	binfile         = "./MultiAsset.bin"
	systemprikey, _ = crypto.HexToECDSA("289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032")
	systemname      = common.Name("ftsystemio")
	to              = common.Name("testtest11")
	newFrom         = common.Name("testtest11")

	contractAddr = common.Name("testtest11")
	assetID      = uint64(1)
	gasLimit = uint64(2000000)

	pub1 = "0x0468cba7890aae10f3dde57d269cf7c4ba14cc0efc2afee86791b0a22b794820febdb2e5c6c56878a308e7f62ad2d75739de40313a72975c993dd76a5301a03d12"
	pri1 = "357a2cbdd91686dcbe2c612e9bed85d4415f62446440839466bf7b2f1ab135b7"

	pub2 = "0x04fa0b2a9b2d0542bf2912c4c6500ba64a26652e302370ed5645b1c32df50fbe7a5f12da0b278638e1df6753a7c6ac09e68cb748cfe6d45102114f52e95e9ed652"
	pri2 = "340cde826336f1adb8673ec945819d073af00cffb5c174542e35ff346445e213"

	pubkey1    = common.HexToPubKey(pub1)
	prikey1, _ = crypto.HexToECDSA(pri1)

	pubkey2    = common.HexToPubKey(pub2)
	prikey2, _ = crypto.HexToECDSA(pri2)
)

func hexToBigInt(hex string) *big.Int {
	n := new(big.Int)
	n, _ = n.SetString(hex[2:], 16)

	return n
}

func input(abifile string, method string, params ...interface{}) (string, error) {
	var abicode string

	hexcode, err := ioutil.ReadFile(abifile)
	if err != nil {
		fmt.Printf("Could not load code from file: %v\n", err)
		return "", err
	}
	abicode = string(bytes.TrimRight(hexcode, "\n"))

	parsed, err := abi.JSON(strings.NewReader(abicode))
	if err != nil {
		fmt.Println("abi.json error ", err)
		return "", err
	}

	input, err := parsed.Pack(method, params...)
	if err != nil {
		fmt.Println("parsed.pack error ", err)
		return "", err
	}
	return common.Bytes2Hex(input), nil
}

func formCreateContractInput(abifile string, binfile string) ([]byte, error) {
	hexcode, err := ioutil.ReadFile(binfile)
	if err != nil {
		jww.INFO.Printf("Could not load code from file: %v\n", err)
		return nil, err
	}
	code := common.Hex2Bytes(string(bytes.TrimRight(hexcode, "\n")))

	createInput, err := input(abifile, "")
	if err != nil {
		jww.INFO.Println("createInput error ", err)
		return nil, err
	}

	createCode := append(code, common.Hex2Bytes(createInput)...)
	return createCode, nil
}

func formIssueAssetInput(abifile string, desc string) ([]byte, error) {
	issueAssetInput, err := input(abifile, "reg", desc)
	if err != nil {
		jww.INFO.Println("createInput error ", err)
		return nil, err
	}
	return common.Hex2Bytes(issueAssetInput), nil
}
func formIssueAssetInput1(abifile string, assetId *big.Int, to common.Address, value *big.Int) ([]byte, error) {
	issueAssetInput, err := input(abifile, "add", assetId, to, value)
	if err != nil {
		jww.INFO.Println("createInput error ", err)
		return nil, err
	}
	return common.Hex2Bytes(issueAssetInput), nil
}
func formSetAssetOwner(abifile string, newOwner common.Address, assetId *big.Int) ([]byte, error) {
	issueAssetInput, err := input(abifile, "changeOwner", newOwner, assetId)
	if err != nil {
		jww.INFO.Println("createInput error ", err)
		return nil, err
	}
	return common.Hex2Bytes(issueAssetInput), nil
}

func formTransferAssetInput(abifile string, toAddr common.Address, assetId *big.Int, value *big.Int) ([]byte, error) {
	transferAssetInput, err := input(abifile, "transAsset", toAddr, assetId, value)
	if err != nil {
		jww.INFO.Println("transferAssetInput error ", err)
		return nil, err
	}
	return common.Hex2Bytes(transferAssetInput), nil
}

func init() {
	jww.SetLogThreshold(jww.LevelTrace)
	jww.SetStdoutThreshold(jww.LevelInfo)

}

func sendDeployContractTransaction() {
	jww.INFO.Println("test sendDeployContractTransaction... ")
	input, err := formCreateContractInput(abifile, binfile)
	if err != nil {
		jww.INFO.Println("sendDeployContractTransaction formCreateContractInput error ... ", err)
		return
	}
	nonce, _ := tc.GetNonce(systemname)
	sendTransferTx(types.CreateContract, systemname, contractAddr, nonce, assetID, big.NewInt(100000000000), input, systemprikey)
}

func sendIssueTransaction() {
	jww.INFO.Println("test sendIssueTransaction... ")
	issueStr := "eth4" + systemname.String() + ",ft,10000000000,10," + systemname.String() //"ft" + contractAddr.String() + ",ft,10000000000,10," + contractAddr.String() + ",9000000000000000," + contractAddr.String() //25560
	input, err := formIssueAssetInput(abifile, issueStr)
	if err != nil {
		jww.INFO.Println("sendIssueTransaction formIssueAssetInput error ... ", err)
		return
	}
	nonce, _ := tc.GetNonce(systemname)
	sendTransferTx(types.CallContract, systemname, contractAddr, nonce, assetID, big.NewInt(0), input, systemprikey)
}

func sendIncreaseIssueTransaction() {
	jww.INFO.Println("test sendIssueTransaction... ")
	input, err := formIssueAssetInput1(abifile, big.NewInt(3), common.BytesToAddress([]byte("testtest1")), big.NewInt(10)) //21976   21848

	if err != nil {
		jww.INFO.Println("sendIssueTransaction formIssueAssetInput error ... ", err)
		return
	}
	nonce, _ := tc.GetNonce(contractAddr)
	sendTransferTx(types.CallContract, contractAddr, contractAddr, nonce, assetID, big.NewInt(0), input, prikey1)
}

func sendSetOwnerIssueTransaction() {
	jww.INFO.Println("test sendIssueTransaction... ")
	input, err := formSetAssetOwner(abifile, common.BytesToAddress([]byte("testtest2")), big.NewInt(3)) //22168

	if err != nil {
		jww.INFO.Println("sendIssueTransaction formIssueAssetInput error ... ", err)
		return
	}

	nonce, _ := tc.GetNonce(contractAddr)
	sendTransferTx(types.CallContract, contractAddr, contractAddr, nonce, assetID, big.NewInt(0), input, prikey1)
}

func sendTransferToContractTransaction() {
	nonce, _ := tc.GetNonce(systemname)
	sendTransferTx(types.Transfer, systemname, contractAddr, nonce, 1, big.NewInt(100), nil, systemprikey)
}

func sendTransferTransaction() {
	jww.INFO.Println("test sendTransferTransaction... ")
	input, err := formTransferAssetInput(abifile, common.BytesToAddress([]byte("testtest2")), big.NewInt(1), big.NewInt(10))
	if err != nil {
		jww.INFO.Println("sendDeployContractTransaction formCreateContractInput error ... ", err)
		return
	}
	nonce, _ := tc.GetNonce(systemname)
	sendTransferTx(types.CallContract, systemname, contractAddr, nonce, assetID, big.NewInt(0), input, systemprikey)
}

func sendTransferTx(txType types.ActionType, from, to common.Name, nonce, assetID uint64, value *big.Int, input []byte, prikey *ecdsa.PrivateKey) {
	action := types.NewAction(txType, from, to, nonce, assetID, gasLimit, value, input)
	gasprice, _ := tc.GasPrice()
	tx := types.NewTransaction(1, gasprice, action)

	signer := types.MakeSigner(big.NewInt(1))
	err := types.SignAction(action, tx, signer, prikey)
	if err != nil {
		jww.ERROR.Fatalln(err)
	}

	rawtx, err := rlp.EncodeToBytes(tx)
	if err != nil {
		jww.ERROR.Fatalln(err)
	}
	hash, _ := tc.SendRawTx(rawtx)
	jww.INFO.Println("result hash: ", hash.Hex())
}

func main() {
	sendDeployContractTransaction()
	time.Sleep(time.Duration(3) * time.Second)
	sendIssueTransaction()
	time.Sleep(time.Duration(3) * time.Second)

	sendIncreaseIssueTransaction()
	time.Sleep(time.Duration(3) * time.Second)

	sendSetOwnerIssueTransaction()
	time.Sleep(time.Duration(3) * time.Second)

	sendTransferToContractTransaction()
	time.Sleep(time.Duration(3) * time.Second)

	sendTransferTransaction()
}
