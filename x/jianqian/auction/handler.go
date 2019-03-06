package auction

import (
	"errors"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/x/common"
	"github.com/QOSGroup/qstars/x/jianqian"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"log"
	"time"

	tmcommon "github.com/tendermint/tendermint/libs/common"
	"strconv"
	"strings"
)

//广告竞买
type AuctionTx struct {
	Wrapper      *txs.TxStd     //已封装好的 TxTransfer 结构体
	ArticleHash  string         // 文章hash
	Address      types.Address  //qos地址
	CoinsType    string         //币种
	OtherAddr    string         //转出方地址
	CoinAmount   types.BigInt   //数量
	Gas          types.BigInt
}

func (tx AuctionTx) ValidateData(ctx context.Context) error {
	if tx.Address.Empty() {
		return errors.New("Auction address must not empty")
	}
	if tx.CoinsType=="" {
		return errors.New("Coins type must not empty")
	}
	articleMapper := ctx.Mapper(jianqian.ArticlesMapperName).(*jianqian.ArticlesMapper)
	a := articleMapper.GetArticle(string(tx.ArticleHash))
	if a==nil{
		return errors.New("竞拍作品不存在")
	}

	if err := checkArticleBase(a, ctx.BlockHeader().Time); err != nil {
		return err
	}

	if a!=nil {
		if a.CoinType != tx.CoinsType {
			return errors.New("竞拍币种与作品要求不符")
		}
	}
	return nil
}


func checkArticleBase(article *jianqian.Articles, now time.Time) error {
	log.Printf("checkArticleBase EndInvestDate:%+v, EndBuyDate:%+v, now:%+v", article.EndInvestDate, article.EndBuyDate, now)
	if article.EndInvestDate.After(now) {
		return errors.New("投资还没结束")
	}
	if article.EndBuyDate.Before(now) {
		return errors.New("超过购买期限")
	}

	return nil
}

//执行业务逻辑,
// crossTxQcp: 需要进行跨链处理的TxQcp。
// 业务端实现中crossTxQcp只需包含`to` 和 `txStd`
func (tx AuctionTx) Exec(ctx context.Context) (result types.Result, crossTxQcps *txs.TxQcp) {
	auctionMapper := ctx.Mapper(jianqian.AuctionMapperName).(*jianqian.AuctionMapper)

	//币种为qos才跨链
	if strings.ToUpper(tx.CoinsType)=="QOS"{
		tx1 := (tmcommon.HexBytes)(tmhash.Sum(ctx.TxBytes()))
		heigth1 := strconv.FormatInt(ctx.BlockHeight(), 10)
		key := "heigth:" + heigth1 + ",hash:" + tx1.String()

		kvMapper := ctx.Mapper(common.QSCResultMapperName).(*common.KvMapper)
		kvMapper.Set([]byte(key), AuctionStub{}.Name())

		//跨链
		crossTxQcps = &txs.TxQcp{}
		crossTxQcps.TxStd = tx.Wrapper
		crossTxQcps.To = config.GetServerConf().QOSChainName
		crossTxQcps.Extends=key
		result = types.Result{
			Code:  types.ABCICodeOK,
		}

		//设置临时状态
		auction:=&jianqian.Auction{tx.ArticleHash,tx.Address,tx.CoinsType,tx.OtherAddr,tx.CoinAmount,ctx.BlockHeader().Time}
		auctionMapper.SetAuctionByKey([]byte(key),auction)

	}else {
		//合并多次投资数额
		auction, ok := auctionMapper.GetAuctionByAddress(tx.ArticleHash, tx.Address.String())
		if ok {
			auction.Amount = auction.Amount.Add(tx.CoinAmount)
			auction.AuctionTime=ctx.BlockHeader().Time
		} else {
			auction=jianqian.Auction{tx.ArticleHash,tx.Address,tx.CoinsType,tx.OtherAddr,tx.CoinAmount,ctx.BlockHeader().Time}
		}
		auctionMapper.SetAuction(auction)
	}
	return
}

func (tx AuctionTx) GetSigner() []types.Address {
	return []types.Address{tx.Address}
}
func (tx AuctionTx) CalcGas() types.BigInt {
	return tx.Gas
}
func (tx AuctionTx) GetGasPayer() types.Address {
	return types.Address{}
}
func (tx AuctionTx) GetSignData() (ret []byte) {
	if tx.Wrapper!=nil {
		ret = append(ret, tx.Wrapper.ITx.GetSignData()...)
	}
	ret = append(ret, []byte(tx.Address)...)
	ret = append(ret, tx.CoinsType...)
	ret = append(ret, tx.OtherAddr...)
	ret = append(ret, types.Int2Byte(tx.CoinAmount.Int64())...)
	return ret
}

func (tx AuctionTx) Name() string {
	return "AuctionTx"
}

func NewAuctionTx(wrapper *txs.TxStd,articleHash string, address types.Address, CoinsType string, OtherAddr string,coinAmount types.BigInt, gas types.BigInt) AuctionTx {

	return AuctionTx{wrapper,articleHash, address, CoinsType,OtherAddr, coinAmount, gas}
}
