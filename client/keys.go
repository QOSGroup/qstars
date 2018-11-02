package client

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	dbm "github.com/tendermint/tendermint/libs/db"
)

// GetKeyBase initializes a keybase based on the given db.
// The KeyBase manages all activity requiring access to a key.
func GetKeyBase(db dbm.DB) keys.Keybase {
	keybase := keys.New(
		db,
	)
	return keybase
}

