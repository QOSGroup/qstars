package article

import (
	"fmt"
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/client/utils"
	"github.com/QOSGroup/qstars/config"
	qstartypes "github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/utility"
	"github.com/QOSGroup/qstars/x/common"
	"github.com/QOSGroup/qstars/x/jianqian"
	"github.com/QOSGroup/qstars/x/jianqian/tx"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"strconv"
)

const (
	ARTICLE_PRIV_ERR            = "201" //上传者私钥获取地址错误
	ARTICLE_ORIGIN_ERR          = "202" //原作者获取地址错误
	ARTICLE_AUTHOR_SHARE_ERR    = "203" // 作者收入比例出错
	ARTICLE_ORIGIN_SHARE_ERR    = "204" // 原创入比例出错
	ARTICLE_COMMUNITY_SHARE_ERR = "205" // 社区收入比例出错
	ARTICLE_INVESTOR_SHARE_ERR  = "206" // 投资者收入比例出错
	ARTICLE_INVESTOR_DATE_ERR   = "207" // 投资期限出错
	ARTICLE_BUY_DATE_ERR        = "208" // 购买期限比例出错
	ARTICLE_ADDRES_ERR          = "209" // 地址转换错误
	ARTICLE_PRIV_AUTHOR_ERR     = "210" // 非作者本人私钥
	ARTICLE_SENDTX_ERR          = "211" //交易出错

	ARTICLE_QUERY_ERR = "212" //查询跨链结果错误
)

//type ResultArticle struct {
//	Code   string          `json:"code"`
//	Reason string          `json:"reason,omitempty"`
//	Result json.RawMessage `json:"result,omitempty"`
//}
//func InternalError(reason string) ResultArticle {
//	return ResultArticle{Code: "-1", Reason: reason}
//}
//func (ri ResultArticle) Marshal() string {
//	jsonBytes, err := json.MarshalIndent(ri, "", "  ")
//	if err != nil {
//		log.Printf("InvestAd err:%s", err.Error())
//		return InternalError(err.Error()).Marshal()
//	}
//	return string(jsonBytes)
//}

