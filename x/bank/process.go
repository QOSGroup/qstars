// Copyright 2018 The QOS Authors

package bank

import (

	"github.com/QOSGroup/qbase/example/basecoin/tx"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qstars/client/utils"

	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/types"
	qbasetypes "github.com/QOSGroup/qbase/types"
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
	_, addrben32 := utility.PubAddrRetrieval(fromstr)

	from, err := types.AccAddressFromBech32(addrben32)
	if err != nil {
		return nil, err
	}
	cliCtx:= *config.GetCLIContext().QOSCliContext

	//cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	//cdc.RegisterConcrete(&ed25519.PubKeyEd25519{}, "ed25519.PubKeyEd25519", nil)
	//cdc.RegisterInterface((*qbaseaccount.Account)(nil), nil)
	//cdc.RegisterConcrete(&qosaccount.QOSAccount{}, "qbase/account/QOSAccount", nil)

	account, err := cliCtx.GetAccount(from,cdc)
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
				Name:   "qos",
				Amount: mount,
			}
		}
	}

	if !qcoins.IsGTE(coins) {
		return nil, errors.Errorf("Address %s doesn't have enough coins to pay for this transaction.", from)
	}

	// build and sign the transaction, then broadcast to Tendermint
	//msg := BuildMsg(from, to, coins, cdc)
	//
	var priv ed25519.PrivKeyEd25519
	bz := utility.Decbase64(fromstr)
	copy(priv[:], bz)


	nn := int64(account.Nonce)
	msg := genStdSendTx(cdc,from,to,cc,priv,nn,"qos")
	response, err := utils.SendTx(cliCtx, cdc,msg,priv)

	result := &SendResult{}
	result.Hash = response

	return result, nil
}

func genStdSendTx(cdc *amino.Codec, sender qbasetypes.Address, receiver qbasetypes.Address, coin qbasetypes.BaseCoin,
	priKey ed25519.PrivKeyEd25519, nonce int64, chainid string) *txs.TxStd {
	sendTx := tx.NewSendTx(sender, receiver, coin)
	gas := qbasetypes.NewInt(int64(0))
	tx := txs.NewTxStd(&sendTx, chainid, gas)
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
