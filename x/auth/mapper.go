package auth

import (
	sdk "github.com/QOSGroup/qstars/types"
		)

var globalAccountNumberKey = []byte("globalAccountNumber")

// Turn an address to key used to get it from the account store
func AddressStoreKey(addr sdk.AccAddress) []byte {
	return append([]byte("account:"), addr.Bytes()...)
}

