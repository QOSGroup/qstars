// Copyright 2018 The QOS Authors

package buyad

import (
	"encoding/json"
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/txs"
	qbasetypes "github.com/QOSGroup/qbase/types"
	qostxs "github.com/QOSGroup/qos/txs/transfer"
	qostypes "github.com/QOSGroup/qos/types"
	"github.com/QOSGroup/qstars/client/utils"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/utility"
	"github.com/QOSGroup/qstars/wire"
	"github.com/QOSGroup/qstars/x/common"
	"github.com/QOSGroup/qstars/x/jianqian"
	"github.com/QOSGroup/qstars/x/jianqian/tx"
	"log"
	"strings"
)

const coinsName = "QOS"



// BuyAd 投资广告
func BuyAd(cdc *wire.Codec, articleHash string) string {
	privatekey := tx.GetConfig().Adbuyermock
	var result common.Result
	result.Code = common.ResultCodeSuccess
	chainId := config.GetCLIContext().Config.QSCChainID
	tx, _ := buyAd(cdc, chainId, articleHash,  privatekey)
	cliCtx := *config.GetCLIContext().QSCCliContext
	_, commitresult, err1 := utils.SendTx(cliCtx, cdc, tx)
	if err1 != nil {
		return common.NewErrorResult(SENDTXERRCode, 0, "", err1.Error()).Marshal()
	}
	return common.NewSuccessResult(cdc, commitresult.Height, commitresult.Hash.String(), "").Marshal()
}


func mergeQSCs(q1, q2 qostypes.QSCs) qostypes.QSCs {
	m := make(map[string]*qbasetypes.BaseCoin)

	for _, v := range q1 {
		m[v.Name] = v
	}

	var res qostypes.QSCs
	for _, v := range q2 {
		if q, ok := m[v.Name]; ok {
			v.Amount.Add(q.Amount)
			m[v.Name] = v
		} else {
			m[v.Name] = v
		}
	}

	for _, v := range m {
		res = append(res, v)
	}

	return res
}

func mergeReceivers(rs []qostxs.TransItem) []qostxs.TransItem {
	var res []qostxs.TransItem
	m := make(map[string]qostxs.TransItem)

	for _, v := range rs {
		if ti, ok := m[v.Address.String()]; ok {
			v.QOS = v.QOS.Add(ti.QOS)
			v.QSCs = mergeQSCs(v.QSCs, ti.QSCs)
			m[v.Address.String()] = v
		} else {
			m[v.Address.String()] = v
		}
	}

	for _, v := range m {
		res = append(res, v)
	}

	log.Printf("buyad.mergeReceivers rs:%+v, res:%+v", rs, res)
	return res
}

func warpperReceivers(cdc *wire.Codec, article *jianqian.Articles, amount qbasetypes.BigInt,
	investors jianqian.Investors, communityAddr qbasetypes.Address) []qostxs.TransItem {
	var result []qostxs.TransItem
	log.Printf("buyad warpperReceivers  article:%+v", article)

	investors = calculateRevenue(cdc, article, amount, investors, communityAddr)

	for _, v := range investors {
		if !v.Revenue.IsZero() {
			addres,_:=types.AccAddressFromBech32(v.OtherAddr)
			result = append(result,	warpperTransItem(addres,[]qbasetypes.BaseCoin{{Name: coinsName, Amount: v.Revenue}}))
		}
	}

	return mergeReceivers(result)
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
				investors[i].OtherAddr, investors[i].Invest.Int64(), revenue.Int64())
		}
	}

	return investors
}

