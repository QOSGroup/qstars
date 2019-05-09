// Copyright 2018 The QOS Authors

package comm

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	qbasetypes "github.com/QOSGroup/qbase/types"
)

const (
	AdvertisersTxFlag = "AdvertisersTx"
	ArticleTxFlag     = "ArticleTx"
	BuyTxFlag         = "BuyTx"
	AuctionTxFlag     = "AuctionTx"
	AOETxFlag         = "AOETx"
	InvestTxFlag      = "InvestTx"
	RechargeTxFlag    = "RechargeTx"
	ExtractTxFlag     = "ExtractTx"
)

var MAX_GAS = qbasetypes.NewInt(20000)

type JianQianTx struct {
	FuncName string               //方法名 路由用
	Address  []qbasetypes.Address //签名者地址
	Args     []string             //参数列表
	Gas      qbasetypes.BigInt
}

var _ txs.ITx = (*JianQianTx)(nil)

func (it JianQianTx) ValidateData(ctx context.Context) error {
	funcName := it.FuncName
	routerTx, err := getStruct(funcName, it.Args)
	if err != nil {
		return err
	}
	return routerTx.ValidateData(ctx)
}

func (it JianQianTx) Exec(ctx context.Context) (result qbasetypes.Result, crossTxQcps *txs.TxQcp) {
	funcName := it.FuncName
	routerTx, err := getStruct(funcName, it.Args)
	if err != nil {
		return
	}
	return routerTx.Exec(ctx)
}

func (it JianQianTx) GetSigner() []qbasetypes.Address {
	return it.Address
}

func (it JianQianTx) CalcGas() qbasetypes.BigInt {
	return it.Gas
}

func (it JianQianTx) GetGasPayer() qbasetypes.Address {
	return it.Address[0]
}

func (it JianQianTx) GetSignData() (ret []byte) {
	ret = append(ret, []byte(it.FuncName)...)
	for _, v := range it.Address {
		ret = append(ret, v.Bytes()...)
	}
	for _, v := range it.Args {
		ret = append(ret, []byte(v)...)
	}
	return
}
