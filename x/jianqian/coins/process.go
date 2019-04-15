package coins

import (
	"encoding/json"
	"fmt"
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/client/context"
	"github.com/QOSGroup/qstars/client/utils"
	"github.com/QOSGroup/qstars/config"
	qstartypes "github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/utility"
	"log"
	"time"

	"github.com/QOSGroup/qstars/wire"
	"github.com/QOSGroup/qstars/x/common"
	"github.com/QOSGroup/qstars/x/jianqian"
	"github.com/QOSGroup/qstars/x/jianqian/tx"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"strings"
)

const (
	COINS_PARA_LEN_ERR     = "101" //参数长度不一致
	COINS_PRIV_ERR         = "102" //私钥获取地址错误
	COINS_SENDTX_ERR       = "103" //交易出错
	COINS_FETCH_RESULT_ERR = "104" //查询跨链结果错误
	COINS_QUERY_ERR        = "105" //查询跨链结果错误
)

//发放活动奖励 一转多
func DispatchSend(cdc *wire.Codec, ctx *config.CLIConfig, privkey string, to []types.Address, amount []types.BigInt, causecode []string, causeStr []string) string {
	tolen := len(to)
	//判断长度是否一致
	if tolen != len(amount) || tolen != len(causecode) || tolen != len(causeStr) {
		return common.NewErrorResult(COINS_PARA_LEN_ERR, 0, "", "Array parameter length is inconsistent").Marshal()
	}

	_, addrben32, priv := utility.PubAddrRetrievalFromAmino(privkey, cdc)
	from, err := qstartypes.AccAddressFromBech32(addrben32)
	if err != nil {
		return common.NewErrorResult(COINS_PRIV_ERR, 0, "", err.Error()).Marshal()
	}
	key := account.AddressStoreKey(from)
	var qscnonce int64 = 0
	acc, err := config.GetCLIContext().QSCCliContext.GetAccount(key, cdc)
	if err != nil {
		qscnonce = 0
	} else {
		qscnonce = int64(acc.Nonce)
	}
	qscnonce++
	fmt.Println("qosnonce", qscnonce)

	addmap := make(map[string]Recipient)
	for i, coin := range amount {
		if v, ok := addmap[to[i].String()]; ok {
			v.Amount.Add(coin)
		} else {
			recipient := Recipient{
				Address: to[i].String(),
				Amount:  coin,
			}
			addmap[to[i].String()] = recipient
		}
	}
	var newccs []Recipient
	for _, v := range addmap {
		newccs = append(newccs, v)
	}
	transtx := AOETx{from, newccs}
	var msg *txs.TxStd
	msg = genStdSendTx(cdc, transtx, priv, config.GetCLIContext().Config.QSCChainID, config.GetCLIContext().Config.QSCChainID, qscnonce)
	cliCtx := *config.GetCLIContext().QSCCliContext
	hash, commitresult, err := utils.SendTx(cliCtx, cdc, msg)
	return common.NewSuccessResult(cdc, commitresult.Height, hash, hash).Marshal()
}

//封装公链交易信息
func genStdSendTx(cdc *amino.Codec, sendTx txs.ITx, priKey ed25519.PrivKeyEd25519, tochainid, fromchainid string, nonce int64) *txs.TxStd {
	gas := types.NewInt(int64(200000))
	stx := txs.NewTxStd(sendTx, tochainid, gas)
	signature, _ := stx.SignTx(priKey, nonce, fromchainid, tochainid)
	stx.Signature = []txs.Signature{txs.Signature{
		Pubkey:    priKey.PubKey(),
		Signature: signature,
		Nonce:     nonce,
	}}
	return stx
}

