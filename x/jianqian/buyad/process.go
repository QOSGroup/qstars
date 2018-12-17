// Copyright 2018 The QOS Authors

package buyad

import (
	"encoding/json"
	"errors"
	"fmt"
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
	"log"
	"strconv"
	"strings"
	"time"
)

type ResultBuy struct {
	Code   string          `json:"code"`
	Reason string          `json:"reason,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
}

func InternalError(reason string) ResultBuy {
	return ResultBuy{Code: "-1", Reason: reason}
}

func NewResultBuy(cdc *wire.Codec, code, reason string, res interface{}) ResultBuy {
	var rawMsg json.RawMessage

	if res != nil {
		var js []byte
		js, err := cdc.MarshalJSON(res)
		if err != nil {
			return InternalError(err.Error())
		}
		rawMsg = json.RawMessage(js)
	}

	var result ResultBuy
	result.Result = rawMsg
	result.Code = code
	result.Reason = reason

	return result
}

func (ri ResultBuy) Marshal() string {
	jsonBytes, err := json.MarshalIndent(ri, "", "  ")
	if err != nil {
		log.Printf("BuyAd err:%s", err.Error())
		return InternalError(err.Error()).Marshal()
	}
	return string(jsonBytes)
}

const coinsName = "QOS"

// BuyAdBackground 提交到链上
func BuyAdBackground(cdc *wire.Codec, txb string, timeout time.Duration) string {
	ts := new(txs.TxStd)
	err := cdc.UnmarshalJSON([]byte(txb), ts)
	log.Printf("buyad.BuyAdBackground ts:%+v, err:%+v", ts, err)
	if err != nil {
		return InternalError(err.Error()).Marshal()
	}

	cliCtx := *config.GetCLIContext().QSCCliContext
	_, commitresult, err := utils.SendTx(cliCtx, cdc, ts)
	log.Printf("buyad.BuyAdBackground SendTx commitresult:%+v, err:%+v", commitresult, err)

	if err != nil {
		return InternalError(err.Error()).Marshal()
	}

	height := strconv.FormatInt(commitresult.Height, 10)
	var code, reason string
	var result interface{}

	waittime, err := strconv.Atoi(config.GetCLIContext().Config.WaitingForQosResult)
	if err != nil {
		panic("WaitingForQosResult should be able to convert to integer." + err.Error())
	}
	counter := 0

	for {
		if counter >= waittime {
			log.Println("time out")
			result = "time out"
			break
		}
		resultstr, err := fetchResult(cdc, height, commitresult.Hash.String())
		log.Printf("fetchResult result:%s, err:%+v\n", resultstr, err)
		if err != nil {
			log.Printf("fetchResult error:%s\n", err.Error())
			reason = err.Error()
			break
		}

		if resultstr != "" && resultstr != (BuyadStub{}).Name() {
			log.Printf("fetchResult result:[%+v]\n", resultstr)
			rs := []rune(resultstr)
			index1 := strings.Index(resultstr, " ")

			reason = ""
			result = string(rs[index1+1:])
			code = string(rs[:index1])
			break
		}
		time.Sleep(500 * time.Millisecond)
		counter++
	}

	return NewResultBuy(cdc, code, reason, result).Marshal()
}

func fetchResult(cdc *wire.Codec, heigth1 string, tx1 string) (string, error) {
	qstarskey := "heigth:" + heigth1 + ",hash:" + tx1
	d, err := config.GetCLIContext().QSCCliContext.QueryStore([]byte(qstarskey), common.QSCResultMapperName)
	if err != nil {
		return "", err
	}
	if d == nil {
		return "", nil
	}
	var res []byte
	err = cdc.UnmarshalBinaryBare(d, &res)
	if err != nil {
		return "", err
	}
	return string(res), err
}

// BuyAd 投资广告
func BuyAd(cdc *wire.Codec, chainId, articleHash, coins, privatekey string, qosnonce, qscnonce int64) string {
	var result ResultBuy
	result.Code = "0"

	tx, err := buyAd(cdc, chainId, articleHash, coins, privatekey, qosnonce, qscnonce)
	if err != nil {
		log.Printf("buyAd err:%s", err.Error())
		result.Code = "-1"
		result.Reason = err.Error()
		return result.Marshal()
	}

	js, err := cdc.MarshalJSON(tx)
	if err != nil {
		log.Printf("buyAd err:%s", err.Error())
		result.Code = "-1"
		result.Reason = err.Error()
		return result.Marshal()
	}
	result.Result = json.RawMessage(js)

	return result.Marshal()
}

func warpperInvestorTx(cdc *wire.Codec, articleHash string, amount int64) []qostxs.TransItem {
	investors, err := jianqian.ListInvestors(config.GetCLIContext().QSCCliContext, cdc, articleHash)
	var result []qostxs.TransItem
	log.Printf("buyAd warpperInvestorTx investors:%+v", investors)

	if err == nil {
		totalInvest := qbasetypes.NewInt(0)
		for _, v := range investors {
			totalInvest = totalInvest.Add(v.Invest)
		}

		log.Printf("buyAd warpperInvestorTx amount:%d, totalInvest:%d", amount, totalInvest.Int64())

		if !totalInvest.IsZero() {
			for _, v := range investors {
				result = append(
					result,
					warpperTransItem(
						v.Address,
						[]qbasetypes.BaseCoin{{Name: coinsName, Amount: qbasetypes.NewInt(amount * v.Invest.Int64() / totalInvest.Int64())}}))
			}
		}
	}

	return result
}

func getCommunityAddr(cdc *wire.Codec) (qbasetypes.Address, error) {
	communityPri := config.GetCLIContext().Config.Community
	if communityPri == "" {
		return nil, errors.New("no community")
	}

	_, addrben32, _ := utility.PubAddrRetrievalFromAmino(communityPri, cdc)
	community, err := types.AccAddressFromBech32(addrben32)
	if err != nil {
		return nil, err
	}

	return community, nil
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

func warpperReceivers(cdc *wire.Codec, article *jianqian.Articles, amount qbasetypes.BigInt) ([]qostxs.TransItem, error) {
	var result []qostxs.TransItem
	log.Printf("buyad warpperReceivers  article:%+v", article)

	investors, err := calculateRevenue(cdc, article, amount)
	if err != nil {
		return nil, err
	}

	for _, v := range investors {
		result = append(
			result,
			warpperTransItem(
				v.Address,
				[]qbasetypes.BaseCoin{{Name: coinsName, Amount: v.Revenue}}))
	}

	return mergeReceivers(result), nil
}

// calculateInvestorRevenue 计算投资者收入
func calculateInvestorRevenue(cdc *wire.Codec, articleHash string, amount qbasetypes.BigInt) (jianqian.Investors, error) {
	investors, err := jianqian.ListInvestors(config.GetCLIContext().QSCCliContext, cdc, articleHash)
	log.Printf("buyAd calculateInvestorRevenue investors:%+v", investors)
	if err != nil {
		return nil, err
	}

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
		}
	}

	return investors, nil
}

// calculateRevenue 计算收入
func calculateRevenue(cdc *wire.Codec, article *jianqian.Articles, amount qbasetypes.BigInt) ([]jianqian.Investor, error) {
	var result []jianqian.Investor
	log.Printf("buyad calculateRevenue  article:%+v", article)

	// 作者地址
	authorTotal := amount.Mul(qbasetypes.NewInt(int64(article.ShareAuthor)).Div(qbasetypes.NewInt(100)))
	result = append(
		result,
		jianqian.Investor{
			InvestorType: jianqian.InvestorTypeAuthor, // 投资者类型
			Address:      article.Authoraddress,       // 投资者地址
			Invest:       qbasetypes.NewInt(0),        // 投资金额
			Revenue:      authorTotal,                 // 投资收益
		})

	// 原创作者地址
	shareOriginalTotal := amount.Mul(qbasetypes.NewInt(int64(article.ShareOriginalAuthor)).Div(qbasetypes.NewInt(100)))
	result = append(
		result,
		jianqian.Investor{
			InvestorType: jianqian.InvestorTypeOriginalAuthor, // 投资者类型
			Address:      article.OriginalAuthor,              // 投资者地址
			Invest:       qbasetypes.NewInt(0),                // 投资金额
			Revenue:      shareOriginalTotal,                  // 投资收益
		})

	// 投资者收入分配
	investorShouldTotal := amount.Mul(qbasetypes.NewInt(int64(article.ShareInvestor)).Div(qbasetypes.NewInt(100)))
	investors, err := calculateInvestorRevenue(cdc, article.ArticleHash, investorShouldTotal)
	if err != nil {
		return nil, err
	}
	result = append(result, investors...)

	shareCommunityAddr, err := getCommunityAddr(cdc)
	if err == nil {
		shareCommunityTotal := amount.Sub(authorTotal).Sub(shareOriginalTotal).Sub(investors.TotalRevenue())
		// 社区收入比例
		result = append(
			result,
			jianqian.Investor{
				InvestorType: jianqian.InvestorTypeCommunity, // 投资者类型
				Address:      shareCommunityAddr,             // 投资者地址
				Invest:       qbasetypes.NewInt(0),           // 投资金额
				Revenue:      shareCommunityTotal,            // 投资收益
			})
	}

	return result, nil
}

// buyAd 投资广告
func buyAd(cdc *wire.Codec, chainId, articleHash, coins, privatekey string, qosnonce, qscnonce int64) (*txs.TxStd, error) {
	article, err := jianqian.QueryArticle(cdc, config.GetCLIContext().QSCCliContext, articleHash)
	log.Printf("buyad.buyAd QueryArticle article:%+v, err:%+v", article, err)
	if err != nil {
		return nil, err
	}

	if err := checkArticleBase(article, time.Now()); err != nil {
		return nil, err
	}

	articleBuy, err := jianqian.QueryArticleBuyer(cdc, config.GetCLIContext().QSCCliContext, articleHash)
	log.Printf("buyad.buyAd QueryArticleBuyer articleBuy:%+v, err:%+v", articleBuy, err)
	if err == nil {
		if articleBuy.CheckStatus != jianqian.CheckStatusFail {
			return nil, errors.New("已被购买")
		}
	}

	cs, err := types.ParseCoins(coins)
	if err != nil {
		return nil, err
	}

	if len(cs) != 1 {
		return nil, errors.New("one coin need")
	}

	for _, v := range cs {
		if v.Denom != coinsName {
			return nil, fmt.Errorf("only support %s", coinsName)
		}
	}

	var amount int64
	_, addrben32, priv := utility.PubAddrRetrievalFromAmino(privatekey, cdc)
	buyer, err := types.AccAddressFromBech32(addrben32)
	var ccs []qbasetypes.BaseCoin
	for _, coin := range cs {
		amount = coin.Amount.Int64()
		ccs = append(ccs, qbasetypes.BaseCoin{
			Name:   coin.Denom,
			Amount: qbasetypes.NewInt(coin.Amount.Int64()),
		})
	}
	qosnonce += 1
	var transferTx qostxs.TxTransfer
	transferTx.Senders = []qostxs.TransItem{warpperTransItem(buyer, ccs)}
	receivers, err := warpperReceivers(cdc, article, qbasetypes.NewInt(amount))
	if err != nil {
		return nil, err
	}
	transferTx.Receivers = receivers
	gas := qbasetypes.NewInt(int64(0))
	stx := txs.NewTxStd(transferTx, config.GetCLIContext().Config.QOSChainID, gas)
	signature, _ := stx.SignTx(priv, qosnonce, config.GetCLIContext().Config.QSCChainID)
	stx.Signature = []txs.Signature{txs.Signature{
		Pubkey:    priv.PubKey(),
		Signature: signature,
		Nonce:     qosnonce,
	}}

	qscnonce += 1
	it := &BuyTx{}
	it.ArticleHash = []byte(articleHash)
	it.Std = stx
	tx2 := txs.NewTxStd(it, config.GetCLIContext().Config.QSCChainID, stx.MaxGas)
	signature2, _ := tx2.SignTx(priv, qscnonce, config.GetCLIContext().Config.QSCChainID)
	tx2.Signature = []txs.Signature{txs.Signature{
		Pubkey:    priv.PubKey(),
		Signature: signature2,
		Nonce:     qscnonce,
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
	var result ResultBuy
	result.Code = "0"

	buyer, err := jianqian.QueryArticleBuyer(cdc, config.GetCLIContext().QSCCliContext, articleHash)
	if err != nil {
		log.Printf("QueryArticleBuyer err:%s", err.Error())
		result.Code = "-1"
		result.Reason = err.Error()
		return result.Marshal()
	}

	js, err := cdc.MarshalJSON(buyer)
	if err != nil {
		log.Printf("buyAd err:%s", err.Error())
		result.Code = "-1"
		result.Reason = err.Error()
		return result.Marshal()
	}
	result.Result = json.RawMessage(js)

	return result.Marshal()
}
