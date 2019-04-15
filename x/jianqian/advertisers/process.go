package advertisers

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

const (
	ADVERTISERS_PARA_LEN_ERR = "601" //数量格式错误
)

// 押金或赎回交易 提交到链上
func AdvertisersBackground(cdc *wire.Codec, txb string, timeout time.Duration) string {
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

//广告商押金或赎回
func Advertisers(cdc *wire.Codec, amount, privatekey, cointype,isDeposit string, qscnonce int64) string {
	var result common.Result
	result.Code = common.ResultCodeSuccess
	tx, berr := advertisers(cdc, amount, privatekey, cointype,isDeposit, qscnonce)
	if berr != "" {
		return berr
	}
	js, err := cdc.MarshalJSON(tx)
	if err != nil {
		log.Printf("Advertisers err:%s", err.Error())
		result.Code = common.ResultCodeInternalError
		result.Reason = err.Error()
		return result.Marshal()
	}
	result.Result = json.RawMessage(js)
	return result.Marshal()
}

// investAd 投资广告
func advertisers(cdc *wire.Codec, coins, privatekey, cointype,isDeposit string,qscnonce int64) (*txs.TxStd, string) {
	amount, ok := qbasetypes.NewIntFromString(coins)
	if !ok {
		return nil, common.NewErrorResult(ADVERTISERS_PARA_LEN_ERR, 0, "", "amount format error").Marshal()
	}
	_, addrben32, priv := utility.PubAddrRetrievalFromAmino(privatekey, cdc)
	investor, _ := types.AccAddressFromBech32(addrben32)
	gas := qbasetypes.NewInt(int64(200000))
	qscnonce += 1
	it := &jianqian.CoinsTx{}
	it.Address = investor
	it.Cointype=cointype
	it.ChangeType=isDeposit
	it.Amount=amount
	tx:=AdvertisersTx{it}
	fmt.Println(investor, amount, cointype, isDeposit)
	tx2 := txs.NewTxStd(tx, config.GetCLIContext().Config.QSCChainID, gas)
	signature2, _ := tx2.SignTx(priv, qscnonce, config.GetCLIContext().Config.QSCChainID, config.GetCLIContext().Config.QSCChainID)
	tx2.Signature = []txs.Signature{txs.Signature{
		Pubkey:    priv.PubKey(),
		Signature: signature2,
		Nonce:     qscnonce,
	}}
	return tx2, ""
}
