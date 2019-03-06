package coins

import (
	"fmt"
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/client/context"
	"github.com/QOSGroup/qstars/client/utils"
	"github.com/QOSGroup/qstars/config"
	qstartypes "github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/utility"
	"github.com/QOSGroup/qstars/wire"
	"github.com/QOSGroup/qstars/x/common"
	"github.com/QOSGroup/qstars/x/jianqian"
	"github.com/QOSGroup/qstars/x/jianqian/tx"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"log"
	"strings"
)

//type SendResult struct {
//	Hash   string `json:"hash"`
//	Error  string `json:"error"`
//	Code   string `json:"code"`
//	Result string `json:"result"`
//	Heigth string `json:"heigth"`
//}
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

	var ccs []jianqian.AoeAccount
	for i, coin := range amount {
		ccs = append(ccs, jianqian.AoeAccount{
			Address: to[i].String(),
			Amount:  coin,
		})
	}

	//合并接收地址
	addmap := make(map[string]jianqian.AoeAccount)
	for i, addv := range to {
		if v, ok := addmap[addv.String()]; ok {
			v.Amount.Add(ccs[i].Amount)
		} else {
			addmap[addv.String()] = ccs[i]
		}
	}
	var newccs []jianqian.AoeAccount
	for _, v := range addmap {
		newccs = append(newccs, v)
	}

	transtx := NewDispatchAOE(newccs, from, to, amount, causecode, causeStr, types.ZeroInt())

	var msg *txs.TxStd

	msg = genStdSendTx(cdc, transtx, priv, config.GetCLIContext().Config.QSCChainID, config.GetCLIContext().Config.QSCChainID, qscnonce)

	return wrapperResult(cdc, msg)
}

//封装公链交易信息
func genStdSendTx(cdc *amino.Codec, sendTx txs.ITx, priKey ed25519.PrivKeyEd25519, tochainid, fromchainid string, nonce int64) *txs.TxStd {
	gas := types.NewInt(int64(0))
	stx := txs.NewTxStd(sendTx, tochainid, gas)

	signature, _ := stx.SignTx(priKey, nonce, fromchainid, tochainid)
	stx.Signature = []txs.Signature{txs.Signature{
		Pubkey:    priKey.PubKey(),
		Signature: signature,
		Nonce:     nonce,
	}}
	return stx
}

//封装奖励发放跨链交易信息
//func genStdWrapTx(cdc *amino.Codec, sendTx txs.ITx, priKey ed25519.PrivKeyEd25519, qosnonce, qscnonce int64, from types.Address, to []types.Address, amount []types.BigInt, causecode []string, causeStr []string) *txs.TxStd {
//	stx := genStdSendTx(cdc, sendTx, priKey, config.GetCLIContext().Config.QOSChainID, config.GetCLIContext().Config.QSCChainID, qosnonce)
//	//tx2 := txs.NewTxStd(sendTx, config.GetCLIContext().Config.QSCChainID, stx.MaxGas)
//	dispatchTx := NewDispatchAOE(stx, from, to, amount, causecode, causeStr, types.ZeroInt())
//	return genStdSendTx(cdc, dispatchTx, priKey, config.GetCLIContext().Config.QSCChainID, config.GetCLIContext().Config.QSCChainID, qscnonce)
//}

func wrapperResult(cdc *wire.Codec, msg *txs.TxStd) string {
	var cliCtx context.CLIContext
	cliCtx = *config.GetCLIContext().QSCCliContext
	hash, commitresult, err := utils.SendTx(cliCtx, cdc, msg)
	if err != nil {
		return common.NewErrorResult(COINS_SENDTX_ERR, 0, "", err.Error()).Marshal()
	}
	return common.NewSuccessResult(cdc, commitresult.Height, hash, hash).Marshal()
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

func GetCoins(cdc *wire.Codec, ctx *context.CLIContext, tx string) string {
	coins, err := jianqian.QueryCoins(cdc, ctx, tx)
	if err != nil {
		return common.NewErrorResult(COINS_QUERY_ERR, 0, "", err.Error()).Marshal()
	}
	if coins == nil || coins.Tx == "" {
		return common.NewErrorResult(COINS_QUERY_ERR, 0, "", fmt.Sprintf("query dispatch coins failure,%s not exist", tx)).Marshal()
	}
	return common.NewSuccessResult(cdc, 0, "", coins).Marshal()
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

