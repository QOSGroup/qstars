// Copyright 2018 The QOS Authors

package bank

import (
	"github.com/QOSGroup/qbase/example/basecoin/tx"
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qstars/client/context"
	"github.com/QOSGroup/qstars/client/utils"

	qbasetypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/utility"
	"github.com/QOSGroup/qstars/wire"
	"github.com/pkg/errors"

	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

type SendResult struct {
	Hash string `json:"hash"`
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

func Send(cdc *wire.Codec, fromstr string, to qbasetypes.Address, coins types.Coins, sopt *SendOptions) (*SendResult, error) {
	_, addrben32, priv:= utility.PubAddrRetrieval(fromstr,cdc)
	from, err := types.AccAddressFromBech32(addrben32)
	key := account.AddressStoreKey(from)
	if err != nil {
		return nil, err
	}

	directTOQOS := config.GetCLIContext().Config.DirectTOQOS
	var cliCtx context.CLIContext
	if directTOQOS==true{
		cliCtx = *config.GetCLIContext().QOSCliContext
	}else {
		cliCtx = *config.GetCLIContext().QSCCliContext
	}

	account, err := config.GetCLIContext().QOSCliContext.GetAccount(key,cdc)

	if err != nil {
		return nil, err
	}
	var cc qbasetypes.BaseCoin
	var qcoins types.Coins
	for _, qsc := range account.Coins {
		amount := qsc.Amount
		qcoins = append(qcoins, types.NewCoin(qsc.Name, types.NewInt(amount.Int64())))

		//TODO-------------------------
		if !amount.IsZero() {
			mount := qbasetypes.NewInt(100)
			cc = qbasetypes.BaseCoin{
				Name:   qsc.Name,
				Amount: mount,
			}
		}
	}

	if !qcoins.IsGTE(coins) {
		return nil, errors.Errorf("Address %s doesn't have enough coins to pay for this transaction.", from)
	}

	var nn int64
	nn = int64(account.Nonce)
	//if directTOQOS==true {
	//	nn = int64(account.Nonce)
	//}else {
	//	qscaccount, err := config.GetCLIContext().QSCCliContext.GetAccount(key,cdc)
	//	if err != nil{
	//		if err.Error()=="Account is not exsit." {
	//			nn = int64(0)
	//		}else{
	//			return nil,err
	//		}
	//	}else{
	//		nn = int64(qscaccount.Nonce)
	//	}
	//}
	nn++
	var msg *txs.TxStd
	if directTOQOS==true {
		msg = genStdSendTx(cdc, from, to, cc, priv, nn)
	}else{
		msg = genStdWrapTx(cdc, from, to, cc, priv, nn)
	}
	response, err := utils.SendTx(cliCtx, cdc,msg)

	result := &SendResult{}
	result.Hash = response

	return result, nil
}

func genStdSendTx(cdc *amino.Codec, sender qbasetypes.Address, receiver qbasetypes.Address, coin qbasetypes.BaseCoin,
	priKey ed25519.PrivKeyEd25519, nonce int64) *txs.TxStd {
	sendTx := tx.NewSendTx(sender, receiver, coin)
	gas := qbasetypes.NewInt(int64(0))
	tx := txs.NewTxStd(&sendTx, config.GetCLIContext().Config.QOSChainID, gas)
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
	priKey ed25519.PrivKeyEd25519, nonce int64 ) *txs.TxStd {
	sendTx := tx.NewSendTx(sender, receiver, coin)
	gas := qbasetypes.NewInt(int64(0))
	tx := txs.NewTxStd(&sendTx, config.GetCLIContext().Config.QOSChainID, gas)
	//priHex, _ := hex.DecodeString(senderPriHex[2:])
	//var priKey ed25519.PrivKeyEd25519
	//cdc.MustUnmarshalBinaryBare(priHex, &priKey)
	signature, _ := tx.SignTx(priKey, nonce)
	tx.Signature = []txs.Signature{txs.Signature{
		Pubkey:    priKey.PubKey(),
		Signature: signature,
		Nonce:     nonce,
	}}


	tx2 := txs.NewTxStd(&sendTx, config.GetCLIContext().Config.QSCChainID, gas)
	tx2.ITx = NewWrapperSendTx(tx)
	return tx2
}
