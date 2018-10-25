// Copyright 2018 The QOS Authors

package auth

import (
	qosacc "github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/wire"
)

// QueryAccount query account by addr
func QueryAccount(cdc *wire.Codec, addr string) (*qosacc.QOSAccount, error) {
	key, err := types.AccAddressFromBech32(addr)
	if err != nil {
		return nil, err
	}

	cliCtx := config.GetCLIContext().QOSCliContext

	acc, err := cliCtx.GetAccount(key)
	if err != nil {
		return nil, err
	}

	return acc, nil
}
