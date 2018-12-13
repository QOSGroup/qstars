package coins

import (
	"fmt"
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qstars/client/context"
	"github.com/QOSGroup/qstars/x/common"
	"github.com/QOSGroup/qstars/x/jianqian"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/client/utils"
	"github.com/QOSGroup/qstars/config"
	qstartypes "github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/utility"
	"github.com/QOSGroup/qstars/wire"
	"github.com/QOSGroup/qstars/x/jianqian/tx"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

//type SendResult struct {
//	Hash   string `json:"hash"`
//	Error  string `json:"error"`
//	Code   string `json:"code"`
//	Result string `json:"result"`
//	Heigth string `json:"heigth"`
//}
const (
	COINS_PARA_LEN_ERR = "101"   //参数长度不一致
	COINS_PRIV_ERR = "102"   //私钥获取地址错误
	COINS_SENDTX_ERR = "103" //交易出错
	COINS_FETCH_RESULT_ERR = "104" //查询跨链结果错误
	COINS_QUERY_ERR = "105" //查询跨链结果错误
)


//发放活动奖励 一转多
func DispatchSend(cdc *wire.Codec, ctx *config.CLIConfig, privkey string, to []types.Address, amount []types.BigInt, causecode []string, causeStr []string) string {
	tolen := len(to)
	//判断长度是否一致
	if tolen != len(amount) || tolen != len(causecode) || tolen != len(causeStr) {
		return common.NewErrorResult(COINS_PARA_LEN_ERR,"Array parameter length is inconsistent").Marshal()
	}

	_, addrben32, priv := utility.PubAddrRetrievalFromAmino(privkey, cdc)
	from, err := qstartypes.AccAddressFromBech32(addrben32)
	if err != nil {
		return common.NewErrorResult(COINS_PRIV_ERR,err.Error()).Marshal()
	}
	key := account.AddressStoreKey(from)
	var qosnonce int64 = 0
	acc, err := config.GetCLIContext().QOSCliContext.GetAccount(key, cdc)
	if err != nil {
		qosnonce = 0
	} else {
		qosnonce = int64(acc.Nonce)
	}
	qosnonce++
	fmt.Println("qosnonce",qosnonce)
	var ccs []types.BaseCoin
	for _, coin := range amount {
		ccs = append(ccs, types.BaseCoin{
			Name:   COINNAME,
			Amount: types.NewInt(coin.Int64()),
		})
	}
	transtx := tx.NewTransfer([]types.Address{from}, to, ccs)
	directTOQOS := config.GetCLIContext().Config.DirectTOQOS
	var msg *txs.TxStd

	if directTOQOS == true {
		//直接连接公链
		msg=genStdSendTx(cdc, transtx, priv, config.GetCLIContext().Config.QOSChainID, config.GetCLIContext().Config.QOSChainID, qosnonce)

	}else{
		//走跨链
		var qscnonce int64 = 0
		qscacc, err := config.GetCLIContext().QSCCliContext.GetAccount(key, cdc)
		if err != nil {
			qscnonce = 0
		} else {
			qscnonce = int64(qscacc.Nonce)
		}
		qscnonce++
		fmt.Println("qscnonce",qscnonce)

		msg= genStdWrapTx(cdc, transtx, priv, qosnonce,qscnonce, from, to, amount, causecode, causeStr)
	}
	//	chainid := ctx.QOSChainID
	//chainid := config.GetCLIContext().Config.QSCChainID
	return wrapperResult(cdc, msg,directTOQOS)
}

//封装公链交易信息
func genStdSendTx(cdc *amino.Codec, sendTx txs.ITx, priKey ed25519.PrivKeyEd25519, tochainid,fromchainid string,  nonce int64) *txs.TxStd {
	gas := types.NewInt(int64(0))
	stx := txs.NewTxStd(sendTx, tochainid, gas)

	signature, _ := stx.SignTx(priKey, nonce, fromchainid)
	stx.Signature = []txs.Signature{txs.Signature{
		Pubkey:    priKey.PubKey(),
		Signature: signature,
		Nonce:     nonce,
	}}
	return stx
}

//封装奖励发放跨链交易信息
func genStdWrapTx(cdc *amino.Codec, sendTx txs.ITx, priKey ed25519.PrivKeyEd25519,  qosnonce,qscnonce int64, from types.Address, to []types.Address, amount []types.BigInt, causecode []string, causeStr []string) *txs.TxStd {
	stx := genStdSendTx(cdc, sendTx, priKey,  config.GetCLIContext().Config.QOSChainID,config.GetCLIContext().Config.QSCChainID, qosnonce)
	//tx2 := txs.NewTxStd(sendTx, config.GetCLIContext().Config.QSCChainID, stx.MaxGas)
	dispatchTx := NewDispatchAOE(stx, from, to, amount, causecode, causeStr, types.ZeroInt())
	return genStdSendTx(cdc, dispatchTx, priKey,config.GetCLIContext().Config.QSCChainID,config.GetCLIContext().Config.QSCChainID,  qscnonce)
}

func wrapperResult(cdc *wire.Codec, msg *txs.TxStd,directTOQOS bool) string {
	var cliCtx context.CLIContext
	if directTOQOS == true {
		cliCtx = *config.GetCLIContext().QOSCliContext
	} else {
		cliCtx = *config.GetCLIContext().QSCCliContext
	}
	hash, commitresult, err := utils.SendTx(cliCtx, cdc, msg)
	if err!=nil{
		return common.NewErrorResult(COINS_SENDTX_ERR,err.Error()).Marshal()
	}
	height := strconv.FormatInt(commitresult.Height, 10)
	waittime, err := strconv.Atoi(config.GetCLIContext().Config.WaitingForQosResult)
	if err != nil {
		waittime=30
	}
	if directTOQOS == false {
		counter := 0
		for {
			if counter >= waittime {
				log.Println("time out")
				return common.ResultQOSTimeoutError().Marshal()
			}
			resultstr, err := fetchResult(cdc, height, commitresult.Hash.String())
			log.Printf("fetchResult result:%s, err:%+v\n", resultstr, err)
			if err != nil {
				log.Printf("fetchResult error:%s\n", err.Error())
				return common.NewErrorResult(COINS_FETCH_RESULT_ERR,err.Error()).Marshal()
			}

			if resultstr != "" && resultstr != (CoinsStub{}).Name() {
				log.Printf("fetchResult result:[%+v]\n", resultstr)
				rs := []rune(resultstr)
				index1 := strings.Index(resultstr, " ")
				reason := string(rs[index1+1:])
				code := string(rs[:index1])
				result:=common.NewErrorResult(code,reason)
				result.Hash=hash
				result.Height=commitresult.Height
				return result.Marshal()
			}
			time.Sleep(500 * time.Millisecond)
			counter++
		}
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
		return "{Code:\"1\",Reason:\"Parameter cannot be empty \"}"
	}
	addrs := strings.Split(address, "|")
	addlen := len(addrs)
	cois := strings.Split(coins, "|")
	codes := strings.Split(causecodes, "|")
	cstrs := strings.Split(causestrings, "|")

	if addlen != len(cois) || addlen != len(codes) || addlen != len(cstrs) {
		return "{Code:\"2\",Reason:\"Parameter lengths are not equal \"}"
	}
	if address == "" || coins == "" || causecodes == "" || causestrings == "" {
		return "{Code:\"1\",Reason:\"Parameter cannot be empty \"}"
	}
	amounts := make([]types.BigInt, len(cois))
	for i, coinsv := range cois {
		if amou, ok := types.NewIntFromString(coinsv); ok {
			amounts[i] = amou
		} else {
			return "{Code:\"2\",Reason:\"amount format error \"}"
		}
	}
	toaddrss := make([]types.Address, addlen)
	for i, addrsv := range addrs {
		to, err := qstartypes.AccAddressFromBech32(addrsv)
		if err != nil {
			return "{Code:\"2\",Reason:\"address format error \"}"
		}
		toaddrss[i] = to
	}
	//cdc := star.MakeCodec()
	privkey := tx.GetConfig().Dappowner
	return  DispatchSend(cdc, ctx, privkey, toaddrss, amounts, codes, cstrs)
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
	if err!=nil{
		return common.NewErrorResult(COINS_QUERY_ERR,err.Error()).Marshal()
	}
	return common.NewSuccessResult(cdc, 0, "", coins).Marshal()
}
