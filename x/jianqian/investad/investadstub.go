package investad

import (
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	ctx "github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/baseapp"
	"github.com/QOSGroup/qstars/x/common"
	"github.com/QOSGroup/qstars/x/jianqian"
	go_amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	"log"
	"strconv"
)

type InvestadStub struct {
}

func NewStub() InvestadStub {
	return InvestadStub{}
}

func (s InvestadStub) StartX(base *baseapp.QstarsBaseApp) error {
	var investMapper = jianqian.NewInvestMapper(jianqian.InvestMapperName)
	base.Baseapp.RegisterMapper(investMapper)

	var investUncheckedMapper = jianqian.NewInvestUncheckedMapper(jianqian.InvestUncheckedMapperName)
	base.Baseapp.RegisterMapper(investUncheckedMapper)

	return nil
}

func (s InvestadStub) RegisterCdc(cdc *go_amino.Codec) {
	cdc.RegisterConcrete(&InvestTx{}, "qstars/InvestTx", nil)
}

func (s InvestadStub) ResultNotify(ctx context.Context, txQcpResult interface{}) *types.Result {
	in := txQcpResult.(*txs.QcpTxResult)
	log.Printf("investad.InvestadStub ResultNotify QcpOriginalSequence:%s, result:%+v", string(in.QcpOriginalSequence), txQcpResult)
	var resultCode types.ABCICodeType
	qcpTxResult, ok := baseabci.ConvertTxQcpResult(txQcpResult)
	if ok == false {
		log.Printf("investad.InvestadStub ResultNotify ConvertTxQcpResult error.")
		return nil
	} else {
		resultCode = qcpTxResult.Result.Code
		key := in.QcpOriginalExtends //orginalTx.abc

		kvMapper := ctx.Mapper(common.QSCResultMapperName).(*common.KvMapper)
		initValue := ""
		kvMapper.Get([]byte(key), &initValue)
		log.Printf("investad.InvestadStub kvMapper-1 key:[%s], value:[%s]", key, initValue)
		if initValue != s.Name() {
			log.Printf("investad.InvestadStub This is not my response.")
			return nil
		}
		c := strconv.FormatInt((int64)(qcpTxResult.Result.Code), 10)
		c = c + " " + qcpTxResult.Result.Log
		log.Printf("investad.InvestadStub kvMapper-2 key:[%s], value:[%s]", key, c)

		kvMapper.Set([]byte(key), c)

		if qcpTxResult.Result.IsOK() {
			log.Printf("investad.InvestadStub ResultNotify update status")

			investUncheckedMapper := ctx.Mapper(jianqian.InvestUncheckedMapperName).(*jianqian.InvestUncheckedMapper)
			log.Printf("investad.InvestadStub investUncheckedMapper :%+v", investUncheckedMapper)

			investUncheckeds, ok := investUncheckedMapper.GetInvestUncheckeds([]byte(key))
			if !ok || investUncheckeds == nil {
				log.Printf("investad.InvestadStub This is not my response.")
				return nil
			}

			investMapper := ctx.Mapper(jianqian.InvestMapperName).(*jianqian.InvestMapper)
			log.Printf("investad.InvestadStub investMapper :%+v", investMapper)

			for k, v := range investUncheckeds {
				log.Printf("investad.InvestadStub investUncheckeds k:%+v, v:%+v\n", k, v)
				if !v.IsChecked {
					key := jianqian.GetInvestKey(v.Article, v.Address)
					investor, ok := investMapper.GetInvestor(key)
					if ok {
						investor.Invest = investor.Invest.Add(v.Invest)
						investor.InvestTime = v.InvestTime
						log.Printf("investad.InvestadStub investor update %+v\n", investor)
						investMapper.SetInvestor(key, investor)
					} else {
						investor = jianqian.Investor{
							Address:    v.Address,
							InvestTime: v.InvestTime,
							Invest:     v.Invest,
						}
						log.Printf("investad.InvestadStub investor create %+v\n", investor)

						investMapper.SetInvestor(key, investor)
					}

					investUncheckeds[k].IsChecked = true
				}
			}
			investUncheckedMapper.SetInvestUncheckeds([]byte(key), investUncheckeds)

			resultCode = types.ABCICodeType(types.CodeOK)
		}
	}

	rr := types.Result{
		Code: resultCode,
	}

	return &rr
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

	investMapper := ctx.Mapper(jianqian.InvestMapperName).(*jianqian.InvestMapper)
	log.Printf("investad.CustomerQuery investMapper:%+v", investMapper)
	result := investMapper.AllInvestors([]byte(key))
	log.Printf("investad.CustomerQuery key:%+v, result:%+v", key, result)

	return investMapper.EncodeObject(result), nil
}

func (s InvestadStub) Name() string {
	return "InvestadStub"
}
