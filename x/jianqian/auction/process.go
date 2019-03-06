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
	"github.com/QOSGroup/qstars/x/jianqian/tx"
	"log"
	"strconv"
	"strings"
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

func AcutionAd(cdc *wire.Codec, articleHash, address, otherAddr, coinsType, coinAmount string, qscnonce, qosnonce int64) string {
	var result common.Result
	result.Code = common.ResultCodeSuccess
	amount, ok := qbasetypes.NewIntFromString(coinAmount)
	if !ok {
		result.Code = common.ResultCodeInternalError
		result.Reason = "AcutionAd invalid amount"
		return result.Marshal()
	}
	tx, err := acutionAd(cdc, articleHash, address, otherAddr, coinsType, amount, qscnonce, qosnonce)
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
func acutionAd(cdc *wire.Codec, articleHash, private, otherAddr, coinsType string, coinAmount qbasetypes.BigInt, qscnonce, qosnonce int64) (*txs.TxStd, error) {
	buyAdPrivatekey := config.GetCLIContext().Config.Adbuyermock
	_, addrben32, _ := utility.PubAddrRetrievalFromAmino(buyAdPrivatekey, cdc)
	adAddress, err := types.AccAddressFromBech32(addrben32)
	if err != nil {
		return nil, err
	}
	_, sendaddrben32, priv := utility.PubAddrRetrievalFromAmino(private, cdc)
	sendAddress, err := types.AccAddressFromBech32(sendaddrben32)
	if err != nil {
		return nil, err
	}

	if articleHash == "" {
		return nil, errors.New("invalid article hash")
	}
	//article, err := jianqian.QueryArticle(cdc, config.GetCLIContext().QSCCliContext, articleHash)
	//log.Printf("AcutionAd. QueryArticle article:%+v, err:%+v", article, err)
	if err != nil {
		return nil, err
	}
	var stx *txs.TxStd
	if strings.ToUpper(strings.TrimSpace(coinsType)) == "QOS" {
		//封装跨链
		var ccs []qbasetypes.BaseCoin
		ccs = append(ccs, qbasetypes.BaseCoin{
			Name:   "QOS",
			Amount: coinAmount,
		})
		qosnonce += 1
		send := []qbasetypes.Address{sendAddress}
		receive := []qbasetypes.Address{adAddress}
		itx := tx.NewTransfer(send, receive, ccs)

		stx = txs.NewTxStd(itx, config.GetCLIContext().Config.QOSChainID, qbasetypes.ZeroInt())
		signature, _ := stx.SignTx(priv, qosnonce, config.GetCLIContext().Config.QSCChainID, config.GetCLIContext().Config.QOSChainID)
		stx.Signature = []txs.Signature{txs.Signature{
			Pubkey:    priv.PubKey(),
			Signature: signature,
			Nonce:     qosnonce,
		}}
	}

	qscnonce += 1
	it := &AuctionTx{}
	it.ArticleHash = articleHash
	it.Wrapper = stx
	it.Address = sendAddress
	it.OtherAddr = otherAddr
	it.CoinsType = coinsType
	it.Gas = qbasetypes.ZeroInt()
	it.CoinAmount = coinAmount
	tx2 := txs.NewTxStd(it, config.GetCLIContext().Config.QSCChainID, qbasetypes.ZeroInt())
	signature2, _ := tx2.SignTx(priv, qscnonce, config.GetCLIContext().Config.QSCChainID, config.GetCLIContext().Config.QSCChainID)
	tx2.Signature = []txs.Signature{txs.Signature{
		Pubkey:    priv.PubKey(),
		Signature: signature2,
		Nonce:     qscnonce,
	}}

	return tx2, nil
}

//竞拍广告 提交上链
func AcutionAdBackground(cdc *wire.Codec, txb string, timeout time.Duration, coinType string) string {

	ts := new(txs.TxStd)
	err := cdc.UnmarshalJSON([]byte(txb), ts)

	if err != nil {
		return common.InternalError(err.Error()).Marshal()
	}

	//判断是否跨链
	auction := ts.ITx.(*AuctionTx)

	cliCtx := *config.GetCLIContext().QSCCliContext
	_, commitresult, err := utils.SendTx(cliCtx, cdc, ts)
	if err != nil {
		return common.NewErrorResult(AUCTION_PARA_LEN_ERR, 0, "", err.Error()).Marshal()
	}
	if strings.ToUpper(strings.TrimSpace(auction.CoinsType)) == "QOS" {
		height := strconv.FormatInt(commitresult.Height, 10)
		code := common.ResultCodeSuccess
		var reason string
		var result interface{}

		waittime, err := strconv.Atoi(config.GetCLIContext().Config.WaitingForQosResult)
		if err != nil {
			panic("WaitingForQosResult should be able to convert to integer." + err.Error())
		}
		counter := 0

		for {
			resultstr, err := fetchResult(cdc, height, commitresult.Hash.String())
			log.Printf("fetchResult result:%s, err:%+v\n", resultstr, err)
			if err != nil {
				log.Printf("fetchResult error:%s\n", err.Error())
				reason = err.Error()
				code = common.ResultCodeInternalError
				break
			}

			if resultstr != "" && resultstr != (AuctionStub{}).Name() {
				log.Printf("fetchResult result:[%+v]\n", resultstr)
				rs := []rune(resultstr)
				index1 := strings.Index(resultstr, " ")

				reason = ""
				result = string(rs[index1+1:])
				code = string(rs[:index1])
				break
			}

			if counter >= waittime {
				log.Println("time out")
				result = "time out"
				if resultstr == "" {
					code = common.ResultCodeQstarsTimeout
				} else {
					code = common.ResultCodeQOSTimeout
				}
				break
			}

			time.Sleep(500 * time.Millisecond)
			counter++
		}

		if code != common.ResultCodeSuccess {
			return common.NewErrorResult(code, commitresult.Height, commitresult.Hash.String(), reason).Marshal()
		}
		return common.NewSuccessResult(cdc, commitresult.Height, commitresult.Hash.String(), result).Marshal()

	}
	return common.NewSuccessResult(cdc, commitresult.Height, commitresult.Hash.String(), "").Marshal()

}

func fetchResult(cdc *wire.Codec, heigth1 string, tx1 string) (string, error) {
	// TODO qbase还没实现
	//qstarskey := "heigth:" + heigth1 + ",hash:" + tx1
	qstarskey := GetResultKey(heigth1, tx1)
	d, err := config.GetCLIContext().QSCCliContext.QueryStore([]byte(qstarskey), common.QSCResultMapperName)
	log.Printf("QueryStore: %+v, %+v\n", d, err)
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

func GetResultKey(heigth1 string, tx1 string) string {
	qstarskey := "heigth:" + heigth1 + ",hash:" + tx1
	return qstarskey
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
