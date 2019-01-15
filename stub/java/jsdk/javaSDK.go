package jsdk

import (
	"encoding/json"
	"fmt"
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/utility"
	"github.com/QOSGroup/qstars/x/common"
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

//for investAdbckaground testing
type ResultInvest struct {
	Code   string          `json:"code"`
	Height int64           `json:"height"`
	Hash   string          `json:"hash,omitempty"`
	Reason string          `json:"reason,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
}

func InvestAdBackground(txb string) string {
	timeout := time.Second * 60
	var ri ResultInvest
	err := json.Unmarshal([]byte(txb), &ri)
	if err != nil {
		return err.Error()
	}
	Txresult := string(ri.Result)
	result := investad.InvestAdBackground(CDC, Txresult, timeout)
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

	qosacc, err1 := config.GetCLIContext().QOSCliContext.GetAccount(key, CDC)
	if err1 != nil {
		return err1.Error()
	}
	qosnonce := int64(qosacc.Nonce)

	qscacc, err2 := config.GetCLIContext().QSCCliContext.GetAccount(key, CDC)
	var qscnonce int64
	if err2 != nil {
		qscnonce = int64(1)
	}
	qscnonce = int64(qscacc.Nonce)
	tx := buyad.BuyAd(CDC, chainid, articleHash, coins, buyer, qosnonce, qscnonce)

	var rb ResultInvest
	if err := json.Unmarshal([]byte(tx), &rb); err != nil {
		return fmt.Sprintf("Unmarshal tx error:%s ", err.Error())
	}

	if rb.Code != "0" {
		return tx
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

func QSCCommitResultCheck(txhash, height string) string {
	qstarskey := "heigth:" + height + ",hash:" + txhash
	d, err := config.GetCLIContext().QSCCliContext.QueryStore([]byte(qstarskey), common.QSCResultMapperName)
	if err != nil {
		return err.Error()
	}
	if d == nil {
		return ""
	}
	var res []byte
	err = CDC.UnmarshalBinaryBare(d, &res)
	if err != nil {
		return err.Error()
	}
	return string(res)
}
