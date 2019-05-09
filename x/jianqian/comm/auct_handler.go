package comm

import (
	"errors"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	qstarstypes "github.com/QOSGroup/qstars/types"

	"github.com/QOSGroup/qstars/x/jianqian"
	"log"
	"strconv"
	"time"
)

//广告竞买
type AuctionTx struct {
	Address     types.Address //qos地址
	CoinAmount  types.BigInt  //数量
	ArticleHash string        // 文章hash
	CoinsType   string        //币种
}

var _ RouterTx = (*AuctionTx)(nil)

func (tx *AuctionTx) ValidateData(ctx context.Context) error {
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
func (tx *AuctionTx) Exec(ctx context.Context) (result types.Result, crossTxQcps *txs.TxQcp) {
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
	accMapper.SubtractBalance(tx.Address.String(), tx.CoinsType, tx.CoinAmount)
	return
}



func (tx *AuctionTx) NewTx(args []string) error {
	args_len := len(args)
	if args_len != para_len_4 {
		return errors.New("AdvertisersTx args len error want " + strconv.Itoa(para_len_4) + " got " + strconv.Itoa(args_len))
	}
	address, err := qstarstypes.AccAddressFromBech32(args[0])
	if err != nil {
		return err
	}
	coins, ok := types.NewIntFromString(args[1])
	if !ok {
		return errors.New("amount format error")
	}
	tx.ArticleHash = args[2]
	tx.Address = address
	tx.CoinsType = args[3]
	tx.CoinAmount = coins

	return nil
}
