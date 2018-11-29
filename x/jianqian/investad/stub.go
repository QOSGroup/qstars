package investad

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/baseapp"
	"github.com/QOSGroup/qstars/x/common"
	"github.com/QOSGroup/qstars/x/jianqian"
	go_amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	ctx "github.com/QOSGroup/qbase/context"
)

type Stub struct {
	baseapp.BaseXTransaction
}

func NewStub() Stub {
	return Stub{}
}

func (s Stub) StartX(base *baseapp.QstarsBaseApp) error {
	var qosMapper = common.NewKvMapper(jianqian.InvestMapperName)
	base.Baseapp.RegisterMapper(qosMapper)

	return nil
}

func (s Stub) RegisterCdc(cdc *go_amino.Codec) {
	cdc.RegisterConcrete(&InvestTx{}, "qstars/InvestTx", nil)
}

func (s Stub) ResultNotify(ctx context.Context, txQcpResult interface{}) *types.Result {
	return nil
}
func (kv Stub) CustomerQuery(ctx ctx.Context, route []string, req abci.RequestQuery) (res []byte, err types.Error){
	return nil,nil
}
