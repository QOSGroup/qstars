// Copyright 2018 The QOS Authors

package jianqian

import (
	"encoding/hex"
	//"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qstars/client/context"
	"github.com/tendermint/go-amino"
	"log"
	"strings"
)

func QueryArticle(cdc *amino.Codec, ctx *context.CLIContext, hash string) (article *Articles, err error) {
	res, err := ctx.QueryStore([]byte(hash), ArticlesMapperName)
	if err != nil {
		return nil, err
	}

	err = cdc.UnmarshalBinaryBare(res, &article)

	return
}

func QueryAllAcution(cdc *amino.Codec, ctx *context.CLIContext, hash string) (auction AuctionMap, err error) {
	res, err := ctx.QueryStore([]byte(hash), AuctionMapperName)
	if err != nil {
		return nil, err
	}
	result := string(res)
	first := strings.Index(result, "{")
	res = res[first:]
	err = json.Unmarshal(res, &auction)
	return
}

func QueryCoins(cdc *amino.Codec, ctx *context.CLIContext, tx string) (coins *Coins, err error) {
	fmt.Println("tx=", tx)
	res, err := ctx.QueryStore([]byte(tx), CoinsMapperName)
	if err != nil {
		return nil, err
	}
	err = cdc.UnmarshalBinaryBare(res, &coins)
	return
}

func QueryBlance(cdc *amino.Codec, ctx *context.CLIContext, tx string) (acc *AOETokens, err error) {
	fmt.Println("tx=", tx)
	res, err := ctx.QueryStore([]byte(tx), AoeAccountMapperName)
	if err != nil {
		return nil, err
	}
	err = cdc.UnmarshalBinaryBare(res, &acc)
	return
}

type Result struct {
	Type  string
	Value interface{}
	Error string
}

func QueryTx(cdc *amino.Codec, ctx *context.CLIContext, txstring string) string {
	result := Result{}
	hash, err := hex.DecodeString(txstring)
	if err != nil {
		result.Error = err.Error()
		resp, _ := json.Marshal(result)
		return string(resp)
	}
	//
	resTx, err := ctx.Client.Tx(hash, !ctx.TrustNode)
	if err != nil {
		result.Error = err.Error()
		resp, _ := json.Marshal(result)
		return string(resp)
	}
	//parse Tx
	var tx *txs.TxStd
	err = cdc.UnmarshalBinaryBare(resTx.Tx, &tx)
	if err != nil {
		result.Error = err.Error()
		resp, _ := json.Marshal(result)
		return string(resp)
	}

	resp, _ := cdc.MarshalJSON(tx.ITxs)
	err = json.Unmarshal(resp, &result)
	if err != nil {
		result.Error = err.Error()
		resp, _ :=json.Marshal(result)
		return string(resp)
	}
	result.Error = ""
	temp, _ := json.Marshal(result)
	return string(temp)
}

func ListInvestors(ctx *context.CLIContext, cdc *amino.Codec, articleHash string) (Investors, error) {
	log.Printf("ListInvestors ctx:%+v, articleHash:%s", ctx, articleHash)
	d, err := ctx.QueryInvestadCustom([]byte(articleHash))
	if err != nil {
		return nil, err
	}
	var investors []Investor
	if err := cdc.UnmarshalBinaryBare(d, &investors); err != nil {
		return nil, err
	}
	return investors, nil
}

func QueryArticleBuyer(cdc *amino.Codec, ctx *context.CLIContext, hash string) (buyer *Buyer, err error) {
	res, err := ctx.QueryStore([]byte(hash), BuyMapperName)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, errors.New("not found")
	}

	err = cdc.UnmarshalBinaryBare(res, &buyer)
	log.Printf("jianqian.QueryArticleBuyer buyer:%+v, key:%+v, res:%+v,err:%+v", buyer, hash, res, err)

	return
}
