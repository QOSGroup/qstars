// Copyright 2018 The QOS Authors

package common

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	qbasetypes "github.com/QOSGroup/qbase/types"
)

const (
	AdvertisersTxFlag = "AdvertisersTx"

	OWNER_ADDRESS=""
)

var MAX_GAS = qbasetypes.NewInt(20000)

type QuZhuanTx struct {
	FuncName string               //方法名 路由用
	Address  []qbasetypes.Address //签名者地址
	Args     []string             //参数列表
	Gas      qbasetypes.BigInt
}

var _ txs.ITx = (*QuZhuanTx)(nil)

func (it QuZhuanTx) ValidateData(ctx context.Context) error {
	funcName := it.FuncName
	routerTx, err := getStruct(funcName, it.Args,it.Address[0])
	if err != nil {
		return err
	}
	return routerTx.ValidateData(ctx)
}

func (it QuZhuanTx) Exec(ctx context.Context) (result qbasetypes.Result, crossTxQcps *txs.TxQcp) {
	funcName := it.FuncName
	routerTx, err := getStruct(funcName, it.Args,it.Address[0])
	if err != nil {
		return
	}
	return routerTx.Exec(ctx)
}

func (it QuZhuanTx) GetSigner() []qbasetypes.Address {
	return it.Address
}

func (it QuZhuanTx) CalcGas() qbasetypes.BigInt {
	return it.Gas
}

func (it QuZhuanTx) GetGasPayer() qbasetypes.Address {
	return it.Address[0]
}

func (it QuZhuanTx) GetSignData() (ret []byte) {
	ret = append(ret, []byte(it.FuncName)...)
	for _, v := range it.Address {
		ret = append(ret, v.Bytes()...)
	}
	for _, v := range it.Args {
		ret = append(ret, []byte(v)...)
	}
	return
}
