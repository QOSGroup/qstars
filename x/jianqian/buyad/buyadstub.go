package buyad

import (
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
	"log"
	"strconv"
)

type BuyadStub struct {
}

func NewStub() BuyadStub {
	return BuyadStub{}
}

func (bs BuyadStub) StartX(base *baseapp.QstarsBaseApp) error {
	var buyMapper = jianqian.NewBuyMapper(jianqian.BuyMapperName)
	base.Baseapp.RegisterMapper(buyMapper)

	return nil
}

func (bs BuyadStub) RegisterCdc(cdc *go_amino.Codec) {
	cdc.RegisterConcrete(&BuyTx{}, "qstars/BuyTx", nil)
}

func (bs BuyadStub) ResultNotify(ctx context.Context, txQcpResult interface{}) *types.Result {
	result := &types.Result{}
	log.Printf("buyad.BuyadStub ResultNotify")
	in := txQcpResult.(*txs.QcpTxResult)
	log.Printf("buyad.BuyadStub ResultNotify QcpOriginalSequence:%s, result:%+v", string(in.QcpOriginalSequence), txQcpResult)
	qcpTxResult, ok := baseabci.ConvertTxQcpResult(txQcpResult)
	if ok == false {
		log.Printf("buyad.BuyadStub ResultNotify ConvertTxQcpResult error.")
		result.Code = types.ABCICodeType(types.CodeTxDecode)
		return result
	}

	log.Printf("buyad.BuyadStub ResultNotify update status")

	key := in.QcpOriginalExtends //orginalTx.abc

	kvMapper := ctx.Mapper(common.QSCResultMapperName).(*common.KvMapper)
	initValue := ""
	kvMapper.Get([]byte(key), &initValue)
	if initValue != bs.Name() {
		log.Printf("buyad.BuyadStub This is not my response.")
		return nil
	}
	log.Printf("buyad.BuyadStub ResultNotify kvMapper get key:%s, value:%s", key, initValue)
	c := strconv.FormatInt((int64)(qcpTxResult.Result.Code), 10)
	c = c + " " + qcpTxResult.Result.Log
	log.Printf("buyad.BuyadStub ResultNotify kvMapper Set key:%s, value:%s", key, c)
	kvMapper.Set([]byte(key), c)

	buyMapper := ctx.Mapper(jianqian.BuyMapperName).(*jianqian.BuyMapper)
	buyerSta, ok := buyMapper.GetBuyer([]byte(key))
	if !ok || buyerSta == nil {
		log.Printf("buyad.BuyadStub unexpected buyerSta.")
		result.Code = types.ABCICodeType(types.CodeTxDecode)
		return result
	}

	if buyerSta.CheckStatus != jianqian.CheckStatusInit {
		log.Printf("buyad.BuyadStub unexpected status.")
		result.Code = types.ABCICodeType(types.CodeTxDecode)
		return result
	}

	buyer, ok := buyMapper.GetBuyer(buyerSta.ArticleHash)
	if !ok || buyer == nil {
		log.Printf("buyad.BuyadStub unexpected buyer.")
		result.Code = types.ABCICodeType(types.CodeTxDecode)
		return result
	}

	if buyer.CheckStatus != jianqian.CheckStatusInit {
		log.Printf("buyad.BuyadStub unexpected status.")
		result.Code = types.ABCICodeType(types.CodeTxDecode)
		return result
	}

	if qcpTxResult.Result.IsOK() {
		buyer.CheckStatus = jianqian.CheckStatusSuccess
		log.Printf("buyad.BuyadStub buyer update key:%+v\n", key)
		buyMapper.SetBuyer(buyerSta.ArticleHash, *buyer)
	} else {
		log.Printf("buyad.BuyadStub buyer delete key:%+v", key)

		buyMapper.DeleteBuyer(buyerSta.ArticleHash)
	}

	// 删除临时状态
	buyMapper.DeleteBuyer([]byte(key))

	result.Code = qcpTxResult.Result.Code

	return result
}

func (bs BuyadStub) EndBlockNotify(ctx context.Context) {

}

func (kv BuyadStub) CustomerQuery(ctx ctx.Context, route []string, req abci.RequestQuery) (res []byte, err types.Error) {
	return nil, nil
}

func (kv BuyadStub) Name() string {
	return "BuyadStub"
}
