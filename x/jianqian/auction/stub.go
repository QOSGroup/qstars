package auction

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

type AuctionStub struct {
}

func NewAuctionStub() AuctionStub {
	return AuctionStub{}
}

func (cstub AuctionStub) StartX(base *baseapp.QstarsBaseApp) error {
	var auctionMapper = jianqian.NewAuctionMapper(jianqian.AuctionMapperName)
	base.Baseapp.RegisterMapper(auctionMapper)
	return nil
}
func (cstub AuctionStub) EndBlockNotify(ctx context.Context) {

}

func (cstub AuctionStub) RegisterCdc(cdc *go_amino.Codec) {
	cdc.RegisterConcrete(&AuctionTx{}, "jianqian/AuctionTx", nil)
}

func (cstub AuctionStub) ResultNotify(ctx context.Context, txQcpResult interface{}) *types.Result {
	//in := txQcpResult.(*txs.QcpTxResult)
	//log.Printf("ResultNotify QcpOriginalSequence:%s, result:%+v", string(in.QcpOriginalSequence), txQcpResult)
	//result := &types.Result{}
	//result.Code = types.
	//qcpTxResult, ok := baseabci.ConvertTxQcpResult(txQcpResult)
	//if ok == false {
	//	log.Printf("auction.AuctionStub ResultNotify ConvertTxQcpResult error.")
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
	//	in := txQcpResult.(*txs.QcpTxResult)
	//	key:=in.QcpOriginalExtends
	//	auctionMapper := ctx.Mapper(jianqian.AuctionMapperName).(*jianqian.AuctionMapper)
	//	tmpacution, ok := auctionMapper.GetTempAuction(key)
	//	if !ok {
	//		log.Printf("acution.AuctionStub unexpected buyerSta.")
	//		return result
	//	}
	//
	//	auction, ok := auctionMapper.GetAuctionByAddress(tmpacution.Article,tmpacution.Address.String())
	//	if ok {
	//		auction.Amount = auction.Amount.Add(tmpacution.Amount)
	//
	//	}else{
	//		auction = jianqian.Auction{tmpacution.Article, tmpacution.Address, tmpacution.CoinsType,  tmpacution.Amount, tmpacution.AuctionTime}
	//
	//	}
	//	auctionMapper.SetAuction(auction)
	//	//删除临时状态
	//	auctionMapper.DeleteAuction([]byte(key))
	//}

	return nil
}

func (cstub AuctionStub) CustomerQuery(ctx ctx.Context, route []string, req abci.RequestQuery) (res []byte, err types.Error) {
	return nil, nil
}

func (cstub AuctionStub) Name() string {
	return "AuctionStub"
}
