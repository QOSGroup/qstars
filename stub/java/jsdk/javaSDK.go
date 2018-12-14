package jsdk

import (
	"encoding/json"
	"fmt"
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/utility"
	"github.com/QOSGroup/qstars/x/jianqian/article"
	"github.com/QOSGroup/qstars/x/jianqian/buyad"
	"github.com/QOSGroup/qstars/x/jianqian/coins"
	"github.com/QOSGroup/qstars/x/jianqian/investad"
	"time"
)

func DispatchCoins(addrs, cns, causecodes, causestrings, gas string) string {
	result := coins.DispatchAOE(CDC, CONF, addrs, cns, causecodes, causestrings, gas)
	//output, err := CDC.MarshalJSON(result)
	//if err != nil {
	//	return err.Error()
	//}
	return result
}

func NewArticle(authorAddress, originAuthor, articleHash, shareAuthor, shareOriginAuthor, shareCommunity, shareInvestor, endInvestDate, endBuyDate string) string {
	result := article.NewArticle(CDC, CONF, authorAddress, originAuthor, articleHash, shareAuthor, shareOriginAuthor, shareCommunity, shareInvestor, endInvestDate, endBuyDate)
	//output, err := CDC.MarshalJSON(result)
	//if err != nil {
	//	return err.Error()
	//}
	return result
}

func InvestAdBackground(txb string) string {
	timeout := time.Second * 60
	result := investad.InvestAdBackground(CDC, txb, timeout)
	//output, err := CDC.MarshalJSON(result)
	//if err != nil {
	//	return err.Error()
	//}
	return result

}

type ResultBuy struct {
	Code   string          `json:"code"`
	Reason string          `json:"reason,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
}

func BuyAd(articleHash, coins, buyer string) string {
	chainid := config.GetCLIContext().Config.QSCChainID
	if len(buyer) == 0 {
		buyer = config.GetCLIContext().Config.Adbuyermock
	}
	_, addrben32, _ := utility.PubAddrRetrievalFromAmino(buyer, CDC)
	from, err := types.AccAddressFromBech32(addrben32)
	key := account.AddressStoreKey(from)

	qosacc, _ := config.GetCLIContext().QOSCliContext.GetAccount(key, CDC)
	qosnonce := int64(qosacc.Nonce)

	qscacc, _ := config.GetCLIContext().QSCCliContext.GetAccount(key, CDC)
	qscnonce := int64(qscacc.Nonce)
	tx := buyad.BuyAd(CDC, chainid, articleHash, coins, buyer, qosnonce, qscnonce)

	var rb ResultBuy
	if err := json.Unmarshal([]byte(tx), &rb); err != nil {
		return fmt.Sprintf("Unmarshal tx error:%s ", err.Error())
	}

	if rb.Code != "0" {
		return fmt.Sprintf("InvestAd tx error:%s ", rb.Reason)
	}

	timeout := time.Second * 60
	result := buyad.BuyAdBackground(CDC, string(rb.Result), timeout)
	//output, err := CDC.MarshalJSON(result)
	if err != nil {
		return err.Error()
	}
	return result
}

func RetrieveInvestors(articleHash string) string {
	result := investad.RetrieveInvestors(CDC, articleHash)
	//output, err := CDC.MarshalJSON(result)
	//if err != nil {
	//	return err.Error()
	//}
	return result
}

func RetrieveBuyer(articleHash string) string {
	result := buyad.RetrieveBuyer(CDC, articleHash)
	//output, err := CDC.MarshalJSON(result)
	//if err != nil {
	//	return err.Error()
	//}
	return result
}

func QueryArticle(articleHash string) string {
	result := article.GetArticle(CDC, articleHash)
	//output, err := CDC.MarshalJSON(result)
	//if err != nil {
	//	return err.Error()
	//}
	return result
}

func QueryCoins(txHash string) string {
	result := coins.GetCoins(CDC, config.GetCLIContext().QSCCliContext, txHash)
	//output, err := CDC.MarshalJSON(result)
	//if err != nil {
	//	return err.Error()
	//}
	return result
}
