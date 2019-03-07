// Copyright 2018 The QOS Authors

package investad

import (
	"encoding/json"
	"fmt"
	"github.com/QOSGroup/qbase/txs"
	qbasetypes "github.com/QOSGroup/qbase/types"
	qostxs "github.com/QOSGroup/qos/module/transfer"
	qostxtype "github.com/QOSGroup/qos/module/transfer/types"
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

const coinsName = "AOE"

var tempAddr = qbasetypes.Address("99999999999999999999")

// InvestAdBackground 提交到链上
func InvestAdBackground(cdc *wire.Codec, txb string, timeout time.Duration) string {
	ts := new(txs.TxStd)
	err := cdc.UnmarshalJSON([]byte(txb), ts)
	fmt.Printf("InvestAdBackground ts:%+v, txb:%s\n", ts, txb)
	if err != nil {
		return common.InternalError(err.Error()).Marshal()
	}

	cliCtx := *config.GetCLIContext().QSCCliContext
	_, commitresult, err := utils.SendTx(cliCtx, cdc, ts)
	fmt.Printf("SendTx commitresult:%+v, err:%+v \n", commitresult, err)
	if err != nil {
		return common.NewErrorResult(common.ResultCodeInternalError, 0, "", err.Error()).Marshal()
	}

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

		if resultstr != "" && resultstr != (InvestadStub{}).Name() {
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

// InvestAd 投资广告
func InvestAd(cdc *wire.Codec, chainId, articleHash, coins, privatekey string, qosnonce, qscnonce int64) string {
	var result common.Result
	result.Code = common.ResultCodeSuccess

	tx, berr := investAd(cdc, chainId, articleHash, coins, privatekey, qosnonce, qscnonce)
	if berr != nil {
		log.Printf("investAd err:%s", berr.Error())
		result.Code = berr.Code()
		result.Reason = berr.Error()
		return result.Marshal()
	}

	js, err := cdc.MarshalJSON(tx)
	if err != nil {
		log.Printf("investAd err:%s", err.Error())
		result.Code = common.ResultCodeInternalError
		result.Reason = err.Error()
		return result.Marshal()
	}
	result.Result = json.RawMessage(js)

	return result.Marshal()
}

// investAd 投资广告
func investAd(cdc *wire.Codec, chainId, articleHash, coins, privatekey string, qosnonce, qscnonce int64) (*txs.TxStd, *InvestadErr) {
	article, err := jianqian.QueryArticle(cdc, config.GetCLIContext().QSCCliContext, articleHash)
	log.Printf("investad.investAd QueryArticle article:%+v, err:%+v", article, err)
	if err != nil {
		return nil, NewInvestadErr(InvalidArticleErrCode, err.Error())
	}

	cs, err := types.ParseCoins(coins)
	if err != nil {
		return nil, NewInvestadErr(CoinsErrCode, err.Error())
	}

	for _, v := range cs {
		if v.Denom != coinsName {
			return nil, CoinsErr
		}
	}

	_, addrben32, priv := utility.PubAddrRetrievalFromAmino(privatekey, cdc)
	investor, err := types.AccAddressFromBech32(addrben32)
	var ccs []qbasetypes.BaseCoin
	for _, coin := range cs {
		ccs = append(ccs, qbasetypes.BaseCoin{
			Name:   coin.Denom,
			Amount: qbasetypes.NewInt(coin.Amount.Int64()),
		})
	}
	qosnonce += 1

	transferTx := NewTransfer(investor, tempAddr, ccs)
	// TODO set zero, temp
	gas := qbasetypes.NewInt(int64(0))
	stx := txs.NewTxStd(transferTx, config.GetCLIContext().Config.QOSChainID, gas)
	signature, _ := stx.SignTx(priv, qosnonce, config.GetCLIContext().Config.QSCChainID, config.GetCLIContext().Config.QOSChainID)
	stx.Signature = []txs.Signature{txs.Signature{
		Pubkey:    priv.PubKey(),
		Signature: signature,
		Nonce:     qosnonce,
	}}

	qscnonce += 1
	it := &InvestTx{}
	it.ArticleHash = []byte(articleHash)
	it.Std = stx
	tx2 := txs.NewTxStd(it, config.GetCLIContext().Config.QSCChainID, stx.MaxGas)
	signature2, _ := tx2.SignTx(priv, qscnonce, config.GetCLIContext().Config.QSCChainID, config.GetCLIContext().Config.QSCChainID)
	tx2.Signature = []txs.Signature{txs.Signature{
		Pubkey:    priv.PubKey(),
		Signature: signature2,
		Nonce:     qscnonce,
	}}

	return tx2, nil
}

func warpperTransItem(addr qbasetypes.Address, coins []qbasetypes.BaseCoin) qostxtype.TransItem {
	var ti qostxtype.TransItem
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

// NewTransfer ...
func NewTransfer(sender qbasetypes.Address, receiver qbasetypes.Address, coin []qbasetypes.BaseCoin) qostxs.TxTransfer {
	var sendTx qostxs.TxTransfer

	sendTx.Senders = append(sendTx.Senders, warpperTransItem(sender, coin))
	sendTx.Receivers = append(sendTx.Receivers, warpperTransItem(receiver, coin))

	return sendTx
}

// RetrieveInvestors 查询投资者
func RetrieveInvestors(cdc *wire.Codec, articleHash string) string {
	var result common.Result
	result.Code = common.ResultCodeSuccess

	investors, err := jianqian.ListInvestors(config.GetCLIContext().QSCCliContext, cdc, articleHash)
	if err != nil {
		log.Printf("ListInvestors err:%s", err.Error())
		result.Code = common.ResultCodeInternalError
		result.Reason = err.Error()
		return result.Marshal()
	}

	js, err := cdc.MarshalJSON(investors)
	if err != nil {
		log.Printf("buyAd err:%s", err.Error())
		result.Code = common.ResultCodeInternalError
		result.Reason = err.Error()
		return result.Marshal()
	}
	result.Result = json.RawMessage(js)

	return result.Marshal()
}
