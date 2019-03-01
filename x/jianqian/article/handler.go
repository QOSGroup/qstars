package article

import (
	"errors"
	"fmt"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/x/jianqian"
	"strings"
	"time"
)

type ArticleTx struct {
	AuthorAddr          types.Address  //作者地址(必填) 0qos 1cosmos
	AuthorOtherAddr     string         //作者其他帐户地址
	ArticleType         int            //是否原创 0原创 1转载
	ArticleHash         string         //作品唯一标识hash
	ShareAuthor         int            //作者收入比例(必填)
	ShareOriginalAuthor int            //原创收入比例(转载作品必填)
	ShareCommunity      int            //社区收入比例(必填)
	ShareInvestor       int            //投资者收入比例(必填)
	InvestHours         int           //可供投资的小时数(必填)
	BuyHours            int           //可供购买广告位的小时数(必填)
	CoinType            string        //币种
	Gas                 types.BigInt
}

func (tx *ArticleTx) ValidateData(ctx context.Context) error {

	if strings.TrimSpace(tx.ArticleHash) == "" {
		return errors.New("Article hash must not empty")
	}
	if strings.TrimSpace(tx.CoinType) == "" {
		return errors.New("Article cointype must not empty")
	}
	if strings.TrimSpace(tx.AuthorAddr.String()) == "" {
		return errors.New("Article Authoraddress must not empty")
	}
	if tx.ShareAuthor > 100 {
		return errors.New("Article ShareAuthor Cannot be greater than 100")
	}
	if tx.ShareOriginalAuthor > 100 {
		return errors.New("Article ShareOriginalAuthor Cannot be greater than 100")
	}
	if tx.ShareInvestor > 100 {
		return errors.New("Article ShareInvestor Cannot be greater than 100")
	}
	articleMapper := ctx.Mapper(jianqian.ArticlesMapperName).(*jianqian.ArticlesMapper)
	if articleMapper.GetArticle(tx.ArticleHash) != nil {
		return errors.New("Article already exist!")
	}

	return nil
}

//执行业务逻辑,
// crossTxQcp: 需要进行跨链处理的TxQcp。
// 业务端实现中crossTxQcp只需包含`to` 和 `txStd`
func (tx *ArticleTx) Exec(ctx context.Context) (result types.Result, crossTxQcp *txs.TxQcp) {
	//本地存储
	articleMapper := ctx.Mapper(jianqian.ArticlesMapperName).(*jianqian.ArticlesMapper)

	buyhours := ctx.BlockHeader().Time.Add(time.Hour * ( time.Duration(tx.BuyHours)))
	investhours := ctx.BlockHeader().Time.Add(time.Hour * ( time.Duration(tx.InvestHours)))

	art := jianqian.Articles{tx.AuthorAddr, tx.AuthorOtherAddr, tx.ArticleType, tx.ArticleHash,tx.ShareAuthor, tx.ShareOriginalAuthor,
		tx.ShareCommunity, tx.ShareInvestor, tx.InvestHours, investhours, tx.BuyHours, buyhours, tx.CoinType,tx.Gas}

	if !articleMapper.SetArticle(tx.ArticleHash, &art) {
		result.Log = "Error: Save Article  error"
		result = types.ErrInternal(result.Log).Result()
	}

	return
}

func (tx *ArticleTx) GetSigner() []types.Address {

	return []types.Address{tx.AuthorAddr}
}
func (tx *ArticleTx) CalcGas() types.BigInt {
	return tx.Gas
}
func (tx *ArticleTx) GetGasPayer() types.Address {

	return types.Address{}
}
func (tx *ArticleTx) GetSignData() (ret []byte) {

	ret = append(ret, tx.ArticleHash...)
	ret = append(ret, tx.AuthorAddr...)
	ret = append(ret, tx.AuthorOtherAddr...)
	ret = append(ret, tx.CoinType...)
	ret = append(ret, fmt.Sprint(tx.ShareInvestor)...)
	ret = append(ret, fmt.Sprint(tx.ShareOriginalAuthor)...)
	ret = append(ret, fmt.Sprint(tx.ShareCommunity)...)
	ret = append(ret, fmt.Sprint(tx.ShareAuthor)...)

	return ret
}

func (tx ArticleTx) Name() string {
	return "ArticleTx"
}

func NewArticlesTx(authoraddress types.Address,authorOtherAddr string, articleType int,articleHash string, shareAuthor, shareOriginalAuthor,
	shareCommunity, shareInvestor, endInvestDate, endBuyDate int, cointype string,gas types.BigInt) *ArticleTx {
	return &ArticleTx{authoraddress, authorOtherAddr,articleType, articleHash, shareAuthor, shareOriginalAuthor,
		shareCommunity, shareInvestor, endInvestDate, endBuyDate,cointype, gas}
}
