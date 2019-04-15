package auction

import (
	"errors"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/x/jianqian"
	"log"
	"time"
)

//广告竞买
type AuctionTx struct {
	ArticleHash string        // 文章hash
	Address     types.Address //qos地址
	CoinsType   string        //币种
	CoinAmount types.BigInt //数量
	Gas        types.BigInt
}

func (tx AuctionTx) ValidateData(ctx context.Context) error {
	if tx.Address.Empty() {
		return errors.New("Auction address must not empty")
	}
	if tx.CoinsType == "" {
		return errors.New("Coins type must not empty")
	}
	articleMapper := ctx.Mapper(jianqian.ArticlesMapperName).(*jianqian.ArticlesMapper)
	a := articleMapper.GetArticle(string(tx.ArticleHash))
	if a == nil {
		return errors.New("竞拍作品不存在")
	}

	if err := checkArticleBase(a, ctx.BlockHeader().Time); err != nil {
		return err
	}

	if a != nil {
		if a.CoinType != tx.CoinsType {
			return errors.New("竞拍币种与作品要求不符")
		}
	}

	accMapper := ctx.Mapper(jianqian.AoeAccountMapperName).(*jianqian.AoeAccountMapper)
	blance := accMapper.GetBalance(tx.Address.String(), tx.CoinsType)
	if tx.CoinAmount.GT(blance) {
		//余额不足
		return errors.New("Insufficient balance")
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
	auction, ok := auctionMapper.GetAuctionByAddress(tx.ArticleHash, tx.Address.String())
	if ok {
		auction.Amount = auction.Amount.Add(tx.CoinAmount)
		auction.AuctionTime = ctx.BlockHeader().Time
	} else {
		auction = jianqian.Auction{tx.ArticleHash, tx.Address, tx.CoinsType, tx.CoinAmount, ctx.BlockHeader().Time}
	}
	auctionMapper.SetAuction(auction)
	accMapper := ctx.Mapper(jianqian.AoeAccountMapperName).(*jianqian.AoeAccountMapper)
	accMapper.SubtractBalance(tx.Address.String(),tx.CoinsType,tx.CoinAmount)
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
	ret = append(ret, []byte(tx.Address)...)
	ret = append(ret, tx.CoinsType...)
	ret = append(ret, types.Int2Byte(tx.CoinAmount.Int64())...)
	return ret
}

func (tx AuctionTx) Name() string {
	return "AuctionTx"
}

func NewAuctionTx(articleHash string, address types.Address, CoinsType string, OtherAddr string, coinAmount types.BigInt, gas types.BigInt) AuctionTx {

	return AuctionTx{articleHash, address, CoinsType, coinAmount, gas}
}
