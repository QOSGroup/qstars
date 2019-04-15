package coins

import (
	"github.com/QOSGroup/qbase/context"
	ctx "github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/baseapp"
	"github.com/QOSGroup/qstars/x/jianqian"
	go_amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
)

const COINNAME = "AOE"

type CoinsStub struct {
}

func NewCoinsStub() CoinsStub {
	return CoinsStub{}
}

func (cstub CoinsStub) StartX(base *baseapp.QstarsBaseApp) error {
	var coinsMapper = jianqian.NewCoinsMapper(jianqian.CoinsMapperName)
	base.Baseapp.RegisterMapper(coinsMapper)
	var aoeMapper = jianqian.NewAccountMapper(jianqian.AoeAccountMapperName)
	base.Baseapp.RegisterMapper(aoeMapper)
	return nil
}
func (cstub CoinsStub) EndBlockNotify(ctx context.Context) {

}

func (cstub CoinsStub) RegisterCdc(cdc *go_amino.Codec) {
	cdc.RegisterConcrete(&CoinsTx{}, "jianqian/CoinsTx", nil)
	cdc.RegisterConcrete(&AOETx{}, "jianqian/AOETx", nil)
}

func (cstub CoinsStub) ResultNotify(ctx context.Context, txQcpResult interface{}) *types.Result {
	//in := txQcpResult.(*txs.QcpTxResult)
	//log.Printf("ResultNotify QcpOriginalSequence:%s, result:%+v", string(in.QcpOriginalSequence), txQcpResult)
	//result := &types.Result{}
	//result.Code = types.CodeOK
	//qcpTxResult, ok := baseabci.ConvertTxQcpResult(txQcpResult)
	//if ok == false {
	//	log.Printf("coins.CoinsStub ResultNotify ConvertTxQcpResult error.")
	//	return result
	//} else {
	//	log.Printf("ResultNotify update status")
	//	orginalTxHash := in.QcpOriginalExtends //orginalTx.abc
	//	kvMapper := ctx.Mapper(common.QSCResultMapperName).(*common.KvMapper)
	//	initValue := ""
	//	kvMapper.Get([]byte(orginalTxHash), &initValue)
	//	if initValue != cstub.Name() {
	//		log.Printf("This is not my response.")
	//		return result
	//	}
	//	//put result to map for client query
	//	c := strconv.FormatInt((int64)(qcpTxResult.Result.Code), 10)
	//	c = c + " " + qcpTxResult.Result.Log
	//	log.Printf("--------update key:" + common.QSCResultMapperName + " key:" + orginalTxHash + " value:" + c)
	//	kvMapper.Set([]byte(orginalTxHash), c)
	//	//根据跨链结果 更新记录结果
	//	coinsMapper := ctx.Mapper(jianqian.CoinsMapperName).(*jianqian.CoinsMapper)
	//	coinsMapper.UpdateCoins(ctx.TxBytes(), c)
	//}

	return nil
}

func (cstub CoinsStub) CustomerQuery(ctx ctx.Context, route []string, req abci.RequestQuery) (res []byte, err types.Error) {
	return nil, nil
}

func (cstub CoinsStub) Name() string {
	return "CoinsStub"
}
