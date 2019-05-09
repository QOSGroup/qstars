// Copyright 2018 The QOS Authors

package comm

import (
	"errors"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	qbasetypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/utility"
	"github.com/QOSGroup/qstars/wire"
	"github.com/QOSGroup/qstars/x/jianqian"
	"log"
	"strconv"
	"strings"
)

type BuyTx struct {
	Addreess    qbasetypes.Address
	ArticleHash []byte `json:"articleHash"` // 文章hash
}

var _ RouterTx = (*BuyTx)(nil)

func (it *BuyTx) ValidateData(ctx context.Context) error {
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
	_, exist := auctioMapper.GetAuction(string(it.ArticleHash))
	if !exist {
		return errors.New("没有竞拍者")
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

func (it *BuyTx) Exec(ctx context.Context) (result qbasetypes.Result, crossTxQcps *txs.TxQcp) {
	buyMapper := ctx.Mapper(jianqian.BuyMapperName).(*jianqian.BuyMapper)
	articleMapper := ctx.Mapper(jianqian.ArticlesMapperName).(*jianqian.ArticlesMapper)
	investMapper := ctx.Mapper(jianqian.InvestMapperName).(*jianqian.InvestMapper)
	auctionMapper := ctx.Mapper(jianqian.AuctionMapperName).(*jianqian.AuctionMapper)
	auction, ok := auctionMapper.GetMaxAuction(string(it.ArticleHash))
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
	investors = calculateRevenue(buyMapper.GetCodec(), article, auction.Amount, investors, communityAddr)

	for _, v := range investors {
		investMapper.SetInvestor(jianqian.GetInvestKey(it.ArticleHash, v.Address.String(), v.InvestorType), v)
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
func (tx *BuyTx) NewTx(args []string) error {
	args_len := len(args)
	if args_len != para_len_2 {
		return errors.New("AdvertisersTx args len error want " + strconv.Itoa(para_len_2) + " got " + strconv.Itoa(args_len))
	}
	address, err := types.AccAddressFromBech32(args[0])
	if err != nil {
		return err
	}

	tx.Addreess = address
	tx.ArticleHash = []byte(args[1])
	return nil
}



// calculateRevenue 计算收入
func calculateRevenue(cdc *wire.Codec, article *jianqian.Articles, amount qbasetypes.BigInt, is jianqian.Investors,
	communityAddr qbasetypes.Address) jianqian.Investors {
	var result []jianqian.Investor
	log.Printf("buyad calculateRevenue  article:%+v, amount:%d", article, amount.Int64())


	//addrstr=article.AuthorAddr.String()
	//orgaddr=communityAddr.String()
	//communitystr=communityAddr.String()




	// 作者地址
	authorTotal := amount.Mul(qbasetypes.NewInt(int64(article.ShareAuthor))).Div(qbasetypes.NewInt(100))
	log.Printf("buyad calculateRevenue  Authoraddress:%s amount:%d", article.AuthorAddr.String(), authorTotal.Int64())
	result = append(
		result,
		jianqian.Investor{
			InvestorType: jianqian.InvestorTypeAuthor, // 投资者类型
			Address:      article.AuthorAddr,             // 投资者地址
			Invest:       qbasetypes.NewInt(0),        // 投资金额
			Revenue:      authorTotal,                 // 投资收益
		})

	// 原创作者地址
	shareOriginalTotal := amount.Mul(qbasetypes.NewInt(int64(article.ShareOriginalAuthor))).Div(qbasetypes.NewInt(100))
	log.Printf("buyad calculateRevenue  OriginalAuthor:%s amount:%d", communityAddr.String(), shareOriginalTotal.Int64())
	result = append(
		result,
		jianqian.Investor{
			InvestorType: jianqian.InvestorTypeOriginalAuthor, // 投资者类型
			Address:      communityAddr,              // 投资者地址
			Invest:       qbasetypes.NewInt(0),                // 投资金额
			Revenue:      shareOriginalTotal,                  // 投资收益
		})

	// 投资者收入分配
	investorShouldTotal := amount.Mul(qbasetypes.NewInt(int64(article.ShareInvestor))).Div(qbasetypes.NewInt(100))
	log.Printf("buyad calculateRevenue investorShouldTotal:%d", investorShouldTotal.Int64())
	investors := calculateInvestorRevenue(cdc, is, investorShouldTotal)
	result = append(result, investors...)

	shareCommunityTotal := amount.Sub(authorTotal).Sub(shareOriginalTotal).Sub(investors.TotalRevenue())
	log.Printf("buyad calculateRevenue  communityAddr:%s amount:%d", communityAddr.String(), shareCommunityTotal.Int64())
	// 社区收入比例
	result = append(
		result,
		jianqian.Investor{
			InvestorType: jianqian.InvestorTypeCommunity, // 投资者类型
			Address:      communityAddr,                  // 投资者地址
			Invest:       qbasetypes.NewInt(0),           // 投资金额
			Revenue:      shareCommunityTotal,            // 投资收益
		})

	return result
}


// calculateInvestorRevenue 计算投资者收入
func calculateInvestorRevenue(cdc *wire.Codec, investors jianqian.Investors, amount qbasetypes.BigInt) jianqian.Investors {
	log.Printf("buyAd calculateInvestorRevenue investors:%+v", investors)

	totalInvest := investors.TotalInvest()
	log.Printf("buyAd calculateInvestorRevenue amount:%s, totalInvest:%d", amount.String(), totalInvest.Int64())

	curAmount := qbasetypes.NewInt(0)
	if !totalInvest.IsZero() {
		l := len(investors)
		for i := 0; i < l; i++ {
			var revenue qbasetypes.BigInt
			if i+1 == l {
				revenue = amount.Sub(curAmount)
			} else {
				revenue = amount.Mul(investors[i].Invest).Div(totalInvest)
			}

			investors[i].Revenue = revenue
			curAmount = curAmount.Add(revenue)
			log.Printf("buyad calculateRevenue  investorAddr:%s invest:%d, revenue:%d",
				investors[i].Address.String(), investors[i].Invest.Int64(), revenue.Int64())
		}
	}

	return investors
}
