// Copyright 2018 The QOS Authors

package buyad

import (
	"errors"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	qbasetypes "github.com/QOSGroup/qbase/types"
	qostxs "github.com/QOSGroup/qos/txs/transfer"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/utility"
	"github.com/QOSGroup/qstars/x/jianqian"
	"log"
	"strings"
	"time"
)

type BuyTx struct {
	Addreess    qbasetypes.Address
	ArticleHash []byte `json:"articleHash"` // 文章hash
}

var _ txs.ITx = (*BuyTx)(nil)

func (it BuyTx) ValidateData(ctx context.Context) error {
	if err := check(ctx, it.ArticleHash); err != nil {
		return err
	}
	if strings.TrimSpace(it.Addreess.String()) == "" {
		return errors.New("签名地址不能为空")
	}
	buyMapper := ctx.Mapper(jianqian.BuyMapperName).(*jianqian.BuyMapper)
	buyer, ok := buyMapper.GetBuyer(it.ArticleHash)
	if ok && buyer.CheckStatus != jianqian.CheckStatusFail {
		return errors.New("已被购买")
	}
	auctioMapper := ctx.Mapper(jianqian.AuctionMapperName).(*jianqian.AuctionMapper)
	_, exist := auctioMapper.GetAuction(it.ArticleHash)
	if !exist {
		return errors.New("没有竞拍者")
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

func check(ctx context.Context, articleKey []byte) error {
	articleMapper := ctx.Mapper(jianqian.ArticlesMapperName).(*jianqian.ArticlesMapper)
	a := articleMapper.GetArticle(string(articleKey))
	if a == nil {
		return errors.New("invalid article")
	}

	log.Printf("--- checkArticle: EndBuyDate:%+v, blockheader:%+v", a.EndBuyDate, ctx.BlockHeader())
	if err := checkArticleBase(a, ctx.BlockHeader().Time); err != nil {
		return err
	}

	return nil
}

func checkRevenue(ctx context.Context, articleKey []byte, totalAmount qbasetypes.BigInt, items []qostxs.TransItem) error {
	articleMapper := ctx.Mapper(jianqian.ArticlesMapperName).(*jianqian.ArticlesMapper)
	a := articleMapper.GetArticle(string(articleKey))
	if a == nil {
		return errors.New("invalid article")
	}
	investMapper := ctx.Mapper(jianqian.InvestMapperName).(*jianqian.InvestMapper)
	investors := investMapper.AllInvestors(articleKey)

	communityPri := config.GetServerConf().Community
	if communityPri == "" {
		return errors.New("no community")
	}

	_, addrben32, _ := utility.PubAddrRetrievalFromAmino(communityPri, articleMapper.GetCodec())
	communityAddr, err := types.AccAddressFromBech32(addrben32)
	if err != nil {
		return err
	}

	receivers := warpperReceivers(articleMapper.GetCodec(), a, totalAmount, investors, communityAddr)

	if len(receivers) != len(items) {
		return errors.New("invalid Receivers")
	}

	receiverMap := make(map[string]qostxs.TransItem)
	for _, v := range receivers {
		if vv, ok := receiverMap[v.Address.String()]; ok {
			vv.QOS = vv.QOS.Add(v.QOS)
			receiverMap[v.Address.String()] = vv
		} else {
			receiverMap[v.Address.String()] = v
		}
	}

	for _, v := range items {
		if vv, ok := receiverMap[v.Address.String()]; ok {
			if !v.QOS.Equal(vv.QOS) {
				return errors.New("qos invalid")
			}

		} else {
			return errors.New("invest not found")
		}
	}

	return nil
}

func (it BuyTx) Exec(ctx context.Context) (result qbasetypes.Result, crossTxQcps *txs.TxQcp) {
	buyMapper := ctx.Mapper(jianqian.BuyMapperName).(*jianqian.BuyMapper)
	articleMapper := ctx.Mapper(jianqian.ArticlesMapperName).(*jianqian.ArticlesMapper)
	investMapper := ctx.Mapper(jianqian.InvestMapperName).(*jianqian.InvestMapper)
	auctionMapper := ctx.Mapper(jianqian.AuctionMapperName).(*jianqian.AuctionMapper)
	auction, ok := auctionMapper.GetMaxAuction(it.ArticleHash)
	if !ok {
		return
	}

	article := articleMapper.GetArticle(string(it.ArticleHash))
	if article == nil {
		log.Printf("buyad.BuyadStub GetArticle error.")
		return
	}

	communityPri := config.GetServerConf().Community
	if communityPri == "" {
		return
	}

	_, addrben32, _ := utility.PubAddrRetrievalFromAmino(communityPri, articleMapper.GetCodec())
	communityAddr, err := types.AccAddressFromBech32(addrben32)
	if err != nil {
		return
	}

	investors := investMapper.AllInvestors(it.ArticleHash)
	investors = calculateRevenue(auctionMapper.GetCodec(), article, auction.Amount, investors, communityAddr)

	for _, v := range investors {
		investMapper.SetInvestor(jianqian.GetInvestKey(it.ArticleHash, v.OtherAddr, v.InvestorType), v)
	}
	buyer := &jianqian.Buyer{}
	buyer.ArticleHash = it.ArticleHash
	buyer.Address = auction.Address
	buyer.Buy = auction.Amount
	buyer.BuyTime = ctx.BlockHeader().Time
	buyer.CheckStatus = jianqian.CheckStatusSuccess
	buyMapper.SetBuyer(it.ArticleHash, *buyer)
	return
}

func (it BuyTx) GetSigner() []qbasetypes.Address {
	return []qbasetypes.Address{it.Addreess}
}

func (it BuyTx) CalcGas() qbasetypes.BigInt {
	return qbasetypes.ZeroInt()
}

func (it BuyTx) GetGasPayer() qbasetypes.Address {
	return it.Addreess
}

func (it BuyTx) GetSignData() (ret []byte) {
	ret = append(ret, it.ArticleHash...)
	ret = append(ret, it.Addreess.Bytes()...)

	return it.ArticleHash
}

func (it BuyTx) Name() string {
	return "BuyTx"
}
