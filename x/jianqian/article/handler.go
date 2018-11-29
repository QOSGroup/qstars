package article

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/x/jianqian"
	"errors"
	"strings"
	"time"
)


type ArticleTx struct {
	Authoraddress       types.Address   //作者地址(必填)
	OriginalAuthor      types.Address   //原创作者地址(为空表示原创)
	ArticleHash         string   //作品唯一标识hash
	ShareAuthor         int   //作者收入比例(必填)
	ShareOriginalAuthor int   //原创收入比例(转载作品必填)
	ShareCommunity      int   //社区收入比例(必填)
	ShareInvestor       int   //投资者收入比例(必填)
	InvestDays          int   //可供投资的天数(必填)
	BuyDays             int          //可供购买广告位的天数(必填)
	Gas                 types.BigInt
}

func (tx *ArticleTx) ValidateData(ctx context.Context) error {

	if strings.TrimSpace(tx.ArticleHash ) == "" {
		return errors.New("Article hash must not empty")
	}
	if strings.TrimSpace(tx.Authoraddress.String())=="" {
		return errors.New("Article Authoraddress must not empty")
	}
	if tx.ShareAuthor>100 {
		return errors.New("Article ShareAuthor Cannot be greater than 100")
	}
	if tx.ShareOriginalAuthor>100{
		return errors.New("Article ShareOriginalAuthor Cannot be greater than 100")
	}
	if tx.ShareInvestor>100 {
		return errors.New("Article ShareInvestor Cannot be greater than 100")
	}
	articleMapper := ctx.Mapper(jianqian.ArticlesMapperName).(*jianqian.ArticlesMapper)
	if articleMapper.GetArticle(tx.ArticleHash)!=nil{
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

	buydays:=ctx.BlockHeader().Time.Add(time.Hour*(24*time.Duration(tx.BuyDays)))
	investdays:=ctx.BlockHeader().Time.Add(time.Hour*(24*time.Duration(tx.InvestDays)))

	art:=jianqian.Articles{tx.Authoraddress,tx.OriginalAuthor,tx.ArticleHash,tx.ShareAuthor,tx.ShareOriginalAuthor,
		tx.ShareCommunity,tx.ShareInvestor,tx.InvestDays,investdays,tx.BuyDays,buydays,tx.Gas}

	if !articleMapper.SetArticle(tx.ArticleHash,&art){
		result.Log = "Error: Save Article  error"
		result = types.ErrInternal(result.Log).Result()
	}

	return
}

func (tx *ArticleTx) GetSigner() []types.Address {

	return []types.Address{tx.Authoraddress}
}
func (tx *ArticleTx) CalcGas() types.BigInt {
	return tx.Gas
}
func (tx *ArticleTx) GetGasPayer() types.Address {

	return types.Address{}
}
func (tx *ArticleTx) GetSignData() (ret []byte) {

	ret = append(ret, tx.ArticleHash...)
	ret = append(ret, tx.Authoraddress...)
	ret = append(ret, tx.OriginalAuthor...)
	ret = append(ret, fmt.Sprint(tx.ShareInvestor)...)
	ret = append(ret, fmt.Sprint(tx.ShareOriginalAuthor)...)
	ret = append(ret, fmt.Sprint(tx.ShareCommunity)...)
	ret = append(ret, fmt.Sprint(tx.ShareAuthor)...)

	return ret
}

func NewArticlesTx(authoraddress, originalAuthor types.Address,  articleHash string,  shareAuthor, shareOriginalAuthor,
shareCommunity,  shareInvestor , endInvestDate,  endBuyDate  int, gas types.BigInt) *ArticleTx {
	return &ArticleTx{authoraddress, originalAuthor,  articleHash,  shareAuthor, shareOriginalAuthor,
		shareCommunity,  shareInvestor,  endInvestDate,  endBuyDate, gas}
}
