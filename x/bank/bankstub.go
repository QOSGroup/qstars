package bank

import (
	"fmt"
	bctypes "github.com/QOSGroup/qbase/example/basecoin/types"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qstars/baseapp"
	go_amino "github.com/tendermint/go-amino"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qbase/context"
)

type BankStub struct {
	baseapp.BaseContract

}


func NewBankStub() BankStub {
	return BankStub{}
}


func (kv BankStub) StartX(base *baseapp.QstarsBaseApp) error{
	return nil
}

func (kv BankStub) RegisterKVCdc(cdc *go_amino.Codec) {
	cdc.RegisterConcrete(&SendTx{}, "basecoin/SendTx",nil)
	cdc.RegisterConcrete(&bctypes.AppAccount{}, "basecoin/AppAccount", nil)
}

func (kv BankStub) ResultNotify(ctx context.Context, txQcpResult interface{}) *types.Result{
	in := txQcpResult.(*txs.QcpTxResult)
	fmt.Println("QcpOriginalSequence:"+string(in.QcpOriginalSequence))
	ok:= types.ABCICodeType(types.CodeOK)
	rr := types.Result{
		Code:ok,
	}
	return &rr
}