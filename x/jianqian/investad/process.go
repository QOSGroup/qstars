// Copyright 2018 The QOS Authors

package investad

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/QOSGroup/qbase/txs"
	qbasetypes "github.com/QOSGroup/qbase/types"
	qostxs "github.com/QOSGroup/qos/txs"
	"github.com/QOSGroup/qstars/client/utils"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/utility"
	"github.com/QOSGroup/qstars/wire"
	"log"
)

type ResultInvest struct {
	Code   string          `json:"code"`
	Reason string          `json:"reason,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
}

func InternalError(reason string) ResultInvest {
	return ResultInvest{Code: "-1", Reason: reason}
}

func NewResultInvest(cdc *wire.Codec, code, reason string, res interface{}) ResultInvest {
	var rawMsg json.RawMessage

	if res != nil {
		var js []byte
		js, err := cdc.MarshalJSON(res)
		if err != nil {
			return InternalError(err.Error())
		}
		rawMsg = json.RawMessage(js)
	}

	var result ResultInvest
	result.Result = rawMsg
	result.Code = code
	result.Reason = reason

	return result
}

func (ri ResultInvest) Marshal() string {
	jsonBytes, err := json.MarshalIndent(ri, "", "  ")
	if err != nil {
		log.Printf("InvestAd err:%s", err.Error())
		return InternalError(err.Error()).Marshal()
	}
	return string(jsonBytes)
}

const coinsName = "AOE"

var tempAddr = qbasetypes.Address("address1wmrup5xemdxzx29jalp5c98t7mywulg8wgxxxx")

// InvestAdBackground 提交到链上
func InvestAdBackground(cdc *wire.Codec, txb string) string {
	result := &ResultInvest{}
	result.Code = "0"

	ts := new(txs.TxStd)
	err := cdc.UnmarshalJSON([]byte(txb), ts)
	fmt.Printf("InvestAdBackground ts:%+v, txb:%s\n", ts, txb)
	if err != nil {
		return InternalError(err.Error()).Marshal()
	}

	cliCtx := *config.GetCLIContext().QSCCliContext
	_, commitresult, err := utils.SendTx(cliCtx, cdc, ts)
	fmt.Printf("SendTx commitresult:%+v, err:%+v \n", commitresult, err)
	if err != nil {
		return InternalError(err.Error()).Marshal()
	}

	//result.Result = []byte(response)
	//height := strconv.FormatInt(commitresult.Height, 10)
	//result.Heigth = height
	return NewResultInvest(cdc, "0", "", commitresult).Marshal()
}

// InvestAd 投资广告
func InvestAd(cdc *wire.Codec, chainId, articleHash, coins, privatekey string, nonce int64) string {
	var result ResultInvest
	result.Code = "0"

	tx, err := investAd(cdc, chainId, articleHash, coins, privatekey, nonce)
	if err != nil {
		log.Printf("investAd err:%s", err.Error())
		result.Code = "-1"
		result.Reason = err.Error()
		return result.Marshal()
	}

	js, err := cdc.MarshalJSON(tx)
	if err != nil {
		log.Printf("investAd err:%s", err.Error())
		result.Code = "-1"
		result.Reason = err.Error()
		return result.Marshal()
	}
	result.Result = json.RawMessage(js)

	return result.Marshal()
}

// investAd 投资广告
func investAd(cdc *wire.Codec, chainId, articleHash, coins, privatekey string, nonce int64) (*txs.TxStd, error) {
	cs, err := types.ParseCoins(coins)
	if err != nil {
		return nil, err
	}

	for _, v := range cs {
		if v.Denom != coinsName {
			return nil, errors.New("only support AOE")
		}
	}

	_, addrben32, priv := utility.PubAddrRetrievalFromAmino(privatekey, cdc)
	investor, err := types.AccAddressFromBech32(addrben32)
	var ccs []qbasetypes.BaseCoin
	for _, coin := range cs {
		ccs = append(ccs, qbasetypes.BaseCoin{
			Name:   coin.Denom,
			Amount: qbasetypes.NewInt(coin.Amount.Int64()),
		})
	}
	nonce++

	transferTx := NewTransfer(investor, tempAddr, ccs)
	// TODO set zero, temp
	gas := qbasetypes.NewInt(int64(0))
	stx := txs.NewTxStd(transferTx, chainId, gas)
	signature, _ := stx.SignTx(priv, nonce)
	stx.Signature = []txs.Signature{txs.Signature{
		Pubkey:    priv.PubKey(),
		Signature: signature,
		Nonce:     nonce,
	}}

	it := &InvestTx{}
	it.ArticleHash = []byte(articleHash)
	it.Std = stx
	tx2 := txs.NewTxStd(it, chainId, stx.MaxGas)

	return tx2, nil
}

func warpperTransItem(addr qbasetypes.Address, coins []qbasetypes.BaseCoin) qostxs.TransItem {
	var ti qostxs.TransItem
	ti.Address = addr
	ti.QOS = qbasetypes.NewInt(0)

	for _, coin := range coins {
		if coin.Name == "qos" {
			ti.QOS = ti.QOS.Add(coin.Amount)
		} else {
			ti.QSCs = append(ti.QSCs, &coin)
		}
	}

	return ti
}

// NewTransfer ...
func NewTransfer(sender qbasetypes.Address, receiver qbasetypes.Address, coin []qbasetypes.BaseCoin) qostxs.TransferTx {
	var sendTx qostxs.TransferTx

	sendTx.Senders = append(sendTx.Senders, warpperTransItem(sender, coin))
	sendTx.Receivers = append(sendTx.Receivers, warpperTransItem(receiver, coin))

	return sendTx
}
