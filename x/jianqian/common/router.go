package common

import (
	"errors"
	"github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/x/jianqian"
	"github.com/QOSGroup/qstars/x/jianqian/common/adve"
	qbasetypes "github.com/QOSGroup/qbase/types"

)

func GetStruct(funcName string, args []string) (RouterTx, error) {
	switch funcName {
	case AdvertisersTxFlag:
		coinsTx, err := getCoins(args[0], args[2], args[3], args[4])
		if err != nil {
			return nil, err
		}
		tx := adve.AdvertisersTx{coinsTx}
		return tx, nil
	case ArticleTxFlag:

	}
	return nil, errors.New(funcName + " funcName not support")
}

func getCoins(addr, cointype, changetype, amount string) (*jianqian.CoinsTx, error) {
	coinsTx := &jianqian.CoinsTx{}
	address, err := types.AccAddressFromBech32(addr)
	if err != nil {
		return nil, err
	}
	coins, ok := qbasetypes.NewIntFromString(amount)
	if !ok {
		return nil, errors.New("amount format error")
	}
	coinsTx.Address = address
	coinsTx.Cointype = cointype
	coinsTx.ChangeType = changetype
	coinsTx.Amount = coins
	return coinsTx, nil
}
