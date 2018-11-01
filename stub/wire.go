package stub

import (
	"github.com/QOSGroup/qstars/star"
	"github.com/QOSGroup/qstars/wire"
)

// Register concrete types on wire codec
var cmCdc *wire.Codec

func init() {
	cmCdc = star.MakeCodec()
}