//活动奖励发放
//address       接收奖励地址(必填) 多个地址用|隔开
//coins         接收奖励数额(必填) 多个地址用|隔开
//causecodes    奖励类型(必填) 多个地址用|隔开
//causestrings  奖励类型描述(必填) 多个地址用|隔开
//gas           gas费 默认为0
func DispatchAOE(cdc *wire.Codec, ctx *config.CLIConfig, address, coins, causecodes, causestrings, gas string) string {
	if address == "" || coins == "" || causecodes == "" || causestrings == "" {
		return common.NewErrorResult(COINS_PARA_LEN_ERR, 0, "", "Dispatch AOE parameters must not empty").Marshal()
	}
	addrs := strings.Split(address, "|")
	addlen := len(addrs)
	cois := strings.Split(coins, "|")
	codes := strings.Split(causecodes, "|")
	cstrs := strings.Split(causestrings, "|")

	if addlen != len(cois) || addlen != len(codes) || addlen != len(cstrs) {
		return common.NewErrorResult(COINS_PARA_LEN_ERR, 0, "", "Array parameter length is inconsistent").Marshal()
	}
	amounts := make([]types.BigInt, len(cois))
	for i, coinsv := range cois {
		if amou, ok := types.NewIntFromString(coinsv); ok {
			amounts[i] = amou
		} else {
			return common.NewErrorResult(COINS_PARA_LEN_ERR, 0, "", "amount format error").Marshal()
		}
	}
	toaddrss := make([]types.Address, addlen)
	for i, addrsv := range addrs {
		to, err := qstartypes.AccAddressFromBech32(addrsv)
		if err != nil {
			return common.NewErrorResult(COINS_PARA_LEN_ERR, 0, "", "address format error").Marshal()
		}
		toaddrss[i] = to
	}
	//cdc := star.MakeCodec()
	privkey := tx.GetConfig().Dappowner
	return DispatchSend(cdc, ctx, privkey, toaddrss, amounts, codes, cstrs)
}

func GetBlance(cdc *wire.Codec, ctx *context.CLIContext, tx string) string {
	coins, err := jianqian.QueryBlance(cdc, ctx, tx)
	if err != nil {
		return common.NewErrorResult(COINS_QUERY_ERR, 0, "", err.Error()).Marshal()
	}
	if coins == nil {
		return common.NewErrorResult(COINS_QUERY_ERR, 0, "", fmt.Sprintf("query blance failure,%s not exist", tx)).Marshal()
	}
	return common.NewSuccessResult(cdc, 0, "", coins).Marshal()
}

// 余额变动交易 提交到链上
func TransferBackground(cdc *wire.Codec, txb string, timeout time.Duration) string {
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
func Transfer(cdc *wire.Codec, amount, privatekey, to,cointype string, qscnonce int64) string {
	var result common.Result
	result.Code = common.ResultCodeSuccess
	tx, berr := transfer(cdc, amount, privatekey, to,cointype, qscnonce)
	if berr != "" {
		return berr
	}
	js, err := cdc.MarshalJSON(tx)
	if err != nil {
		log.Printf("CoinsChange err:%s", err.Error())
		result.Code = common.ResultCodeInternalError
		result.Reason = err.Error()
		return result.Marshal()
	}
	result.Result = json.RawMessage(js)
	return result.Marshal()
}

// 转账
func transfer(cdc *wire.Codec, coins, privatekey, to,cointype string, qscnonce int64) (*txs.TxStd, string) {
	amount, ok := types.NewIntFromString(coins)
	if !ok {
		return nil, common.NewErrorResult(COINS_PARA_LEN_ERR, 0, "", "amount format error").Marshal()
	}
	_, addrben32, priv := utility.PubAddrRetrievalFromAmino(privatekey, cdc)
	from, _ := qstartypes.AccAddressFromBech32(addrben32)
	gas := types.NewInt(int64(0))
	qscnonce += 1
	it := &CoinsTx{}
	it.From=from
	it.CoinType=cointype
	recipient:=Recipient{to,amount}
	it.To=[]Recipient{recipient}
	fmt.Println(it,to,cointype, amount)
	tx2 := txs.NewTxStd(it, config.GetCLIContext().Config.QSCChainID, gas)
	signature2, _ := tx2.SignTx(priv, qscnonce, config.GetCLIContext().Config.QSCChainID, config.GetCLIContext().Config.QSCChainID)
	tx2.Signature = []txs.Signature{txs.Signature{
		Pubkey:    priv.PubKey(),
		Signature: signature2,
		Nonce:     qscnonce,
	}}
	return tx2, ""
}