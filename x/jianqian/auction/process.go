package auction

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/QOSGroup/qbase/txs"
	qbasetypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/client/utils"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/utility"
	"github.com/QOSGroup/qstars/wire"
	"github.com/QOSGroup/qstars/x/common"
	"github.com/QOSGroup/qstars/x/jianqian"
	"log"
	"time"
)

//type SendResult struct {
//	Hash   string `json:"hash"`
//	Error  string `json:"error"`
//	Code   string `json:"code"`
//	Result string `json:"result"`
//	Heigth string `json:"heigth"`
//}
const (
	AUCTION_PARA_LEN_ERR     = "501" //参数长度不一致
	AUCTION_PRIV_ERR         = "502" //私钥获取地址错误
	AUCTION_SENDTX_ERR       = "503" //交易出错
	AUCTION_FETCH_RESULT_ERR = "504" //查询跨链结果错误
	AUCTION_QUERY_ERR        = "505" //查询跨链结果错误
)

func AcutionAd(cdc *wire.Codec, articleHash, address,  coinsType, coinAmount string, qscnonces int64) string {
	var result common.Result
	result.Code = common.ResultCodeSuccess
	amount, ok := qbasetypes.NewIntFromString(coinAmount)
	if !ok {
		result.Code = common.ResultCodeInternalError
		result.Reason = "AcutionAd invalid amount"
		return result.Marshal()
	}
	tx, err := acutionAd(cdc, articleHash, address,  coinsType, amount, qscnonces)
	if err != nil {
		log.Printf("AcutionAd err:%s", err.Error())
		result.Code = common.ResultCodeInternalError
		result.Reason = err.Error()
		return result.Marshal()
	}

	js, err := cdc.MarshalJSON(tx)
	if err != nil {
		log.Printf("AcutionAd err:%s", err.Error())
		result.Code = common.ResultCodeInternalError
		result.Reason = err.Error()
		return result.Marshal()
	}
	result.Result = json.RawMessage(js)
	return result.Marshal()
}

// investAd 投资广告
func acutionAd(cdc *wire.Codec, articleHash, private,  coinsType string, coinAmount qbasetypes.BigInt, qscnonce int64) (*txs.TxStd, error) {
	_, sendaddrben32, priv := utility.PubAddrRetrievalFromAmino(private, cdc)
	sendAddress, err := types.AccAddressFromBech32(sendaddrben32)
	if err != nil {
		return nil, err
	}

	if articleHash == "" {
		return nil, errors.New("invalid article hash")
	}

	if err != nil {
		return nil, err
	}

	qscnonce += 1
	it := &AuctionTx{}
	it.ArticleHash = articleHash
	it.Address = sendAddress
	it.CoinsType = coinsType
	it.Gas = qbasetypes.ZeroInt()
	it.CoinAmount = coinAmount
	tx2 := txs.NewTxStd(it, config.GetCLIContext().Config.QSCChainID, qbasetypes.NewInt(200000))
	signature2, _ := tx2.SignTx(priv, qscnonce, config.GetCLIContext().Config.QSCChainID, config.GetCLIContext().Config.QSCChainID)
	tx2.Signature = []txs.Signature{txs.Signature{
		Pubkey:    priv.PubKey(),
		Signature: signature2,
		Nonce:     qscnonce,
	}}

	return tx2, nil
}

//竞拍广告 提交上链
func AcutionAdBackground(cdc *wire.Codec, txb string, timeout time.Duration) string {
	ts := new(txs.TxStd)
	err := cdc.UnmarshalJSON([]byte(txb), ts)
	if err != nil {
		return common.InternalError(err.Error()).Marshal()
	}
	cliCtx := *config.GetCLIContext().QSCCliContext
	_, commitresult, err := utils.SendTx(cliCtx, cdc, ts)
	if err != nil {
		return common.NewErrorResult(AUCTION_PARA_LEN_ERR, 0, "", err.Error()).Marshal()
	}
	return common.NewSuccessResult(cdc, commitresult.Height, commitresult.Hash.String(), "").Marshal()

}


func QueryMaxAcution(cdc *wire.Codec, key string) string {
	auctionMap, err := jianqian.QueryAllAcution(cdc, config.GetCLIContext().QSCCliContext, key)
	if err != nil {
		return common.NewErrorResult(AUCTION_QUERY_ERR, 0, "", err.Error()).Marshal()
	}
	if auctionMap == nil {
		return common.NewErrorResult(AUCTION_QUERY_ERR, 0, "", fmt.Sprintf("query auction failure,%s not exist", key)).Marshal()
	}
	auction, exist := auctionMap[jianqian.MAXPRICEKEY]
	if !exist {
		return common.NewErrorResult(AUCTION_QUERY_ERR, 0, "", fmt.Sprintf("query auction failure,%s not exist", key)).Marshal()
	}
	return common.NewSuccessResult(cdc, 0, "", auction).Marshal()
}

func QueryAllAcution(cdc *wire.Codec, key string) string {
	auctionMap, err := jianqian.QueryAllAcution(cdc, config.GetCLIContext().QSCCliContext, key)
	if err != nil {
		return common.NewErrorResult(AUCTION_QUERY_ERR, 0, "", err.Error()).Marshal()
	}
	if auctionMap == nil {
		return common.NewErrorResult(AUCTION_QUERY_ERR, 0, "", fmt.Sprintf("QueryAllAcution failure,%s not exist", key)).Marshal()
	}
	delete(auctionMap, jianqian.MAXPRICEKEY)
	return common.NewSuccessResult(cdc, 0, "", auctionMap).Marshal()
}
