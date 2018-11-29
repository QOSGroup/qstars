package article

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/baseapp"
	"github.com/QOSGroup/qstars/x/jianqian"
	go_amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	ctx "github.com/QOSGroup/qbase/context"
)

const ArticlesMapper = "article"

type AricleStub struct {
	baseapp.BaseXTransaction
}

func NewCoinsStub() AricleStub {
	return AricleStub{}
}

func (astub AricleStub) StartX(base *baseapp.QstarsBaseApp) error {
	var aricleMapper = jianqian.NewArticlesMapper(ArticlesMapper)
	base.Baseapp.RegisterMapper(aricleMapper)
	return nil
}
func (astub AricleStub) EndBlockNotify(ctx context.Context) {

}

func (astub AricleStub) RegisterCdc(cdc *go_amino.Codec) {
	cdc.RegisterConcrete(&ArticleTx{}, "jianqian/ArticleTx", nil)
}

func (astub AricleStub) ResultNotify(ctx context.Context, txQcpResult interface{}) *types.Result {

	return nil
}

func (kv AricleStub) CustomerQuery(ctx ctx.Context, route []string, req abci.RequestQuery) (res []byte, err types.Error){
	return nil,nil
}
