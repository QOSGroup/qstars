// Copyright 2018 The HSB Authors

// Package pkg comments for pkg auth
// auth ...
package auth

import (
	"github.com/QOSGroup/qbase/account"
	qosacc "github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qstars/client/context"
	"github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/wire"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

// QueryAccount query account by addr
func QueryAccount(cdc *wire.Codec, addr string) (*qosacc.QOSAccount, error) {
	key, err := types.AccAddressFromBech32(addr)
	if err != nil {
		return nil, err
	}

	cliCtx := context.NewOQSCLIContext().WithCodec(cdc)

	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(&ed25519.PubKeyEd25519{}, "ed25519.PubKeyEd25519", nil)
	cdc.RegisterInterface((*account.Account)(nil), nil)
	cdc.RegisterConcrete(&qosacc.QOSAccount{}, "qbase/account/QOSAccount", nil)

	acc, err := cliCtx.GetAccount(key)
	if err != nil {
		return nil, err
	}

	return acc, nil
}
