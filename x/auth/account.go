package auth

import (
	qos "github.com/QOSGroup/qos/account"
	types "github.com/QOSGroup/qstars/types"
)

// Account is a standard account using a sequence number for replay protection
// and a pubkey for authentication.
type QAccount interface {
	GetQOSAccount() qos.QOSAccount
	GetCoins() types.Coins
}

// AccountDecoder unmarshals account bytes
type AccountDecoder func(accountBytes []byte) (QAccount, error)

//-----------------------------------------------------------
// QStarsAccount

type QStarsAccount struct {
	QosAccount qos.QOSAccount
	QCoins     types.Coins
}

func (acc QStarsAccount) GetQOSAccount() qos.QOSAccount {
	return acc.QosAccount
}

func (acc QStarsAccount) GetCoins() types.Coins {
	return acc.QCoins
}
