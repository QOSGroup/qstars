// Copyright 2018 The QOS Authors

package auth

import (
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/types"
	qostype "github.com/QOSGroup/qos/types"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/wire"
)

// QueryAccount query account by addr
func QueryAccount(cdc *wire.Codec, addr string) (*qostype.QOSAccount, error) {
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

//Query QSCAccount by addr
func QSCQueryAccount(cdc *wire.Codec, addr string) (*qostype.QOSAccount, error) {
	address, err := types.GetAddrFromBech32(addr)
	if err != nil {
		return nil, err
	}

	key := account.AddressStoreKey(address)

	cliCtx := config.GetCLIContext().QSCCliContext

	acc, err := cliCtx.GetAccount(key, cdc)
	if err != nil {
		return nil, err
	}

	return acc, nil
}
