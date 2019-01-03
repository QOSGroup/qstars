package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/txs"
	qbasetypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/star"
	"github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/utility"
	"github.com/QOSGroup/qstars/wire"
	"github.com/QOSGroup/qstars/x/bank/tx"
	"github.com/pkg/errors"
	"strings"

	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
	sdk "github.com/QOSGroup/qstars/types"
)

//func main1(){
//	InitJNI()
//	//send --from=rpt3O80wAFI1+ZqNYt8DqJ5PaQ+foDq7G/InFfycoFYT8tgGFJLp+BSVELW2fTQNGZ/yTzTIXbu9fg33gOmmzA== --to=address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355 --amount=2qos
//	r := SendByJNI("rpt3O80wAFI1+ZqNYt8DqJ5PaQ+foDq7G/InFfycoFYT8tgGFJLp+BSVELW2fTQNGZ/yTzTIXbu9fg33gOmmzA==","address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355","2qos")
//	print(r)
//}

func main(){
	cdc := star.MakeCodec()
	tostrs := []qbasetypes.Address {}
	toaddressesStr := "address13hkg8nva06hntmnhfupy29c2l9aq9zs879jhez;address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355"
	toaddresses := strings.Split(toaddressesStr,";")
	fmt.Println("to address:",toaddresses)
	for i:=0;i< len(toaddresses);i++{
		to, err := sdk.AccAddressFromBech32(toaddresses[i])
		if err != nil {
			fmt.Println(err)
		}
		tostrs = append(tostrs, to)
	}

	fromstrs := []string{}
	fromaddressesStr := "Ey+2bNFF2gTUV6skSBgRy3rZwo9nS4Dw0l2WpLrhVvV8MuMRbjN4tUK8orHiJgHTR+enkxyXcA8giVrsrIRM4Q==;31PlT2p6UICjV63dG7Nh3Mh9W0b+7FAEU+KOAxyNbZ29rwqNzxQJlQPh59tZpbS1EdIT6TE5N6L72se9BUe9iw=="

	fromaddresses := strings.Split(fromaddressesStr,";")
	fmt.Println("private key:",fromaddresses)
	for i:=0;i< len(fromaddresses);i++{
		fromstrs = append(fromstrs, fromaddresses[i])
	}

	scoins :=[]types.Coins {}
	fromamountstr := "100AOE,4qos;2qos"
	// parse coins trying to be sent
	fromamounts := strings.Split(fromamountstr,";")
	fmt.Println("from mounts:",fromamounts)
	for i:=0;i< len(fromamounts);i++ {
		secointmp, err := sdk.ParseCoins(fromamounts[i])
		if err != nil {
			fmt.Println(err)
		}
		scoins = append(scoins, secointmp)
	}


	rcoins :=[]types.Coins {}
	toamountstr := "100AOE,1qos;5qos"
	toamounts := strings.Split(toamountstr,";")
	fmt.Println("to mounts:",toamounts)
	for i:=0;i< len(toamounts);i++ {
		// parse coins trying to be sent
		recointmp, err := sdk.ParseCoins(toamounts[i])
		if err != nil {
			fmt.Println(err)
		}
		rcoins = append(rcoins, recointmp)
	}

	 err := Send(cdc, fromstrs, tostrs, scoins, rcoins)
	if err != nil {
		fmt.Println(err)
	}
}
// Send 支持一次多种币 coins.Len() == 1;
func Send(cdc *wire.Codec, fromstrs []string, tos []qbasetypes.Address, scoins []types.Coins, rcoins []types.Coins) error {
	if (len(scoins) == 0)||(len(rcoins)==0) {
		return  errors.New("coins不能为空")
	}
	privs := []ed25519.PrivKeyEd25519{}
	froms := []qbasetypes.Address{}
	for _,fromstr := range fromstrs{
		_, addrben32, priv := utility.PubAddrRetrievalFromAmino(fromstr, cdc)
		from, err := types.AccAddressFromBech32(addrben32)
		account.AddressStoreKey(from)
		if err != nil {
			return  err
		}
		privs = append(privs, priv)
		froms = append(froms,from)
	}


	var	sendcoins [][]qbasetypes.BaseCoin
	for _, echosendercoin := range scoins{
		var sendcoinstmp []qbasetypes.BaseCoin
		for _, coin := range echosendercoin {
			sendcoinstmp = append(sendcoinstmp, qbasetypes.BaseCoin{
				Name:   coin.Denom,
				Amount: qbasetypes.NewInt(coin.Amount.Int64()),
			})
		}
		sendcoins = append(sendcoins,sendcoinstmp)
	}

	var	receivecoins [][]qbasetypes.BaseCoin
	for _, echoreceivecoin := range rcoins {
		var receivecoinstmp []qbasetypes.BaseCoin
		for _, coin := range echoreceivecoin {
			receivecoinstmp = append(receivecoinstmp, qbasetypes.BaseCoin{
				Name:   coin.Denom,
				Amount: qbasetypes.NewInt(coin.Amount.Int64()),
			})
		}
		receivecoins = append(receivecoins,receivecoinstmp)
	}

	var nn int64
	nn = int64(35)
	nn++
	nns := []int64{}
	nns = append(nns, nn)
	nn = 1
	nns = append(nns, nn)

	t := tx.NewTransferMultiple(froms, tos, sendcoins,receivecoins)
	var msg *txs.TxStd
	tochainid := "qos-test"
	fromchainid := "qstars-test"
	msg = genStdSendMultiTx(cdc, t, privs, fromchainid,tochainid, nns)
	//--------------------------------------------------------------------
	rrr, _ := cdc.MarshalJSON(msg)
	fmt.Println()
	fmt.Println((string)(rrr))

	//--------------------------------------------------------------------
	ccc, _ := cdc.MarshalBinaryBare(msg)
	fmt.Println()
	fmt.Println(" \n",ccc)
	fmt.Println()

	for _,fromstr1 := range fromstrs{
		fmt.Println("Private key:         ", fromstr1)
	}
	fmt.Println("nonce:              ",nn)
	//-----------------------------Signature-----------------------------------
	i := 0
	for _,nn := range nns{
		//sigdata := append(msg.BuildSignatureBytes(nn,tochainid), Int2Byte(nn)...)
		//sigdata,_ := cdc.MarshalBinaryBare(msg.Signature)
		//-------------------------------------------------------------------------

		privateStr := hex.EncodeToString(froms[i].Bytes())
		fmt.Println("from:            ", (privateStr))

		sigdata := msg.BuildSignatureBytes(nn, tochainid)
		encodedStr := hex.EncodeToString(sigdata)
		fmt.Println("Need to signdata hex:  ", (encodedStr))
		fmt.Println("Need to signdata byte: ",sigdata);

		signed1, _ := privs[i].Sign(sigdata)
		encodedStr = hex.EncodeToString(signed1)
		fmt.Println("signature hex:     ", (encodedStr))
		fmt.Println("signature byte:    ", (signed1))
		i++
	}

	for _,to := range tos {
		encodedStr := hex.EncodeToString(to.Bytes())
		fmt.Println("to: ", (encodedStr))
	}


	return nil
}


