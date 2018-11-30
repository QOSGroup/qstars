package coins

import (
	"strconv"

	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/baseapp"
	"github.com/QOSGroup/qstars/x/common"
	"github.com/QOSGroup/qstars/x/jianqian"
	"github.com/prometheus/common/log"
	go_amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"

)

const QSCResultMapperName = "coinsResult"
const COINNAME = "AOE"

type CoinsStub struct {
	baseapp.BaseXTransaction
}

func NewCoinsStub() CoinsStub {
	return CoinsStub{}
}

func (cstub CoinsStub) StartX(base *baseapp.QstarsBaseApp) error {
	var coinsMapper = jianqian.NewCoinsMapper(jianqian.CoinsMapperName)
	base.Baseapp.RegisterMapper(coinsMapper)

	var qosMapper = jianqian.NewCoinsMapper(QSCResultMapperName)
	base.Baseapp.RegisterMapper(qosMapper)
	return nil
}
func (cstub CoinsStub) EndBlockNotify(ctx context.Context) {

}

func (cstub CoinsStub) RegisterCdc(cdc *go_amino.Codec) {
	cdc.RegisterConcrete(&CoinAOETx{}, "jianqian/CoinAOETx", nil)
	cdc.RegisterConcrete(&DispatchAOETx{}, "jianqian/DispatchAOETx", nil)
}

func (cstub CoinsStub) ResultNotify(ctx context.Context, txQcpResult interface{}) *types.Result {
	in := txQcpResult.(*txs.QcpTxResult)
	log.Debugf("ResultNotify QcpOriginalSequence:%s, result:%+v", string(in.QcpOriginalSequence), txQcpResult)
	var resultCode types.ABCICodeType
	qcpTxResult, ok := baseabci.ConvertTxQcpResult(txQcpResult)
	if ok == false {
		log.Errorf("ResultNotify ConvertTxQcpResult error.")
		resultCode = types.ABCICodeType(types.CodeTxDecode)
	} else {
		log.Errorf("ResultNotify update status")

		orginalTxHash := in.QcpOriginalExtends //orginalTx.abc
		kvMapper := ctx.Mapper(QSCResultMapperName).(*common.KvMapper)
		initValue := ""
		kvMapper.Get([]byte(orginalTxHash), &initValue)
		if initValue != "-1" {
			log.Info("This is not my response.")
			return nil
		}
		//put result to map for client query
		c := strconv.FormatInt((int64)(qcpTxResult.Result.Code), 10)
		c = c + " " + qcpTxResult.Result.Log
		log.Errorf("--------update key:" + QSCResultMapperName + " key:" + orginalTxHash + " value:" + c)
		kvMapper.Set([]byte(orginalTxHash), c)
		resultCode = types.ABCICodeType(types.CodeOK)
	}
	rr := types.Result{
		Code: resultCode,
	}
	return &rr
}
func (cstub CoinsStub) CustomerQuery(ctx context.Context, route []string, req abci.RequestQuery) (res []byte, err types.Error){
	return nil,nil
}