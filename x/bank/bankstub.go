package bank

import (

	"github.com/QOSGroup/qstars/baseapp"
	go_amino "github.com/tendermint/go-amino"
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

}