//add the string input chainid
func genStdSendTx(cdc *amino.Codec, sendTx txs.ITx, priKey ed25519.PrivKeyEd25519, fromchainid string, tochainid string, nonce int64) *txs.TxStd {
	gas := qbasetypes.NewInt(int64(0))
	stx := txs.NewTxStd(sendTx, tochainid, gas)
	signature, _ := stx.SignTx(priKey, nonce,tochainid, tochainid)
	stx.Signature = []txs.Signature{txs.Signature{
		Pubkey:    priKey.PubKey(),
		Signature: signature,
		Nonce:     nonce,
	}}

	return stx
}

//add the string input chainid
func genStdSendMultiTx(cdc *amino.Codec, sendTx txs.ITx, priKeys []ed25519.PrivKeyEd25519, fromchainid string, tochainid string, nonces []int64 ) *txs.TxStd {
	gas := qbasetypes.NewInt(int64(0))
	stx := txs.NewTxStd(sendTx, tochainid, gas)
	i := 0
	for _,priKey := range priKeys {
		fmt.Println("signed data: ",i,stx.BuildSignatureBytes(nonces[i], tochainid));
		signature, _ := stx.SignTx(priKey, nonces[i],tochainid, tochainid)
		stx.Signature = append(stx.Signature, txs.Signature{
			Pubkey:    priKey.PubKey(),
			Signature: signature,
			Nonce:     nonces[i],
		})
		i++;
	}

	return stx
}

func Int2Byte(in int64) []byte {
	var ret = bytes.NewBuffer([]byte{})
	err := binary.Write(ret, binary.BigEndian, in)
	if err != nil {
		fmt.Printf("Int2Byte error:%s", err.Error())
		return nil
	}

	return ret.Bytes()
}