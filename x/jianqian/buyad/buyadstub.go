package buyad

import (
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	ctx "github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	qbasetypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/baseapp"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/utility"
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

func (bs BuyadStub) ResultNotify(ctx context.Context, txQcpResult interface{}) *qbasetypes.Result {
	result := &qbasetypes.Result{}
	result.Code = qbasetypes.ABCICodeType(qbasetypes.CodeOK)

	log.Printf("buyad.BuyadStub ResultNotify")
	in := txQcpResult.(*txs.QcpTxResult)
	log.Printf("buyad.BuyadStub ResultNotify QcpOriginalSequence:%s, result:%+v", string(in.QcpOriginalSequence), txQcpResult)
	qcpTxResult, ok := baseabci.ConvertTxQcpResult(txQcpResult)
	if ok == false {
		log.Printf("buyad.BuyadStub ResultNotify ConvertTxQcpResult error.")
		return result
	}

	log.Printf("buyad.BuyadStub ResultNotify update status")
	key := in.QcpOriginalExtends //orginalTx.abc

	kvMapper := ctx.Mapper(common.QSCResultMapperName).(*common.KvMapper)
	initValue := ""
	kvMapper.Get([]byte(key), &initValue)
	if initValue != bs.Name() {
		log.Printf("buyad.BuyadStub This is not my response.")
		return result
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
		return result
	}

	buyer, ok := buyMapper.GetBuyer(buyerSta.ArticleHash)
	if !ok || buyer == nil {
		log.Printf("buyad.BuyadStub unexpected buyer.")
		return result
	}

	if buyerSta.CheckStatus != jianqian.CheckStatusInit {
		log.Printf("buyad.BuyadStub unexpected status.")
		return result
	}

	if buyer.CheckStatus != jianqian.CheckStatusInit {
		log.Printf("buyad.BuyadStub unexpected status.")
		return result
	}

	if !qcpTxResult.Result.IsOK() {
		buyerSta.CheckStatus = jianqian.CheckStatusFail
		log.Printf("buyad.BuyadStub buyerSta update key:%+v\n", key)
		buyMapper.SetBuyer(buyerSta.ArticleHash, *buyerSta)

		buyer.CheckStatus = jianqian.CheckStatusFail
		log.Printf("buyad.BuyadStub buyer update key:%+v\n", key)
		buyMapper.SetBuyer(buyerSta.ArticleHash, *buyer)

		return result
	}

	articleMapper := ctx.Mapper(jianqian.ArticlesMapperName).(*jianqian.ArticlesMapper)
	investMapper := ctx.Mapper(jianqian.InvestMapperName).(*jianqian.InvestMapper)

	article := articleMapper.GetArticle(string(buyerSta.ArticleHash))
	if article == nil {
		log.Printf("buyad.BuyadStub GetArticle error.")
		return result
	}

	communityPri := config.GetServerConf().Community
	if communityPri == "" {
		return result
	}

	_, addrben32, _ := utility.PubAddrRetrievalFromAmino(communityPri, articleMapper.GetCodec())
	communityAddr, err := types.AccAddressFromBech32(addrben32)
	if err != nil {
		return result
	}

	investors := investMapper.AllInvestors(buyerSta.ArticleHash)
	investors = calculateRevenue(buyMapper.GetCodec(), article, buyer.Buy, investors, communityAddr)

	for _, v := range investors {
		investMapper.SetInvestor(jianqian.GetInvestKey(buyer.ArticleHash, v.Address, v.InvestorType), v)
	}

	buyer.CheckStatus = jianqian.CheckStatusSuccess
	log.Printf("buyad.BuyadStub buyer update key:%+v\n", key)
	buyMapper.SetBuyer(buyerSta.ArticleHash, *buyer)

	// 删除临时状态
	buyMapper.DeleteBuyer([]byte(key))

	return result
}

func (bs BuyadStub) EndBlockNotify(ctx context.Context) {

}

func (kv BuyadStub) CustomerQuery(ctx ctx.Context, route []string, req abci.RequestQuery) (res []byte, err qbasetypes.Error) {
	return nil, nil
}

func (kv BuyadStub) Name() string {
	return "BuyadStub"
}
