package recharge

import (
	"fmt"
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/txs"
	qbasetypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/client/utils"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/utility"
	"github.com/QOSGroup/qstars/wire"
	"github.com/QOSGroup/qstars/x/common"
	"github.com/QOSGroup/qstars/x/jianqian"
)

const (
	COINS_PARA_LEN_ERR     = "701" //参数长度不一致

	COINS_QUERY_ERR        = "702" //查询跨链结果错误
)
//
//// 余额变动交易 提交到链上
//func RechargeBackground(cdc *wire.Codec, txb string, timeout time.Duration) string {
//	ts := new(txs.TxStd)
//	err := cdc.UnmarshalJSON([]byte(txb), ts)
//	if err != nil {
//		return common.InternalError(err.Error()).Marshal()
//	}
//	cliCtx := *config.GetCLIContext().QSCCliContext
//	_, commitresult, err := utils.SendTx(cliCtx, cdc, ts)
//	if err != nil {
//		return common.NewErrorResult(common.ResultCodeInternalError, 0, "", err.Error()).Marshal()
//	}
//	return common.NewSuccessResult(cdc, commitresult.Height, commitresult.Hash.String(), "").Marshal()
//}

func Recharge(cdc *wire.Codec, amount, privatekey,address, cointype,isDeposit string) string {
	var result common.Result
	result.Code = common.ResultCodeSuccess
	tx, berr := recharge(cdc, amount, privatekey,address, cointype,isDeposit)
	if berr != "" {
		return berr
	}
	cliCtx := *config.GetCLIContext().QSCCliContext
	_, commitresult, err := utils.SendTx(cliCtx, cdc, tx)
	if err != nil {
		return common.NewErrorResult(common.ResultCodeInternalError, 0, "", err.Error()).Marshal()
	}
	return common.NewSuccessResult(cdc, commitresult.Height, commitresult.Hash.String(), "").Marshal()
}

func recharge(cdc *wire.Codec, coins, privatekey,address, cointype,isDeposit string) (*txs.TxStd, string) {
	amount, ok := qbasetypes.NewIntFromString(coins)
	if !ok {
		return nil, common.NewErrorResult(COINS_PARA_LEN_ERR, 0, "", "amount format error").Marshal()
	}
	_, addrben32, priv := utility.PubAddrRetrievalFromAmino(privatekey, cdc)
	from, _ := types.AccAddressFromBech32(addrben32)
	gas := qbasetypes.NewInt(int64(200000))

	key := account.AddressStoreKey(from)
	var qscnonce int64 = 0
	qscacc, err := getQSCAcc(key, cdc)
	if err != nil {
		qscnonce = 0
	} else {
		qscnonce = int64(qscacc.Nonce)
	}

	qscnonce += 1
	it := &jianqian.CoinsTx{}
	it.Address = from
	it.Cointype=cointype
	it.ChangeType=isDeposit
	it.Amount=amount
	tx:=RechargeTx{address,it}
	fmt.Println(from, amount, cointype, isDeposit)
	tx2 := txs.NewTxStd(tx, config.GetCLIContext().Config.QSCChainID, gas)
	signature2, _ := tx2.SignTx(priv, qscnonce, config.GetCLIContext().Config.QSCChainID, config.GetCLIContext().Config.QSCChainID)
	tx2.Signature = []txs.Signature{txs.Signature{
		Pubkey:    priv.PubKey(),
		Signature: signature2,
		Nonce:     qscnonce,
	}}
	return tx2, ""
}
