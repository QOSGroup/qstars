package jsdk

import (
	"github.com/QOSGroup/qstars/x/jianqian/article"
	"github.com/QOSGroup/qstars/x/jianqian/buyad"
	"github.com/QOSGroup/qstars/x/jianqian/coins"
	"github.com/QOSGroup/qstars/x/jianqian/investad"
	"time"
)

func DispatchCoins(addrs, cns, causecodes, causestrings, gas string) string {
	result := coins.DispatchAOE(CDC, CONF, addrs, cns, causecodes, causestrings, gas)
	output, err := CDC.MarshalJSON(result)
	if err != nil {
		return err.Error()
	}
	return string(output)
}

func NewArticle(authorAddress, originAuthor, articleHash, shareAuthor, shareOriginAuthor, shareCommunity, shareInvestor, endInvestDate, endBuyDate string) string {
	result := article.NewArticle(CDC, CONF, authorAddress, originAuthor, articleHash, shareAuthor, shareOriginAuthor, shareCommunity, shareInvestor, endInvestDate, endBuyDate)
	output, err := CDC.MarshalJSON(result)
	if err != nil {
		return err.Error()
	}
	return string(output)
}

func InvestAdBackground(txb string) string {
	timeout := time.Second * 60
	result := investad.InvestAdBackground(CDC, txb, timeout)
	output, err := CDC.MarshalJSON(result)
	if err != nil {
		return err.Error()
	}
	return string(output)

}

func BuyAd(chainId, articleHash, coins, privatekey string, nonce int64) string {
	result := buyad.BuyAd(CDC, chainId, articleHash, coins, privatekey, nonce, nonce)
	output, err := CDC.MarshalJSON(result)
	if err != nil {
		return err.Error()
	}
	return string(output)
}
