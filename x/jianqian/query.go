package jianqian

import (
	"github.com/QOSGroup/qstars/client/context"
	"github.com/tendermint/go-amino"
)

func QueryArticle(cdc *amino.Codec,ctx *context.CLIContext, hash string) (article *Articles,err error) {
	res, err := ctx.QueryStore([]byte(hash),ArticlesMapperName)
	if err != nil {
		return nil, err
	}
	err =cdc.UnmarshalBinaryBare(res, &article)
	return
}
