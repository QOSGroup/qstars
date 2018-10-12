package star

import (
	"github.com/QOSGroup/qstars/wire"
	"github.com/QOSGroup/qstars/baseapp"

	"github.com/QOSGroup/qstars/x/kvstore"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	abci "github.com/tendermint/tendermint/abci/types"
	"io"
)

func NewApp(logger log.Logger, db dbm.DB, storeTracer io.Writer) abci.Application {
	app:= baseapp.NewAPP()
	app.Register(kvstore.NewKVStub())

	app.Start()
	return app.Baseapp
}

func MakeCodec() *wire.Codec {
	cdc := wire.NewCodec()

	return cdc
}