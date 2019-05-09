package comm

import (
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/txs"
	qosaccount "github.com/QOSGroup/qos/types"
	"github.com/QOSGroup/qstars/client/utils"
	"github.com/QOSGroup/qstars/config"
	qstarstypes "github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/utility"
	"github.com/QOSGroup/qstars/wire"
	"github.com/QOSGroup/qstars/x/common"

	"github.com/QOSGroup/qbase/types"
)

func CommHandler(cdc *wire.Codec, funcName, privatekey string, args []string) string {
	var result common.Result
	result.Code = common.ResultCodeSuccess
	tx, berr := commHandler(cdc, funcName, privatekey, args)

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

func commHandler(cdc *wire.Codec, funcName, privatekey string, args []string) (*txs.TxStd, string) {
	_, addrben32, priv := utility.PubAddrRetrievalFromAmino(privatekey, cdc)
	from, _ := qstarstypes.AccAddressFromBech32(addrben32)
	gas := types.NewInt(int64(200000))
	key := account.AddressStoreKey(from)
	var qscnonce int64 = 0
	qscacc, err := getQSCAcc(key, cdc)
	if err != nil {
		qscnonce = 0
	} else {
		qscnonce = int64(qscacc.Nonce)
	}
	qscnonce += 1
	tx := &JianQianTx{}
	tx.Address = []types.Address{from}
	tx.FuncName = funcName
	tx.Args = args
	tx.Gas = gas

	qscchainid:=config.GetCLIContext().Config.QSCChainID

	tx2 := txs.NewTxStd(tx, qscchainid, gas)
	signature2, _ := tx2.SignTx(priv, qscnonce,qscchainid, qscchainid)
	tx2.Signature = []txs.Signature{txs.Signature{
		Pubkey:    priv.PubKey(),
		Signature: signature2,
		Nonce:     qscnonce,
	}}
	return tx2, ""
}

func getQSCAcc(address []byte, cdc *wire.Codec) (*qosaccount.QOSAccount, error) {
	return config.GetCLIContext().QSCCliContext.GetAccount(address, cdc)
}
