// Copyright 2018 The QOS Authors

package investad

import (
	"encoding/json"
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

const coinsName = "AOE"

var tempAddr = qbasetypes.Address("99999999999999999999")

// InvestAdBackground 提交到链上
func InvestAdBackground(cdc *wire.Codec, txb string, timeout time.Duration) string {
	ts := new(txs.TxStd)
	err := cdc.UnmarshalJSON([]byte(txb), ts)
	if err != nil {
		return common.InternalError(err.Error()).Marshal()
	}

	cliCtx := *config.GetCLIContext().QSCCliContext
	_, commitresult, err := utils.SendTx(cliCtx, cdc, ts)
	if err != nil {
		return common.NewErrorResult(common.ResultCodeInternalError, 0, "", err.Error()).Marshal()
	}

	return common.NewSuccessResult(cdc, commitresult.Height, commitresult.Hash.String(), "").Marshal()
}

//func fetchResult(cdc *wire.Codec, heigth1 string, tx1 string) (string, error) {
//	qstarskey := "heigth:" + heigth1 + ",hash:" + tx1
//	d, err := config.GetCLIContext().QSCCliContext.QueryStore([]byte(qstarskey), common.QSCResultMapperName)
//	if err != nil {
//		return "", err
//	}
//	if d == nil {
//		return "", nil
//	}
//	var res []byte
//	err = cdc.UnmarshalBinaryBare(d, &res)
//	if err != nil {
//		return "", err
//	}
//	return string(res), err
//}

// InvestAd 投资广告
func InvestAd(cdc *wire.Codec, articleHash, amount, privatekey,otheraddr string, qscnonce int64) string {
	var result common.Result
	result.Code = common.ResultCodeSuccess

	tx, berr := investAd(cdc, articleHash, amount, privatekey, otheraddr, qscnonce)
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
func investAd(cdc *wire.Codec, articleHash, coins, privatekey ,otheraddr string, qscnonce int64) (*txs.TxStd, *InvestadErr) {
	article, err := jianqian.QueryArticle(cdc, config.GetCLIContext().QSCCliContext, articleHash)
	log.Printf("investad.investAd QueryArticle article:%+v, err:%+v", article, err)
	if err != nil {
		return nil, NewInvestadErr(InvalidArticleErrCode, err.Error())
	}
	amount,ok:=qbasetypes.NewIntFromString(coins)
	if !ok{
		return nil, NewInvestadErr(CoinsErrCode, err.Error())
	}
	_, addrben32, priv := utility.PubAddrRetrievalFromAmino(privatekey, cdc)
	investor, _ := types.AccAddressFromBech32(addrben32)
	gas := qbasetypes.NewInt(int64(200000))
	qscnonce += 1
	it := &InvestTx{}
	it.ArticleHash = []byte(articleHash)
	it.Invest = amount
	it.Address=investor
	//it.OtherAddr=otheraddr
	fmt.Println(articleHash,amount,investor,otheraddr)

	tx2 := txs.NewTxStd(it, config.GetCLIContext().Config.QSCChainID, gas)
	signature2, _ := tx2.SignTx(priv, qscnonce, config.GetCLIContext().Config.QSCChainID, config.GetCLIContext().Config.QSCChainID)
	tx2.Signature = []txs.Signature{txs.Signature{
		Pubkey:    priv.PubKey(),
		Signature: signature2,
		Nonce:     qscnonce,
	}}
	return tx2, nil
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
