package recharge

import (
	"github.com/QOSGroup/qbase/context"
	ctx "github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/baseapp"
	go_amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
)


type RechargeStub struct {
}

func NewRechargeStub() RechargeStub {
	return RechargeStub{}
}

func (cstub RechargeStub) StartX(base *baseapp.QstarsBaseApp) error {
	return nil
}
func (cstub RechargeStub) EndBlockNotify(ctx context.Context) {

}

func (cstub RechargeStub) RegisterCdc(cdc *go_amino.Codec) {
	cdc.RegisterConcrete(&RechargeTx{}, "jianqian/RechargeTx", nil)
}

func (cstub RechargeStub) ResultNotify(ctx context.Context, txQcpResult interface{}) *types.Result {
	return nil
}

func (cstub RechargeStub) CustomerQuery(ctx ctx.Context, route []string, req abci.RequestQuery) (res []byte, err types.Error) {
	return nil, nil
}

func (cstub RechargeStub) Name() string {
	return "RechargeStub"
}
