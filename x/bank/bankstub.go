package bank

import (
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	qosapp "github.com/QOSGroup/qos/app"
	"github.com/QOSGroup/qstars/baseapp"
	"github.com/QOSGroup/qstars/x/common"
	"github.com/prometheus/common/log"
	go_amino "github.com/tendermint/go-amino"
	"strconv"
	abci "github.com/tendermint/tendermint/abci/types"
	ctx "github.com/QOSGroup/qbase/context"
)

type BankStub struct {
	baseapp.BaseXTransaction
}

func NewBankStub() BankStub {
	return BankStub{}
}

func (kv BankStub) StartX(base *baseapp.QstarsBaseApp) error {

	var qosMapper = common.NewKvMapper(common.QSCResultMapperName)
	base.Baseapp.RegisterMapper(qosMapper)

	return nil
}

func (kv BankStub) RegisterCdc(cdc *go_amino.Codec) {
	cdc.RegisterConcrete(&WrapperSendTx{}, "qstars/WrapperSendTx", nil)
	qosapp.RegisterCodec(cdc)
}

func (kv BankStub) EndBlockNotify(ctx context.Context){

}

func (kv BankStub) CustomerQuery(ctx ctx.Context, route []string, req abci.RequestQuery) (res []byte, err types.Error){
	return nil,nil
}

func (kv BankStub) ResultNotify(ctx context.Context, txQcpResult interface{}) *types.Result {
	in := txQcpResult.(*txs.QcpTxResult)
	log.Debugf("ResultNotify QcpOriginalSequence:%s, result:%+v", string(in.QcpOriginalSequence), txQcpResult)
	var resultCode types.ABCICodeType
	qcpTxResult, ok := baseabci.ConvertTxQcpResult(txQcpResult)
	if ok == false {
		log.Errorf("ResultNotify ConvertTxQcpResult error.")
		resultCode = types.ABCICodeType(types.CodeTxDecode)
	} else {
		//get original cross chain transaction
		//orginalSeq := qcpTxResult.QcpOriginalSequence
		//orginalTx := baseabci.GetQcpMapper(ctx).GetChainOutTxs("", orginalSeq)
		//if orginalTx.IsResult == false {
		//	log.Errorf("ResultNotify Cross chain result is not a type of result.")
		//	resultCode = types.ABCICodeType(types.CodeInternal)
		//} else {
		//
		//}
		log.Errorf("ResultNotify update status")

		orginalTxHash := in.QcpOriginalExtends //orginalTx.abc
		kvMapper := ctx.Mapper(common.QSCResultMapperName).(*common.KvMapper)
		initValue :=""
		kvMapper.Get([]byte(orginalTxHash),&initValue)
		if initValue!="-1"{
			log.Info("This is not my response.")
			return nil
		}
		//put result to map for client query
		c := strconv.FormatInt((int64)(qcpTxResult.Result.Code), 10)
		c = c+" "+qcpTxResult.Result.Log
		log.Errorf("--------update key:"+common.QSCResultMapperName+" key:"+ orginalTxHash +" value:" + c)
		kvMapper.Set([]byte(orginalTxHash), c)
		resultCode = types.ABCICodeType(types.CodeOK)
	}
	rr := types.Result{
		Code: resultCode,
	}
	return &rr
}
