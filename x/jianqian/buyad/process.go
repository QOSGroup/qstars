// Copyright 2018 The QOS Authors

package buyad

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

type ResultBuy struct {
	Code   string          `json:"code"`
	Reason string          `json:"reason,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
}

func InternalError(reason string) ResultBuy {
	return ResultBuy{Code: "-1", Reason: reason}
}

func NewResultBuy(cdc *wire.Codec, code, reason string, res interface{}) ResultBuy {
	var rawMsg json.RawMessage

	if res != nil {
		var js []byte
		js, err := cdc.MarshalJSON(res)
		if err != nil {
			return InternalError(err.Error())
		}
		rawMsg = json.RawMessage(js)
	}

	var result ResultBuy
	result.Result = rawMsg
	result.Code = code
	result.Reason = reason

	return result
}

func (ri ResultBuy) Marshal() string {
	jsonBytes, err := json.MarshalIndent(ri, "", "  ")
	if err != nil {
		log.Printf("BuyAd err:%s", err.Error())
		return InternalError(err.Error()).Marshal()
	}
	return string(jsonBytes)
}

const coinsName = "QOS"
const tempAddr = "address1wmrup5xemdxzx29jalp5c98t7mywulg8wgxxxx"

// BuyAdBackground 提交到链上
func BuyAdBackground(cdc *wire.Codec, txb string) string {
	result := &ResultBuy{}
	result.Code = "0"

	ts := new(txs.TxStd)
	err := cdc.UnmarshalJSON([]byte(txb), ts)
	if err != nil {
		return InternalError(err.Error()).Marshal()
	}

	cliCtx := *config.GetCLIContext().QOSCliContext
	_, commitresult, err := utils.SendTx(cliCtx, cdc, ts)
	if err != nil {
		return InternalError(err.Error()).Marshal()
	}

	//result.Result = []byte(response)
	//height := strconv.FormatInt(commitresult.Height, 10)
	//result.Heigth = height
	return NewResultBuy(cdc, "0", "", commitresult).Marshal()
}

// BuyAd 投资广告
func BuyAd(cdc *wire.Codec, chainId, articleHash, coins, privatekey string, nonce int64) string {
	var result ResultBuy

	tx, err := buyAd(cdc, chainId, articleHash, coins, privatekey, nonce)
	if err != nil {
		log.Printf("buyAd err:%s", err.Error())
		result.Code = "-1"
		result.Reason = err.Error()
		return result.Marshal()
	}

	js, err := cdc.MarshalJSON(tx)
	if err != nil {
		log.Printf("buyAd err:%s", err.Error())
		result.Code = "-1"
		result.Reason = err.Error()
		return result.Marshal()
	}
	result.Result = json.RawMessage(js)

	return result.Marshal()
}

func getReceivers(articleHash string, amount int64) []qostxs.TransItem {
	return nil
}

// buyAd 投资广告
func buyAd(cdc *wire.Codec, chainId, articleHash, coins, privatekey string, nonce int64) (*txs.TxStd, error) {
	cs, err := types.ParseCoins(coins)
	if err != nil {
		return nil, err
	}

	if len(cs) != 1 {
		return nil, errors.New("one coin need")
	}

	for _, v := range cs {
		if v.Denom != coinsName {
			return nil, fmt.Errorf("only support %s", coinsName)
		}
	}

	var amount int64
	_, addrben32, priv := utility.PubAddrRetrievalFromAmino(privatekey, cdc)
	buyer, err := types.AccAddressFromBech32(addrben32)
	var ccs []qbasetypes.BaseCoin
	for _, coin := range cs {
		amount = coin.Amount.Int64()
		ccs = append(ccs, qbasetypes.BaseCoin{
			Name:   coin.Denom,
			Amount: qbasetypes.NewInt(coin.Amount.Int64()),
		})
	}
	nonce++
	var transferTx qostxs.TransferTx
	transferTx.Senders = []qostxs.TransItem{warpperTransItem(buyer, ccs)}
	transferTx.Receivers = getReceivers(articleHash, amount)
	gas := qbasetypes.NewInt(int64(0))
	stx := txs.NewTxStd(transferTx, chainId, gas)
	signature, _ := stx.SignTx(priv, nonce)
	stx.Signature = []txs.Signature{txs.Signature{
		Pubkey:    priv.PubKey(),
		Signature: signature,
		Nonce:     nonce,
	}}

	it := &BuyTx{}
	it.ArticleHash = []byte(articleHash)
	it.Std = stx
	//tx2 := txs.NewTxStd(it, config.GetCLIContext().Config.QSCChainID, stx.MaxGas)
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
