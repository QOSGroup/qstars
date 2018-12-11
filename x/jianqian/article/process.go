package article

import (
	"encoding/json"
	"fmt"
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/client/utils"
	"github.com/QOSGroup/qstars/config"
	qstartypes "github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/utility"
	"github.com/QOSGroup/qstars/x/jianqian"
	"github.com/QOSGroup/qstars/x/jianqian/tx"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"log"
	"strconv"
	"strings"
)

type ResultArticle struct {
	Code   string          `json:"code"`
	Reason string          `json:"reason,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
}
func InternalError(reason string) ResultArticle {
	return ResultArticle{Code: "-1", Reason: reason}
}
func (ri ResultArticle) Marshal() string {
	jsonBytes, err := json.MarshalIndent(ri, "", "  ")
	if err != nil {
		log.Printf("InvestAd err:%s", err.Error())
		return InternalError(err.Error()).Marshal()
	}
	return string(jsonBytes)
}


//上传新作品
//
//authoraddress          作者地址(必填)
//originalAuthor          原创作者地址(为空表示原创)
//articleHash            作品唯一标识hash
//shareAuthor            作者收入比例(必填)
//shareOriginalAuthor    原创收入比例(转载作品必填)
//shareCommunity         社区收入比例(必填)
//shareInvestor          投资者收入比例(必填)
//endInvestDate          投资结束时间(必填)
//endBuyDate             广告位购买结果时间(必填)
func NewArticle(cdc *amino.Codec, ctx *config.CLIConfig, authorAddress, originalAuthor, articleHash, shareAuthor, shareOriginalAuthor,
	shareCommunity, shareInvestor, endInvestDate, endBuyDate string) string {
	privkey := tx.GetConfig().Dappowner
	authorAddr, err := qstartypes.AccAddressFromBech32(authorAddress)
	if err != nil {
		return err.Error()
	}
	var originaladdr types.Address
	if strings.TrimSpace(originalAuthor) != "" {
		originaladdr, _ = qstartypes.AccAddressFromBech32(originalAuthor)
		if err != nil {
			return err.Error()
		}
	}

	authshare, _ := strconv.Atoi(shareAuthor)
	origshare, _ := strconv.Atoi(shareOriginalAuthor)
	commushare, _ := strconv.Atoi(shareCommunity)
	invesshare, _ := strconv.Atoi(shareInvestor)
	investDays, _ := strconv.Atoi(endInvestDate)
	buydays, _ := strconv.Atoi(endBuyDate)

	tx := NewArticlesTx(authorAddr, originaladdr, articleHash, authshare, origshare, commushare, invesshare, investDays, buydays, types.ZeroInt())
	_, addrben32, priv := utility.PubAddrRetrievalFromAmino(privkey, cdc)
	from, err := qstartypes.AccAddressFromBech32(addrben32)
	if err != nil {
		return InternalError(err.Error()).Marshal()
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
	chainid := config.GetCLIContext().Config.QSCChainID
	txsd := genStdSendTx(cdc, tx, priv, chainid, nonce)
	cliCtx := *config.GetCLIContext().QSCCliContext
	_, _, err1 := utils.SendTx(cliCtx, cdc, txsd)

	if err1 != nil {
		return InternalError(err1.Error()).Marshal()
	}
	return "{Code:\"0\",Reason:\"\"}"
}

//封装公链交易信息
func genStdSendTx(cdc *amino.Codec, sendTx txs.ITx, priKey ed25519.PrivKeyEd25519, chainid string, nonce int64) *txs.TxStd {
	gas := types.NewInt(int64(0))
	stx := txs.NewTxStd(sendTx, chainid, gas)
	signature, _ := stx.SignTx(priKey, nonce, chainid)
	stx.Signature = []txs.Signature{txs.Signature{
		Pubkey:    priKey.PubKey(),
		Signature: signature,
		Nonce:     nonce,
	}}
	return stx
}

// GetArticle process of get Article
func GetArticle(cdc *amino.Codec, key string) string {
	var result ResultArticle
	result.Code = "0"
	article, err := jianqian.QueryArticle(cdc, config.GetCLIContext().QSCCliContext, key)
	if err!=nil{
		return InternalError(err.Error()).Marshal()
	}
	if article.ArticleHash==""{
		return InternalError(key+" ArticleHash not exist").Marshal()
	}
	js, err := cdc.MarshalJSON(article)
	if err != nil {
		log.Printf("GetCoins err:%s", err.Error())
		result.Code = "-1"
		result.Reason = err.Error()
		return result.Marshal()
	}
	result.Result = json.RawMessage(js)
	return result.Marshal()
}
