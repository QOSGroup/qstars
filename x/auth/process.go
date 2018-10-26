// Copyright 2018 The QOS Authors

package auth

import (
	"github.com/QOSGroup/qbase/account"
	bctypes "github.com/QOSGroup/qbase/example/basecoin/types"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/wire"
)

// QueryAccount query account by addr
func QueryAccount(cdc *wire.Codec, addr string) (*bctypes.AppAccount, error) {
	address, err := types.GetAddrFromBech32(addr)
	if err != nil {
		return nil, err
	}
	key := account.AddressStoreKey(address)

	cliCtx := config.GetCLIContext().QOSCliContext

	acc, err := cliCtx.GetAccount(key,cdc)
	if err != nil {
		return nil, err
	}

	return acc, nil
}
