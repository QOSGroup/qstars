package comm

import (
	"github.com/QOSGroup/qbase/context"
	ctx "github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/baseapp"
	go_amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
)


type JianQianStub struct {
}

func NewJianQianStub() JianQianStub {
	return JianQianStub{}
}

func (cstub JianQianStub) StartX(base *baseapp.QstarsBaseApp) error {
	return nil
}
func (cstub JianQianStub) EndBlockNotify(ctx context.Context) {

}

func (cstub JianQianStub) RegisterCdc(cdc *go_amino.Codec) {
	cdc.RegisterConcrete(&JianQianTx{}, "jianqian/JianQianTx", nil)
}


func (cstub JianQianStub) ResultNotify(ctx context.Context, txQcpResult interface{}) *types.Result {
	return nil
}

func (cstub JianQianStub) CustomerQuery(ctx ctx.Context, route []string, req abci.RequestQuery) (res []byte, err types.Error) {
	return nil, nil
}

func (cstub JianQianStub) Name() string {
	return "JianQianStub"
}
