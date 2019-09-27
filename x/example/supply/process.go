package supply

import (
	"fmt"
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/txs"
	qbasetypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/client/context"
	"github.com/QOSGroup/qstars/client/utils"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/utility"
	"github.com/QOSGroup/qstars/wire"
	"github.com/QOSGroup/qstars/x/bank/tx"
	"github.com/QOSGroup/qstars/x/common"
	"github.com/prometheus/common/log"
	"strconv"
	"strings"
	"time"

	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

type SendResult struct {
	Hash   string `json:"hash"`
	Error  string `json:"error"`
	Code   string `json:"code"`
	Result string `json:"result"`
	Heigth string `json:"heigth"`
}

type SendOptions struct {
	fee string
	gas int64
}

func fee(f string) func(*SendOptions) {
	return func(opt *SendOptions) {
		opt.fee = f
	}
}

func gas(g int64) func(*SendOptions) {
	return func(opt *SendOptions) {
		opt.gas = g
	}
}

func NewSendOptions(opts ...func(*SendOptions)) *SendOptions {
	sopt := new(SendOptions)

	if opts != nil {
		for _, opt := range opts {
			opt(sopt)
		}
	}

	return sopt
}

func queryQSCAccount(cdc *wire.Codec, key []byte) (*SendResult, error, int64) {
	chainid := config.GetCLIContext().Config.QSCChainID
	accqsc, errqsc := config.GetCLIContext().QSCCliContext.GetAccount(key, cdc)

	if errqsc != nil {
		fmt.Println(errqsc.Error())
	}
	fmt.Println("---------" + chainid)
	result := &SendResult{}
	if errqsc != nil {
		if errqsc.Error() != context.ACCOUNT_NOT_EXIST {
			result.Hash = ""
			result.Error = errqsc.Error()
			result.Code = "1"
			return result, errqsc, 0
		}
	}
	qscnonce := int64(0)
	if accqsc == nil {
		qscnonce = 0
	} else {
		qscnonce = int64(accqsc.Nonce)
	}
	return nil, nil, qscnonce
}

// Send 支持一次多种币 coins.Len() == 1;
func Send(cdc *wire.Codec, privatestr string, to qbasetypes.Address, coins types.Coins, id string, sopt *SendOptions) (*SendResult, error) {
	_, addrben32, priv := utility.PubAddrRetrievalFromAmino(privatestr, cdc)
	from, err := types.AccAddressFromBech32(addrben32)
	if err != nil {
		return nil,err
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

	var ccs []qbasetypes.BaseCoin
	for _, coin := range coins {
		ccs = append(ccs, qbasetypes.BaseCoin{
			Name:   coin.Denom,
			Amount: qbasetypes.NewInt(coin.Amount.Int64()),
		})
	}

	t := tx.NewTransfer(from, to, ccs)
	directTOQOS := config.GetCLIContext().Config.DirectTOQOS
	var msg *txs.TxStd
	var cliCtx context.CLIContext
	if directTOQOS == true {
		cliCtx = *config.GetCLIContext().QOSCliContext

		msg = genStdSendTx(cdc, t, priv, config.GetCLIContext().Config.QOSChainID, config.GetCLIContext().Config.QOSChainID, qosnonce)
	} else {
		cliCtx = *config.GetCLIContext().QSCCliContext

		result, err1, qscnonce := queryQSCAccount(cdc, key)
		if result != nil {
			return result, err1
		}
		qscnonce++

		order := &OrderTx{Address: from, OrderTo: to, OrderAmount: ccs[0].Amount, Gas: qbasetypes.NewInt(0)}

		msg = genStdWrapTx(cdc, t, priv, config.GetCLIContext().Config.QOSChainID, config.GetCLIContext().Config.QSCChainID, qosnonce, qscnonce, order)
	}
	response, commitresult, err := utils.SendTx(cliCtx, cdc, msg)

	result := &SendResult{}
	if err != nil {
		result.Hash = ""
		result.Error = err.Error()
		result.Code = "1"
		return result, nil
	}
	result.Hash = response
	height := strconv.FormatInt(commitresult.Height, 10)
	result.Heigth = height
	waittime, err := strconv.Atoi(config.GetCLIContext().Config.WaitingForQosResult)
	if err != nil {
		panic("WaitingForQosResult should be able to convert to integer." + err.Error())
	}

	if directTOQOS == false {
		counter := 0
		for {
			if counter >= waittime {
				fmt.Println("time out")
				result.Error = "time out"
				break
			}
			resultstr, err := fetchResult(cdc, height, commitresult.Hash.String())
			if err != nil {
				fmt.Println("get result error:" + err.Error())
				result.Error = err.Error()
			}
			if resultstr == "BankStub" {
				result.Error = ""
				result.Result = resultstr
				result.Code = "-1"
			} else if resultstr != "" && resultstr != "-1" {
				fmt.Printf("get result:[%+v]\n", resultstr)
				rs := []rune(resultstr)
				index1 := strings.Index(resultstr, " ")

				result.Error = ""
				result.Result = string(rs[index1+1:])
				result.Code = string(rs[:index1])
				break
			}
			time.Sleep(500 * time.Millisecond)
			counter++
		}
	}
	return result, nil
}

func fetchResult(cdc *wire.Codec, heigth1 string, tx1 string) (string, error) {
	qstarskey := GetResultKey(heigth1, tx1)
	d, err := config.GetCLIContext().QSCCliContext.QueryStore([]byte(qstarskey), common.QSCResultMapperName)
	log.Infof("QueryStore: %+v, %+v\n", d, err)
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

//add the string input chainid
func genStdSendTx(cdc *amino.Codec, sendTx txs.ITx, priKey ed25519.PrivKeyEd25519, tochainid string, fromchainid string, nonce int64) *txs.TxStd {
	gas := qbasetypes.NewInt(int64(config.MaxGas))
	stx := txs.NewTxStd(sendTx, tochainid, gas)
	signature, _ := stx.SignTx(priKey, nonce, fromchainid, tochainid)
	stx.Signature = []txs.Signature{txs.Signature{
		Pubkey:    priKey.PubKey(),
		Signature: signature,
		Nonce:     nonce,
	}}

	return stx
}

//add the string input chainid
func genStdWrapTx(cdc *amino.Codec, sendTx txs.ITx, priKey ed25519.PrivKeyEd25519, tochainid string, fromchainid string, qosnonce int64, qscnonce int64, order *OrderTx) *txs.TxStd {
	stx := genStdSendTx(cdc, sendTx, priKey, tochainid, fromchainid, qosnonce)
	tx2 := txs.NewTxStd(nil, fromchainid, stx.MaxGas)
	order.Wrapper = stx
	tx2.ITxs[0] = order
	signature, _ := tx2.SignTx(priKey, qscnonce, fromchainid, fromchainid)
	tx2.Signature = []txs.Signature{txs.Signature{
		Pubkey:    priKey.PubKey(),
		Signature: signature,
		Nonce:     qscnonce,
	}}

	return tx2
}
