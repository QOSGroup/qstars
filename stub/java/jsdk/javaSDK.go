package jsdk

import (
	"encoding/json"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/x/common"
	"github.com/QOSGroup/qstars/x/jianqian/article"
	"github.com/QOSGroup/qstars/x/jianqian/auction"
	"github.com/QOSGroup/qstars/x/jianqian/coins"
	"github.com/QOSGroup/qstars/x/jianqian/buyad"

	"github.com/QOSGroup/qstars/x/jianqian/investad"
	"time"
)

//发放活动奖励 aoe
//addrs        奖励地址
//cns          数额
//causecodes   活动类型编号
//causestrings 活动描述
func DispatchCoins(addrs, cns, causecodes, causestrings, gas string) string {
	result := coins.DispatchAOE(CDC, CONF, addrs, cns, causecodes, causestrings, gas)
	return result
}

//创建文章
//AuthorAddr                   //作者地址(必填) 0qos 1cosmos
//AuthorOtherAddr              //作者其他帐户地址
//ArticleType                  //是否原创 0原创 1转载
//ArticleHash                  //作品唯一标识hash
//ShareAuthor                  //作者收入比例(必填)
//ShareOriginalAuthor          //原创收入比例(转载作品必填)
//ShareCommunity               //社区收入比例(必填)
//ShareInvestor                //投资者收入比例(必填)
//InvestHours                  //可供投资的小时数(必填)
//BuyHours                     //可供购买广告位的小时数(必填)
//CoinType                     //币种
func NewArticle(authorAddress, authorOtherAddr, articleType, articleHash, shareAuthor, shareOriginAuthor, shareCommunity, shareInvestor, endInvestDate, endBuyDate, cointype string) string {
	result := article.NewArticle(CDC, CONF, authorAddress, authorOtherAddr, articleType, articleHash, shareAuthor, shareOriginAuthor, shareCommunity, shareInvestor, endInvestDate, endBuyDate, cointype)
	return result
}

//竞拍广告上链
func AcutionAdBackground(txb string) string {
	timeout := time.Second * 60
	var ri ResultInvest
	err := json.Unmarshal([]byte(txb), &ri)
	if err != nil {
		return err.Error()
	}
	Txresult := string(ri.Result)
	result := auction.AcutionAdBackground(CDC, Txresult, timeout)
	return result
}

//查询当前最高出价
//txHash 广告位标识
func QueryMaxAcution(txHash string) string {
	result := auction.QueryMaxAcution(CDC, config.GetCLIContext().QSCCliContext, txHash)

	return result
}

//查询全部竞拍信息
//txHash 广告位标识
func QueryAllAcution(txHash string) string {
	result := auction.QueryAllAcution(CDC, config.GetCLIContext().QSCCliContext, txHash)
	return result
}

//分配利润 竞拍期过后分配出价最高者分配  竞拍失败部分原路退回
//txHash 广告位标识
func Distribution(txHash string) string {
	result:=buyad.BuyAd(CDC,txHash)
	return result
}

//投资上链
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

//查询投资信息
func RetrieveInvestors(articleHash string) string {
	result := investad.RetrieveInvestors(CDC, articleHash)
	return result
}

//查询文章信息
func QueryArticle(articleHash string) string {
	result := article.GetArticle(CDC, articleHash)
	return result
}

//查询活动奖励信息
func QueryCoins(txHash string) string {
	result := coins.GetCoins(CDC, config.GetCLIContext().QSCCliContext, txHash)
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

//for investAdbckaground testing
type ResultInvest struct {
	Code   string          `json:"code"`
	Height int64           `json:"height"`
	Hash   string          `json:"hash,omitempty"`
	Reason string          `json:"reason,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
}
