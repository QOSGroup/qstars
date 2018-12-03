package buyad

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
	go_amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
)

type BuyadStub struct {
	baseapp.BaseXTransaction
}

func NewStub() BuyadStub {
	return BuyadStub{}
}

func (bs BuyadStub) StartX(base *baseapp.QstarsBaseApp) error {
	var qosMapper = common.NewKvMapper(jianqian.BuyMapperName)
	base.Baseapp.RegisterMapper(qosMapper)

	return nil
}

func (bs BuyadStub) RegisterCdc(cdc *go_amino.Codec) {
	cdc.RegisterConcrete(&BuyTx{}, "qstars/BuyTx", nil)
}

func (bs BuyadStub) ResultNotify(ctx context.Context, txQcpResult interface{}) *types.Result {
	result := &types.Result{}

	in := txQcpResult.(*txs.QcpTxResult)
	fmt.Printf("buyad.BuyadStub ResultNotify QcpOriginalSequence:%s, result:%+v", string(in.QcpOriginalSequence), txQcpResult)
	qcpTxResult, ok := baseabci.ConvertTxQcpResult(txQcpResult)
	if ok == false {
		fmt.Printf("ResultNotify ConvertTxQcpResult error.")
		result.Code = types.ABCICodeType(types.CodeTxDecode)
		return result
	}

	fmt.Printf("buyad.BuyadStub ResultNotify update status")

	key := in.QcpOriginalExtends //orginalTx.abc
	buyMapper := ctx.Mapper(jianqian.BuyMapperName).(*jianqian.BuyMapper)
	buyer, ok := buyMapper.GetBuyer([]byte(key))
	if !ok || buyer == nil {

		fmt.Printf("unexpected buyer.")
		result.Code = types.ABCICodeType(types.CodeTxDecode)
		return result
	}

	if buyer.CheckStatus != jianqian.CheckStatusInit {
		fmt.Printf("unexpected status.")
		result.Code = types.ABCICodeType(types.CodeTxDecode)
		return result
	}

	if qcpTxResult.Result.IsOK() {
		buyer.CheckStatus = jianqian.CheckStatusSuccess
		buyMapper.SetBuyer([]byte(key), *buyer)
	} else {
		buyMapper.DeleteBuyer([]byte(key))
	}

	result.Code = qcpTxResult.Result.Code

	return result
}

func (bs BuyadStub) EndBlockNotify(ctx context.Context) {

}

func (kv BuyadStub) CustomerQuery(ctx ctx.Context, route []string, req abci.RequestQuery) (res []byte, err types.Error) {
	return nil, nil
}
