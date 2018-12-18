// Copyright 2018 The QOS Authors

package jianqian

import (
	"errors"
	"fmt"
	"github.com/QOSGroup/qstars/client/context"
	"github.com/tendermint/go-amino"
	"log"
)

func QueryArticle(cdc *amino.Codec, ctx *context.CLIContext, hash string) (article *Articles, err error) {
	res, err := ctx.QueryStore([]byte(hash), ArticlesMapperName)
	if err != nil {
		return nil, err
	}

	err = cdc.UnmarshalBinaryBare(res, &article)

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
