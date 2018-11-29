package baseapp

import (
	"github.com/QOSGroup/qbase/mapper"
	go_amino "github.com/tendermint/go-amino"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qbase/context"
	abci "github.com/tendermint/tendermint/abci/types"
	ctx "github.com/QOSGroup/qbase/context"
)

type BaseXTransaction interface {
	mapper.IMapper
	RegisterCdc(cdc *go_amino.Codec)
	StartX(base *QstarsBaseApp) error
	ResultNotify(ctx context.Context, txQcpResult interface{}) *types.Result
	EndBlockNotify(ctx context.Context)
    CustomerQuery(ctx ctx.Context, route []string, req abci.RequestQuery) (res []byte, err types.Error)
}