//上传新作品
//
//authoraddress          作者地址(必填)
//originalAuthor          原创作者地址(为空表示原创)
//articleHash            作品唯一标识hash
//shareAuthor            作者收入比例(必填)
//shareOriginalAuthor    原创收入比例(转载作品必填)
//shareCommunity         社区收入比例(必填)
//shareInvestor          投资者收入比例(必填)
//endInvestDate          投资结束时间(必填) 单位 小时
//endBuyDate             广告位购买结果时间(必填) 单位 小时
func NewArticle(cdc *amino.Codec, ctx *config.CLIConfig, authorAddress, authorOtherAddr ,articleType, articleHash, shareAuthor, shareOriginalAuthor,
	shareCommunity, shareInvestor, endInvestDate, endBuyDate string,coinType string) string {

	privkey := tx.GetConfig().Dappowner


	authorAddr, err := qstartypes.AccAddressFromBech32(authorAddress)
	if err != nil {
		return common.NewErrorResult(ARTICLE_ADDRES_ERR, 0, "", err.Error()).Marshal()
	}

	//var originaladdr types.Address
	//if strings.TrimSpace(originalAuthor) != "" {
	//	originaladdr, _ = qstartypes.AccAddressFromBech32(originalAuthor)
	//	if err != nil {
	//		return common.NewErrorResult(ARTICLE_ORIGIN_ERR, 0, "", err.Error()).Marshal()
	//	}
	//}

	articletype, err := strconv.Atoi(articleType)
	if err != nil {
		return common.NewErrorResult(ARTICLE_AUTHOR_SHARE_ERR, 0, "", err.Error()).Marshal()
	}

	authshare, err := strconv.Atoi(shareAuthor)
	if err != nil {
		return common.NewErrorResult(ARTICLE_AUTHOR_SHARE_ERR, 0, "", err.Error()).Marshal()
	}
	origshare, err := strconv.Atoi(shareOriginalAuthor)
	if err != nil {
		return common.NewErrorResult(ARTICLE_ORIGIN_SHARE_ERR, 0, "", err.Error()).Marshal()
	}
	commushare, err := strconv.Atoi(shareCommunity)
	if err != nil {
		return common.NewErrorResult(ARTICLE_COMMUNITY_SHARE_ERR, 0, "", err.Error()).Marshal()
	}
	invesshare, err := strconv.Atoi(shareInvestor)
	if err != nil {
		return common.NewErrorResult(ARTICLE_INVESTOR_SHARE_ERR, 0, "", err.Error()).Marshal()
	}
	investhours, err := strconv.Atoi(endInvestDate)
	if err != nil {
		return common.NewErrorResult(ARTICLE_INVESTOR_DATE_ERR, 0, "", err.Error()).Marshal()
	}
	if investhours <= 0 {
		return common.NewErrorResult(ARTICLE_INVESTOR_DATE_ERR, 0, "", "投资期需大于0").Marshal()
	}
	buyhours, err := strconv.Atoi(endBuyDate)
	if err != nil {
		return common.NewErrorResult(ARTICLE_BUY_DATE_ERR, 0, "", err.Error()).Marshal()
	}
	if buyhours <= 0 {
		return common.NewErrorResult(ARTICLE_BUY_DATE_ERR, 0, "", "广告竞拍期需大于0").Marshal()
	}

	_, addrben32, priv := utility.PubAddrRetrievalFromAmino(privkey, cdc)
	from, err := qstartypes.AccAddressFromBech32(addrben32)
	if err != nil {
		return common.NewErrorResult(ARTICLE_PRIV_ERR, 0, "", err.Error()).Marshal()
	}

	if authorAddr.String() != from.String() {
		//非作者本人私钥
		return common.NewErrorResult(ARTICLE_PRIV_AUTHOR_ERR, 0, "", "非作者本人私钥").Marshal()
	}

	key := account.AddressStoreKey(from)

	var nonce int64 = 0
	acc, err := config.GetCLIContext().QSCCliContext.GetAccount(key, cdc)
	if err != nil {
		nonce = 0
	} else {
		nonce = int64(acc.Nonce)
	}
	fmt.Println("nonce=", nonce)
	nonce++
	fromchainid := config.GetCLIContext().Config.QSCChainID
	tochainid := config.GetCLIContext().Config.QOSChainID
	tx := NewArticlesTx(authorAddr, authorOtherAddr,articletype, articleHash, authshare, origshare, commushare, invesshare, investhours, buyhours,coinType, types.ZeroInt())
	txsd := genStdSendTx(cdc, tx, priv, fromchainid, tochainid, nonce)
	cliCtx := *config.GetCLIContext().QSCCliContext
	_, commitresult, err1 := utils.SendTx(cliCtx, cdc, txsd)
	if err1 != nil {
		return common.NewErrorResult(ARTICLE_SENDTX_ERR, 0, "", err1.Error()).Marshal()
	}
	return common.NewSuccessResult(cdc, commitresult.Height, commitresult.Hash.String(), "").Marshal()
}

//封装公链交易信息
func genStdSendTx(cdc *amino.Codec, sendTx txs.ITx, priKey ed25519.PrivKeyEd25519, fromchainid string, tochainid string, nonce int64) *txs.TxStd {
	gas := types.NewInt(int64(0))
	stx := txs.NewTxStd(sendTx, fromchainid, gas)
	signature, _ := stx.SignTx(priKey, nonce, fromchainid, fromchainid)
	stx.Signature = []txs.Signature{txs.Signature{
		Pubkey:    priKey.PubKey(),
		Signature: signature,
		Nonce:     nonce,
	}}
	return stx
}

// GetArticle process of get Article
func GetArticle(cdc *amino.Codec, key string) string {
	article, err := jianqian.QueryArticle(cdc, config.GetCLIContext().QSCCliContext, key)
	if err != nil {
		return common.NewErrorResult(ARTICLE_QUERY_ERR, 0, "", err.Error()).Marshal()
	}
	if article == nil || article.ArticleHash == "" {
		return common.NewErrorResult(ARTICLE_QUERY_ERR, 0, "", fmt.Sprintf("query article failure,%s not exist", key)).Marshal()
	}
	return common.NewSuccessResult(cdc, 0, "", article).Marshal()

}
