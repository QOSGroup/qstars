// Copyright 2018 The QOS Authors

package auth

import (
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/types"
	qosaccount "github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/wire"
)

// QueryAccount query account by addr
func QueryAccount(cdc *wire.Codec, addr string) (*qosaccount.QOSAccount, error) {
	address, err := types.GetAddrFromBech32(addr)
	if err != nil {
		return nil, err
	}
	key := account.AddressStoreKey(address)

	cliCtx := config.GetCLIContext().QOSCliContext

	acc, err := cliCtx.GetAccount(key, cdc)
	if err != nil {
		return nil, err
	}

	return acc, nil
}
