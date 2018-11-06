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
)

type BankStub struct {
	baseapp.BaseXTransaction
}

func NewBankStub() BankStub {
	return BankStub{}
}

func (kv BankStub) StartX(base *baseapp.QstarsBaseApp) error {

	var qosMapper = common.NewKvMapper(QSCResultMapperName)
	base.Baseapp.RegisterMapper(qosMapper)

	return nil
}

func (kv BankStub) RegisterCdc(cdc *go_amino.Codec) {
	cdc.RegisterConcrete(&WrapperSendTx{}, "qstars/WrapperSendTx", nil)
	qosapp.RegisterCodec(cdc)
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
		orginalSeq := qcpTxResult.QcpOriginalSequence
		orginalTx := baseabci.GetQcpMapper(ctx).GetChainOutTxs("", orginalSeq)
		if orginalTx.IsResult == true {
			log.Errorf("ResultNotify Cross chain result is not a type of result.")
			resultCode = types.ABCICodeType(types.CodeInternal)
		} else {
			log.Debugf("ResultNotify update status")

			orginalTxHash := in.QcpOriginalExtends //orginalTx.abc
			kvMapper := ctx.Mapper(QSCResultMapperName).(*common.KvMapper)
			//put result to map for client query
			c := strconv.FormatInt((int64)(qcpTxResult.Result.Code), 10)
			log.Info("--------update key:"+QSCResultMapperName+" key:"+ orginalTxHash +" value:" + c)
			kvMapper.Set([]byte(orginalTxHash), c)
			resultCode = types.ABCICodeType(types.CodeOK)
		}
	}
	rr := types.Result{
		Code: resultCode,
	}
	return &rr
}