// calculateRevenue 计算收入
func calculateRevenue(cdc *wire.Codec, article *jianqian.Articles, amount qbasetypes.BigInt, is jianqian.Investors,
	communityAddr qbasetypes.Address) jianqian.Investors {
	var result []jianqian.Investor
	log.Printf("buyad calculateRevenue  article:%+v, amount:%d", article, amount.Int64())


	var addrstr,orgaddr,communitystr string
	if article.CoinType=="QOS"{
		addrstr=article.AuthorAddr.String()
		orgaddr=communityAddr.String()
		communitystr=communityAddr.String()
	}else{
		addrstr=article.AuthorOtherAddr
		orgaddr=config.GetCLIContext().Config.OrgAuthor_Other
		communitystr=config.GetCLIContext().Config.Community_Other
	}



	// 作者地址
	authorTotal := amount.Mul(qbasetypes.NewInt(int64(article.ShareAuthor))).Div(qbasetypes.NewInt(100))
	log.Printf("buyad calculateRevenue  Authoraddress:%s amount:%d", article.AuthorAddr.String(), authorTotal.Int64())
	result = append(
		result,
		jianqian.Investor{
			InvestorType: jianqian.InvestorTypeAuthor, // 投资者类型
			OtherAddr:      addrstr,             // 投资者地址
			Invest:       qbasetypes.NewInt(0),        // 投资金额
			Revenue:      authorTotal,                 // 投资收益
		})

	// 原创作者地址
	shareOriginalTotal := amount.Mul(qbasetypes.NewInt(int64(article.ShareOriginalAuthor))).Div(qbasetypes.NewInt(100))
	log.Printf("buyad calculateRevenue  OriginalAuthor:%s amount:%d", orgaddr, shareOriginalTotal.Int64())
	result = append(
		result,
		jianqian.Investor{
			InvestorType: jianqian.InvestorTypeOriginalAuthor, // 投资者类型
			OtherAddr:      orgaddr,              // 投资者地址
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
			OtherAddr:      communitystr,                  // 投资者地址
			Invest:       qbasetypes.NewInt(0),           // 投资金额
			Revenue:      shareCommunityTotal,            // 投资收益
		})

	return result
}

// buyAd 投资广告
func buyAd(cdc *wire.Codec, chainId, articleHash,  privatekey string) (*txs.TxStd, *BuyadErr) {
	_, addrben32, priv := utility.PubAddrRetrievalFromAmino(privatekey, cdc)

	from, err := types.AccAddressFromBech32(addrben32)

	key := account.AddressStoreKey(from)

	var nonce int64 = 0
	acc, err := config.GetCLIContext().QSCCliContext.GetAccount(key, cdc)
	if err != nil {
		nonce = 0
	} else {
		nonce = int64(acc.Nonce)
	}
	nonce++
	it := &BuyTx{}
	it.ArticleHash = []byte(articleHash)
	tx2 := txs.NewTxStd(it, config.GetCLIContext().Config.QSCChainID, qbasetypes.ZeroInt())
	signature2, _ := tx2.SignTx(priv, nonce, config.GetCLIContext().Config.QSCChainID, config.GetCLIContext().Config.QSCChainID)
	tx2.Signature = []txs.Signature{txs.Signature{
		Pubkey:    priv.PubKey(),
		Signature: signature2,
		Nonce:     nonce,
	}}
	return tx2, nil
}

func warpperTransItem(addr qbasetypes.Address, coins []qbasetypes.BaseCoin) qostxs.TransItem {
	var ti qostxs.TransItem
	ti.Address = addr
	ti.QOS = qbasetypes.NewInt(0)

	for _, coin := range coins {
		if strings.ToUpper(coin.Name) == "QOS" {
			ti.QOS = ti.QOS.Add(coin.Amount)
		} else {
			ti.QSCs = append(ti.QSCs, &coin)
		}
	}

	return ti
}

// RetrieveBuyer 查询购买者
func RetrieveBuyer(cdc *wire.Codec, articleHash string) string {
	var result common.Result
	result.Code = common.ResultCodeSuccess

	buyer, err := jianqian.QueryArticleBuyer(cdc, config.GetCLIContext().QSCCliContext, articleHash)
	if err != nil {
		log.Printf("QueryArticleBuyer err:%s", err.Error())
		result.Code = common.ResultCodeInternalError
		result.Reason = err.Error()
		return result.Marshal()
	}

	js, err := cdc.MarshalJSON(buyer)
	if err != nil {
		log.Printf("buyAd err:%s", err.Error())
		result.Code = common.ResultCodeInternalError
		result.Reason = err.Error()
		return result.Marshal()
	}
	result.Result = json.RawMessage(js)

	return result.Marshal()
}
