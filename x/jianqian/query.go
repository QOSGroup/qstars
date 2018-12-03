// Copyright 2018 The QOS Authors

package jianqian

import (
	"fmt"
	"github.com/QOSGroup/qstars/client/context"
	"github.com/tendermint/go-amino"
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

func ListInvestors(ctx *context.CLIContext, cdc *amino.Codec, articleHash string) ([]Investor, error) {
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
