package common

import (
	"github.com/QOSGroup/qbase/context"
	ctx "github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/baseapp"
	"github.com/QOSGroup/qstars/x/quzhuan"
	go_amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
)


type QuZhuanStub struct {
}

func NewJianQianStub() QuZhuanStub {
	return QuZhuanStub{}
}

func (cstub QuZhuanStub) StartX(base *baseapp.QstarsBaseApp) error {
	var userMapper = quzhuan.NewUsersMapper(quzhuan.UsersMapperName)
	base.Baseapp.RegisterMapper(userMapper)
	var scenesMapper = quzhuan.NewUsersMapper(quzhuan.ScenesMapperName)
	base.Baseapp.RegisterMapper(scenesMapper)
	var coinsMapper = quzhuan.NewUsersMapper(quzhuan.CoinsMapperName)
	base.Baseapp.RegisterMapper(coinsMapper)

	return nil
}
func (cstub QuZhuanStub) EndBlockNotify(ctx context.Context) {

}

func (cstub QuZhuanStub) RegisterCdc(cdc *go_amino.Codec) {
	cdc.RegisterConcrete(&QuZhuanTx{}, "quzhuan/QuZhuanTx", nil)
}


func (cstub QuZhuanStub) ResultNotify(ctx context.Context, txQcpResult interface{}) *types.Result {
	return nil
}

func (cstub QuZhuanStub) CustomerQuery(ctx ctx.Context, route []string, req abci.RequestQuery) (res []byte, err types.Error) {
	return nil, nil
}

func (cstub QuZhuanStub) Name() string {
	return "JianQianStub"
}
