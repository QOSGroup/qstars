package comm

import (
	"errors"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	qbasetypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/types"

	"github.com/QOSGroup/qstars/x/jianqian"
	"strconv"
	"strings"
	"time"
)

type ArticleTx struct {
	AuthorAddr          qbasetypes.Address //作者地址(必填) 0qos 1cosmos
	ArticleType         int                //是否原创 0原创 1转载
	ArticleHash         string             //作品唯一标识hash
	ShareAuthor         int                //作者收入比例(必填)
	ShareOriginalAuthor int                //原创收入比例(转载作品必填)
	ShareCommunity      int                //社区收入比例(必填)
	ShareInvestor       int                //投资者收入比例(必填)
	InvestHours         int                //可供投资的小时数(必填)
	BuyHours            int                //可供购买广告位的小时数(必填)
	CoinType            string             //币种
}

var _ RouterTx = (*ArticleTx)(nil)

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
func (tx *ArticleTx) Exec(ctx context.Context) (result qbasetypes.Result, crossTxQcp *txs.TxQcp) {
	//本地存储
	articleMapper := ctx.Mapper(jianqian.ArticlesMapperName).(*jianqian.ArticlesMapper)

	buyhours := ctx.BlockHeader().Time.Add(time.Hour * (time.Duration(tx.BuyHours)))
	investhours := ctx.BlockHeader().Time.Add(time.Hour * (time.Duration(tx.InvestHours)))

	art := jianqian.Articles{tx.AuthorAddr, tx.ArticleType, tx.ArticleHash, tx.ShareAuthor, tx.ShareOriginalAuthor,
		tx.ShareCommunity, tx.ShareInvestor, tx.InvestHours, investhours, tx.BuyHours, buyhours, tx.CoinType, qbasetypes.NewInt(0)}

	if !articleMapper.SetArticle(tx.ArticleHash, &art) {
		result.Log = "Error: Save Article  error"
		result = qbasetypes.ErrInternal(result.Log).Result()
	}

	return
}

func (tx *ArticleTx) NewTx(args []string) error {
	args_len := len(args)
	if args_len != para_len_10 {
		return errors.New("AdvertisersTx args len error want " + strconv.Itoa(para_len_10) + " got " + strconv.Itoa(args_len))
	}

	address, err := types.AccAddressFromBech32(args[0])
	if err != nil {
		return err
	}

	articleType, err := strconv.Atoi(args[1])
	if err != nil {
		return err
	}

	shareAuthor, err := strconv.Atoi(args[3])
	if err != nil {
		return err
	}
	shareOriginalAuthor, err := strconv.Atoi(args[4])
	if err != nil {
		return err
	}
	shareCommunity, err := strconv.Atoi(args[5])
	if err != nil {
		return err
	}

	shareInvestor, err := strconv.Atoi(args[6])
	if err != nil {
		return err
	}
	investHours, err := strconv.Atoi(args[7])
	if err != nil {
		return err
	}
	buyHours, err := strconv.Atoi(args[8])
	if err != nil {
		return err
	}

	tx.AuthorAddr = address
	tx.ArticleType = articleType
	tx.ArticleHash = args[2]
	tx.ShareAuthor = shareAuthor
	tx.ShareOriginalAuthor = shareOriginalAuthor
	tx.ShareCommunity = shareCommunity
	tx.ShareInvestor = shareInvestor
	tx.InvestHours = investHours
	tx.BuyHours = buyHours
	tx.CoinType = args[9]
	return nil

}
