package bank

import (
	"github.com/QOSGroup/qstars/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
}

var msgCdc = wire.NewCodec()

func init() {
	RegisterWire(msgCdc)
}
