package investad

import (
	"github.com/QOSGroup/qbase/context"
	ctx "github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/baseapp"
	"github.com/QOSGroup/qstars/x/jianqian"
	go_amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	"log"
)

type InvestadStub struct {


}

func NewStub() InvestadStub {
	return InvestadStub{}
}

func (s InvestadStub) StartX(base *baseapp.QstarsBaseApp) error {
	var investMapper = jianqian.NewInvestMapper(jianqian.InvestMapperName)
	base.Baseapp.RegisterMapper(investMapper)

	//var investUncheckedMapper = jianqian.NewInvestUncheckedMapper(jianqian.InvestUncheckedMapperName)
	//base.Baseapp.RegisterMapper(investUncheckedMapper)

	return nil
}

func (s InvestadStub) RegisterCdc(cdc *go_amino.Codec) {
	cdc.RegisterConcrete(&InvestTx{}, "qstars/InvestTx", nil)
}

func (s InvestadStub) ResultNotify(ctx context.Context, txQcpResult interface{}) *types.Result {
	return nil
}

func (s InvestadStub) EndBlockNotify(ctx context.Context) {

}
func (s InvestadStub) CustomerQuery(ctx ctx.Context, route []string, req abci.RequestQuery) (res []byte, err types.Error) {
	log.Printf("investad.CustomerQuery route:%+v", route)

	// jianqian, investad
	if len(route) != 2 {
		return nil, nil
	}

	if route[0] != "jianqian" || route[1] != "investad" {
		return nil, nil
	}

	key := req.Data
	if key == nil || len(key) == 0 {
		return nil, nil
	}

	investMapper := ctx.Mapper(jianqian.InvestMapperName).(*jianqian.InvestMapper)
	log.Printf("investad.CustomerQuery investMapper:%+v", investMapper)
	result := investMapper.AllInvestors([]byte(key))
	log.Printf("investad.CustomerQuery key:%+v, result:%+v", key, result)

	return investMapper.EncodeObject(result), nil
}

func (s InvestadStub) Name() string {
	return "InvestadStub"
}
