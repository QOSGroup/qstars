package advertisers

import (
"github.com/QOSGroup/qbase/context"
ctx "github.com/QOSGroup/qbase/context"
"github.com/QOSGroup/qbase/types"
"github.com/QOSGroup/qstars/baseapp"
	"github.com/QOSGroup/qstars/x/jianqian"
	go_amino "github.com/tendermint/go-amino"
abci "github.com/tendermint/tendermint/abci/types"
)


type AdvertisersStub struct {
}

func NewAdvertisersStub() AdvertisersStub {
	return AdvertisersStub{}
}

func (cstub AdvertisersStub) StartX(base *baseapp.QstarsBaseApp) error {
	var advertisersMapper = jianqian.NewAdvertisersMapper(jianqian.AdvertisersMapperName)
	base.Baseapp.RegisterMapper(advertisersMapper)
	return nil
}
func (cstub AdvertisersStub) EndBlockNotify(ctx context.Context) {

}

func (cstub AdvertisersStub) RegisterCdc(cdc *go_amino.Codec) {
	cdc.RegisterConcrete(&AdvertisersTx{}, "jianqian/AdvertisersTx", nil)
}

func (cstub AdvertisersStub) ResultNotify(ctx context.Context, txQcpResult interface{}) *types.Result {
	return nil
}

func (cstub AdvertisersStub) CustomerQuery(ctx ctx.Context, route []string, req abci.RequestQuery) (res []byte, err types.Error) {
	return nil, nil
}

func (cstub AdvertisersStub) Name() string {
	return "AdvertisersStub"
}
