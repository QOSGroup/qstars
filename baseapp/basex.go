package baseapp

import (
	"github.com/QOSGroup/qbase/context"
	ctx "github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/types"
	go_amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
)

type BaseXTransaction interface {
	RegisterCdc(cdc *go_amino.Codec)
	StartX(base *QstarsBaseApp) error
	ResultNotify(ctx context.Context, txQcpResult interface{}) *types.Result
	EndBlockNotify(ctx context.Context)
	CustomerQuery(ctx ctx.Context, route []string, req abci.RequestQuery) (res []byte, err types.Error)
	Name() string
}
