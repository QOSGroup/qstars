// Copyright 2018 The QOS Authors

package bank

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
	"github.com/pkg/errors"
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

func TxJsonSend(cdc *wire.Codec, txb []byte) (*SendResult, error) {
	ts := new(txs.TxStd)
	err := cdc.UnmarshalJSON(txb, ts)
	if err != nil {
		return nil, err
	}

	cliCtx := *config.GetCLIContext().QOSCliContext
	response, commitresult, err := utils.SendTx(cliCtx, cdc, ts)
	result := &SendResult{}
	if err != nil {
		result.Error = err.Error()
		return result, nil
	}

	result.Hash = response
	height := strconv.FormatInt(commitresult.Height, 10)
	result.Heigth = height
	return result, nil
}

func TxSend(cdc *wire.Codec, txb []byte) (*SendResult, error) {
	ts := new(txs.TxStd)
	err := cdc.UnmarshalJSON(txb, ts)
	if err != nil {
		return nil, err
	}

	cliCtx := *config.GetCLIContext().QOSCliContext
	response, commitresult, err := utils.SendTx(cliCtx, cdc, ts)
	result := &SendResult{}
	if err != nil {
		//result.Error = err.Error()
		//return result, nil
		return nil, err
	}

	result.Hash = response
	height := strconv.FormatInt(commitresult.Height, 10)
	result.Heigth = height
	return result, nil
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
func Send(cdc *wire.Codec, fromstr string, to qbasetypes.Address, coins types.Coins, sopt *SendOptions) (*SendResult, error) {
	if coins.Len() == 0 {
		return nil, errors.New("coins不能为空")
	}

	_, addrben32, priv := utility.PubAddrRetrievalFromAmino(fromstr, cdc)
	from, err := types.AccAddressFromBech32(addrben32)
	key := account.AddressStoreKey(from)
	if err != nil {
		return nil, err
	}

	directTOQOS := config.GetCLIContext().Config.DirectTOQOS
	var cliCtx context.CLIContext
	if directTOQOS == true {
		cliCtx = *config.GetCLIContext().QOSCliContext
	} else {
		cliCtx = *config.GetCLIContext().QSCCliContext
	}

	acc, err := config.GetCLIContext().QOSCliContext.GetAccount(key, cdc)
	if err != nil {
		return nil, err
	}

	var ccs []qbasetypes.BaseCoin
	for _, coin := range coins {
		ccs = append(ccs, qbasetypes.BaseCoin{
			Name:   coin.Denom,
			Amount: qbasetypes.NewInt(coin.Amount.Int64()),
		})
	}

	var qcoins types.Coins
	for _, qsc := range acc.QSCs {
		amount := qsc.Amount
		qcoins = append(qcoins, types.NewCoin(qsc.Name, types.NewInt(amount.Int64())))
	}
	qcoins = append(qcoins, types.NewCoin("qos", types.NewInt(acc.QOS.Int64())))

	qcoins.Sort()
	coins.Sort()
	if !qcoins.IsGTE(coins) {
		return nil, errors.Errorf("Address %s doesn't have enough coins to pay for this transaction.", from)
	}

	var nn int64
	nn = int64(acc.Nonce)
	nn++

	t := tx.NewTransfer(from, to, ccs)
	var msg *txs.TxStd
	if directTOQOS == true {
		msg = genStdSendTx(cdc, t, priv, config.GetCLIContext().Config.QOSChainID, config.GetCLIContext().Config.QOSChainID, nn)
	} else {
		result, err1, qscnonce := queryQSCAccount(cdc, key)
		if result != nil {
			return result, err1
		}
		qscnonce++
		msg = genStdWrapTx(cdc, t, priv, config.GetCLIContext().Config.QOSChainID, config.GetCLIContext().Config.QSCChainID, nn, qscnonce)
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

// Send 暂时只支持一次只转一种币 coins.Len() == 1; add chainid string input
func Approve(cdc *wire.Codec, command string, fromstr string, tostr string, coins types.Coins, chainid string,
	sopt *SendOptions) (*SendResult, error) {
	if command != "cancel" {
		if coins.Len() == 0 {
			return nil, errors.New("coins不能为空")
		}
	}

	fmt.Printf("---command:%s, from:%s, to:%s\n", command, fromstr, tostr)
	var priv ed25519.PrivKeyEd25519
	var from, to qbasetypes.Address
	var nonce int64
	var err error
	var key []byte
	if command == "use" {
		from, err = types.AccAddressFromBech32(fromstr)
		if err != nil {
			return nil, err
		}

		var addrben32 string
		_, addrben32, priv = utility.PubAddrRetrievalFromAmino(tostr, cdc)
		to, err = types.AccAddressFromBech32(addrben32)
		if err != nil {
			return nil, err
		}

		key = account.AddressStoreKey(to)
		if err != nil {
			return nil, err
		}
		acc, err := config.GetCLIContext().QOSCliContext.GetAccount(key, cdc)
		if err != nil {
			return nil, err
		}
		nonce = int64(acc.Nonce)
		nonce++

	} else {
		var addrben32 string
		_, addrben32, priv = utility.PubAddrRetrievalFromAmino(fromstr, cdc)
		from, err = types.AccAddressFromBech32(addrben32)
		if err != nil {
			return nil, err
		}

		to, err = types.AccAddressFromBech32(tostr)
		if err != nil {
			return nil, err
		}

		key = account.AddressStoreKey(from)
		if err != nil {
			return nil, err
		}

		acc, err := config.GetCLIContext().QOSCliContext.GetAccount(key, cdc)
		if err != nil {
			return nil, err
		}
		var qcoins types.Coins
		for _, qsc := range acc.QSCs {
			amount := qsc.Amount
			qcoins = append(qcoins, types.NewCoin(qsc.Name, types.NewInt(amount.Int64())))
		}
		qcoins = append(qcoins, types.NewCoin("qos", types.NewInt(acc.QOS.Int64())))
		if command != "cancel" {
			if !qcoins.IsGTE(coins) {
				return nil, errors.Errorf("Address %s doesn't have enough coins to pay for this transaction.", from)
			}
		}

		nonce = int64(acc.Nonce)
		nonce++
	}

	var ccs []qbasetypes.BaseCoin
	for _, coin := range coins {
		ccs = append(ccs, qbasetypes.BaseCoin{
			Name:   coin.Denom,
			Amount: qbasetypes.NewInt(coin.Amount.Int64()),
		})
	}

	atx := tx.NewApproveTx(from, to)
	var t txs.ITx
	switch command {
	case "create":
		t = atx.Create(ccs)
	case "increase":
		t = atx.Increase(ccs)
	case "decrease":
		t = atx.Decrease(ccs)
	case "use":
		t = atx.Use(ccs)
	case "cancel":
		t = atx.Cancel()
	default:
		return nil, errors.New("command not support")
	}

	var msg *txs.TxStd
	directTOQOS := config.GetCLIContext().Config.DirectTOQOS
	if directTOQOS == true {
		msg = genStdSendTx(cdc, t, priv, config.GetCLIContext().Config.QOSChainID, config.GetCLIContext().Config.QOSChainID, nonce)
	} else {
		result, err1, qscnonce := queryQSCAccount(cdc, key)
		if result != nil {
			return result, err1
		}
		qscnonce++
		msg = genStdWrapTx(cdc, t, priv, config.GetCLIContext().Config.QSCChainID, config.GetCLIContext().Config.QOSChainID, nonce, qscnonce)
	}

	var cliCtx context.CLIContext
	if directTOQOS == true {
		cliCtx = *config.GetCLIContext().QOSCliContext
	} else {
		cliCtx = *config.GetCLIContext().QSCCliContext
	}
	response, commitresult, err := utils.SendTx(cliCtx, cdc, msg)

	result := &SendResult{}
	result.Hash = response
	height := strconv.FormatInt(commitresult.Height, 10)
	result.Heigth = height
	if directTOQOS == false {
		counter := 0
		for {
			if counter >= 10 {
				fmt.Println("time out")
				result.Error = "time out"
				break
			}
			resultstr, err := fetchResult(cdc, height, commitresult.Hash.String())
			if err != nil {
				fmt.Println("get result error:" + err.Error())
				result.Error = err.Error()
			}
			if resultstr != "" && resultstr != "-1" {
				fmt.Printf("get result:[%+v]\n", resultstr)
				result.Result = resultstr
				break
			}
			time.Sleep(500 * time.Millisecond)
			counter++
		}
	}
	return result, nil
}

func fetchResult(cdc *wire.Codec, heigth1 string, tx1 string) (string, error) {

	// TODO qbase还没实现
	//qstarskey := "heigth:" + heigth1 + ",hash:" + tx1
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
	gas := qbasetypes.NewInt(int64(0))
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
func genStdWrapTx(cdc *amino.Codec, sendTx txs.ITx, priKey ed25519.PrivKeyEd25519, tochainid string, fromchainid string, qosnonce int64, qscnonce int64) *txs.TxStd {
	stx := genStdSendTx(cdc, sendTx, priKey, tochainid, fromchainid, qosnonce)
	tx2 := txs.NewTxStd(nil, fromchainid, stx.MaxGas)
	tx2.ITx = NewWrapperSendTx(stx)
	signature, _ := tx2.SignTx(priKey, qscnonce, fromchainid, fromchainid)
	tx2.Signature = []txs.Signature{txs.Signature{
		Pubkey:    priKey.PubKey(),
		Signature: signature,
		Nonce:     qscnonce,
	}}

	return tx2
}
