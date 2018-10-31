package bank

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/example/basecoin/tx"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qstars/baseapp"
	"github.com/QOSGroup/qstars/x/common"
	"github.com/prometheus/common/log"
	go_amino "github.com/tendermint/go-amino"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/example/basecoin/tx"
	"strconv"
)

type BankStub struct {
	baseapp.BaseContract
}

func NewBankStub() BankStub {
	return BankStub{}
}


func (kv BankStub) StartX(base *baseapp.QstarsBaseApp) error{

	var qosMapper = common.NewKvMapper(QSCResultMapperName)
	base.Baseapp.RegisterMapper(qosMapper)

	if err := base.Baseapp.LoadLatestVersion(); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (kv BankStub) RegisterKVCdc(cdc *go_amino.Codec) {
	cdc.RegisterConcrete(&tx.SendTx{}, "basecoin/SendTx", nil)
	cdc.RegisterConcrete(&WrapperSendTx{}, "qstars/WrapperSendTx", nil)
	cdc.RegisterConcrete(&account.QOSAccount{}, "qos/QOSAccount", nil)
}

func (kv BankStub) ResultNotify(ctx context.Context, txQcpResult interface{}) *types.Result {
	in := txQcpResult.(*txs.QcpTxResult)
	fmt.Println("QcpOriginalSequence:"+string(in.QcpOriginalSequence))
	var resultCode types.ABCICodeType
	qcpTxResult, ok := baseabci.ConvertTxQcpResult(txQcpResult)
	if ok == false{
		fmt.Println("ConvertTxQcpResult error.")
		resultCode = types.ABCICodeType(types.CodeTxDecode)
	}else {
		//get original cross chain transaction
		orginalSeq := qcpTxResult.QcpOriginalSequence
		orginalTx := baseabci.GetQcpMapper(ctx).GetChainOutTxs("",orginalSeq)
		if orginalTx.IsResult == false{
			log.Error("Cross chain result is not a type of result.")
			resultCode = types.ABCICodeType(types.CodeInternal)
		}else{
			//orginalTx.
			orginalTxHash := "123"  //orginalTx.abc
			kvMapper := ctx.Mapper(QSCResultMapperName).(*common.KvMapper)
			//put result to map for client query
			c := strconv.FormatInt(qcpTxResult.Code,10)
			kvMapper.Set([]byte(orginalTxHash),c)
			resultCode = types.ABCICodeType(types.CodeOK)
		}
	}
	rr := types.Result{
		Code:resultCode,

	}
	return &rr
}
