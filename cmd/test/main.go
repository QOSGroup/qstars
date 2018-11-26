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

	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
	sdk "github.com/QOSGroup/qstars/types"
)

func main(){
	InitJNI()
	//send --from=rpt3O80wAFI1+ZqNYt8DqJ5PaQ+foDq7G/InFfycoFYT8tgGFJLp+BSVELW2fTQNGZ/yTzTIXbu9fg33gOmmzA== --to=address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355 --amount=2qos
	r := SendByJNI("rpt3O80wAFI1+ZqNYt8DqJ5PaQ+foDq7G/InFfycoFYT8tgGFJLp+BSVELW2fTQNGZ/yTzTIXbu9fg33gOmmzA==","address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355","2qos")
	print(r)
}

func main1(){
	cdc := star.MakeCodec()
	to, err := sdk.AccAddressFromBech32("address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355")
	if err != nil {
		fmt.Println(err)
	}
	fromstr := "rpt3O80wAFI1+ZqNYt8DqJ5PaQ+foDq7G/InFfycoFYT8tgGFJLp+BSVELW2fTQNGZ/yTzTIXbu9fg33gOmmzA=="

	amount := "2qos"
	// parse coins trying to be sent
	coins, err := sdk.ParseCoins(amount)
	if err != nil {
		fmt.Println(err)
	}

	 err = Send(cdc, fromstr, to, coins)
	if err != nil {
		fmt.Println(err)
	}
}
// Send 支持一次多种币 coins.Len() == 1;
func Send(cdc *wire.Codec, fromstr string, to qbasetypes.Address, coins types.Coins) error {
	if coins.Len() == 0 {
		return  errors.New("coins不能为空")
	}

	_, addrben32, priv := utility.PubAddrRetrievalFromAmino(fromstr, cdc)
	from, err := types.AccAddressFromBech32(addrben32)
	account.AddressStoreKey(from)
	if err != nil {
		return  err
	}


	var ccs []qbasetypes.BaseCoin
	for _, coin := range coins {
		ccs = append(ccs, qbasetypes.BaseCoin{
			Name:   coin.Denom,
			Amount: qbasetypes.NewInt(coin.Amount.Int64()),
		})
	}



	var nn int64
	nn = int64(6)
	nn++

	t := tx.NewTransfer(from, to, ccs)
	var msg *txs.TxStd
	chainid := "qos-test"
	msg = genStdSendTx(cdc, t, priv, chainid, nn)
	rrr, _ := cdc.MarshalJSON(msg)
	fmt.Println()
	fmt.Println((string)(rrr))
	fmt.Println()

	sigdata := append(msg.GetSignData(), Int2Byte(nn)...)
	encodedStr := hex.EncodeToString(sigdata)
	fmt.Println(":need to signdata ", (encodedStr))

	encodedStr = hex.EncodeToString(from.Bytes())
	fmt.Println("from: ", (encodedStr))
	encodedStr = hex.EncodeToString(to.Bytes())
	fmt.Println("to: ", (encodedStr))

	signed1, _ := priv.Sign(sigdata)
	encodedStr = hex.EncodeToString(signed1)
	fmt.Println("sign content: ", (encodedStr))
	return nil
}


//add the string input chainid
func genStdSendTx(cdc *amino.Codec, sendTx txs.ITx, priKey ed25519.PrivKeyEd25519, chainid string, nonce int64) *txs.TxStd {
	gas := qbasetypes.NewInt(int64(0))
	stx := txs.NewTxStd(sendTx, chainid, gas)
	signature, _ := stx.SignTx(priKey, nonce)
	stx.Signature = []txs.Signature{txs.Signature{
		Pubkey:    priKey.PubKey(),
		Signature: signature,
		Nonce:     nonce,
	}}

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