package investad

import (
	"fmt"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	ctx "github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/baseapp"
	"github.com/QOSGroup/qstars/x/common"
	"github.com/QOSGroup/qstars/x/jianqian"
	"github.com/prometheus/common/log"
	go_amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	"strconv"
)

type InvestadStub struct {
	baseapp.BaseXTransaction
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
	fmt.Printf("investad.InvestadStub ResultNotify QcpOriginalSequence:%s, result:%+v", string(in.QcpOriginalSequence), txQcpResult)
	var resultCode types.ABCICodeType
	qcpTxResult, ok := baseabci.ConvertTxQcpResult(txQcpResult)
	if ok == false {
		fmt.Printf("ResultNotify ConvertTxQcpResult error.")
		resultCode = types.ABCICodeType(types.CodeTxDecode)
	} else {
		resultCode = qcpTxResult.Result.Code
		key := in.QcpOriginalExtends //orginalTx.abc

		kvMapper := ctx.Mapper(common.QSCResultMapperName).(*common.KvMapper)
		initValue := ""
		kvMapper.Get([]byte(key), &initValue)
		if initValue != "-1" {
			log.Info("This is not my response.")
			return nil
		}
		c := strconv.FormatInt((int64)(qcpTxResult.Result.Code), 10)
		c = c + " " + qcpTxResult.Result.Log
		kvMapper.Set([]byte(key), c)

		if qcpTxResult.Result.IsOK() {
			fmt.Printf("investad.InvestadStub ResultNotify update status")

			investUncheckedMapper := ctx.Mapper(jianqian.InvestUncheckedMapperName).(*jianqian.InvestUncheckedMapper)
			investUncheckeds, ok := investUncheckedMapper.GetInvestUncheckeds([]byte(key))
			if !ok || investUncheckeds == nil {
				fmt.Printf("This is not my response.")
				return nil
			}

			investMapper := ctx.Mapper(jianqian.InvestMapperName).(*jianqian.InvestMapper)
			for k, v := range investUncheckeds {
				if !v.IsChecked {
					key := jianqian.GetInvestKey(v.Article, v.Address)
					investor, ok := investMapper.GetInvestor(key)
					if ok {
						investor.Invest = investor.Invest.Add(v.Invest)
						investor.InvestTime = v.InvestTime
						investMapper.SetInvestor(key, investor)
					} else {
						investor = jianqian.Investor{
							Address:    v.Address,
							InvestTime: v.InvestTime,
							Invest:     v.Invest,
						}
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
	// jianqian, investad, key
	if len(route) != 3 {
		return nil, nil
	}

	if route[0] != "jianqian" || route[1] != "investad" {
		return nil, nil
	}

	key := route[2]

	investMapper := ctx.Mapper(jianqian.InvestMapperName).(*jianqian.InvestMapper)
	fmt.Printf("%+v", investMapper)
	result := investMapper.AllInvestors([]byte(key))

	return investMapper.EncodeObject(result), nil
}
