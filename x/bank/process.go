// Copyright 2018 The QOS Authors

package bank

import (
	qbaseaccount "github.com/QOSGroup/qbase/account"
	qosaccount "github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qstars/client/utils"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/utility"
	"github.com/QOSGroup/qstars/wire"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto"
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

func Send(cdc *wire.Codec, fromstr string, to types.AccAddress, coins types.Coins, sopt *SendOptions) (*SendResult, error) {
	_, addrben32 := utility.PubAddrRetrieval(fromstr)

	from, err := types.AccAddressFromBech32(addrben32)
	if err != nil {
		return nil, err
	}
	cliCtx:= *config.GetCLIContext().QSCCliContext

	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(&ed25519.PubKeyEd25519{}, "ed25519.PubKeyEd25519", nil)
	cdc.RegisterInterface((*qbaseaccount.Account)(nil), nil)
	cdc.RegisterConcrete(&qosaccount.QOSAccount{}, "qbase/account/QOSAccount", nil)

	account, err := cliCtx.GetAccount(from)
	if err != nil {
		return nil, err
	}

	var qcoins types.Coins
	for _, qsc := range account.QscList {
		amount := qsc.Amount
		qcoins = append(qcoins, types.NewCoin(qsc.Name, types.NewInt(amount.Int64())))
	}

	if !qcoins.IsGTE(coins) {
		return nil, errors.Errorf("Address %s doesn't have enough coins to pay for this transaction.", from)
	}

	// build and sign the transaction, then broadcast to Tendermint
	msg := BuildMsg(from, to, coins, cdc)

	var priv ed25519.PrivKeyEd25519
	bz := utility.Decbase64(fromstr)
	copy(priv[:], bz)
	response, err := utils.SendTx(cliCtx, cdc, msg, priv)
	result := &SendResult{}
	result.Hash = response

	return result, nil
}
