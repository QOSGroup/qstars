// Copyright 2018 The QOS Authors

package bank

import (
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/txs"
	qbasetypes "github.com/QOSGroup/qbase/types"
	qostxs "github.com/QOSGroup/qos/txs"
	"github.com/QOSGroup/qstars/client/context"
	"github.com/QOSGroup/qstars/client/utils"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/utility"
	"github.com/QOSGroup/qstars/wire"
	"github.com/pkg/errors"
	"strconv"
	"time"

	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

type SendResult struct {
	Hash   string `json:"hash"`
	Error  string `json:"error"`
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

// Send 暂时只支持一次只转一种币 coins.Len() == 1
func Send(cdc *wire.Codec, fromstr string, to qbasetypes.Address, coins types.Coins, sopt *SendOptions) (*SendResult, error) {
	_, addrben32, priv := utility.PubAddrRetrieval(fromstr, cdc)

	// TODO 暂时只支持一次只转一种币
	if coins.Len() == 0 {
		return nil, errors.New("coins不能为空")
	}

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

	account, err := config.GetCLIContext().QOSCliContext.GetAccount(key, cdc)

	if err != nil {
		return nil, err
	}
	var cc qbasetypes.BaseCoin
	// TODO 暂时只支持一次只转一种币
	cc = qbasetypes.BaseCoin{
		Name:   coins[0].Denom,
		Amount: qbasetypes.NewInt(coins[0].Amount.Int64()),
	}

	var qcoins types.Coins
	for _, qsc := range account.QSCs {
		amount := qsc.Amount
		qcoins = append(qcoins, types.NewCoin(qsc.Name, types.NewInt(amount.Int64())))
	}

	if !qcoins.IsGTE(coins) {
		return nil, errors.Errorf("Address %s doesn't have enough coins to pay for this transaction.", from)
	}

	var nn int64
	nn = int64(account.Nonce)

	nn++
	var msg *txs.TxStd
	if directTOQOS == true {
		msg = genStdSendTx(cdc, from, to, cc, priv, nn)
	} else {
		msg = genStdWrapTx(cdc, from, to, cc, priv, nn)
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
			resultstr, err := fetchResult(height, commitresult.Hash.String())
			if err != nil {
				fmt.Println("get result error:" + err.Error())
				result.Error = err.Error()
			}
			if resultstr != "-1" {
				fmt.Println("get result")
				result.Result = resultstr
				break
			}
			time.Sleep(500 * time.Millisecond)
			counter++
		}
	}
	return result, nil
}

func fetchResult(heigth1 string, tx1 string) (string, error) {

	qstarskey := "heigth:" + heigth1 + ",hash:" + tx1
	res, err := config.GetCLIContext().QSCCliContext.QueryStore([]byte(qstarskey), QSCResultMapperName)
	re := string(res)
	return re, err
}
func newQOSTx(sender qbasetypes.Address, receiver qbasetypes.Address, coin qbasetypes.BaseCoin) *qostxs.TransferTx {
	sendTx := qostxs.TransferTx{}

	sendTx.Senders = append(sendTx.Senders,
		qostxs.TransItem{Address: sender, QOS: qbasetypes.NewInt(0), QSCs: qbasetypes.BaseCoins{&coin}})

	sendTx.Receivers = append(sendTx.Receivers,
		qostxs.TransItem{Address: receiver, QOS: qbasetypes.NewInt(0), QSCs: qbasetypes.BaseCoins{&coin}})

	return &sendTx
}

func genStdSendTx(cdc *amino.Codec, sender qbasetypes.Address, receiver qbasetypes.Address, coin qbasetypes.BaseCoin,
	priKey ed25519.PrivKeyEd25519, nonce int64) *txs.TxStd {
	sendTx := newQOSTx(sender, receiver, coin)
	gas := qbasetypes.NewInt(int64(0))
	tx := txs.NewTxStd(sendTx, config.GetCLIContext().Config.QOSChainID, gas)
	//priHex, _ := hex.DecodeString(senderPriHex[2:])
	//var priKey ed25519.PrivKeyEd25519
	//cdc.MustUnmarshalBinaryBare(priHex, &priKey)
	signature, _ := tx.SignTx(priKey, nonce)
	tx.Signature = []txs.Signature{txs.Signature{
		Pubkey:    priKey.PubKey(),
		Signature: signature,
		Nonce:     nonce,
	}}
	return tx
}

func genStdWrapTx(cdc *amino.Codec, sender qbasetypes.Address, receiver qbasetypes.Address, coin qbasetypes.BaseCoin,
	priKey ed25519.PrivKeyEd25519, nonce int64) *txs.TxStd {
	sendTx := newQOSTx(sender, receiver, coin)
	gas := qbasetypes.NewInt(int64(0))
	tx := txs.NewTxStd(sendTx, config.GetCLIContext().Config.QOSChainID, gas)
	//priHex, _ := hex.DecodeString(senderPriHex[2:])
	//var priKey ed25519.PrivKeyEd25519
	//cdc.MustUnmarshalBinaryBare(priHex, &priKey)
	signature, _ := tx.SignTx(priKey, nonce)
	tx.Signature = []txs.Signature{txs.Signature{
		Pubkey:    priKey.PubKey(),
		Signature: signature,
		Nonce:     nonce,
	}}

	tx2 := txs.NewTxStd(sendTx, config.GetCLIContext().Config.QSCChainID, gas)
	tx2.ITx = NewWrapperSendTx(tx)
	return tx2
}
