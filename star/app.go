package star

import (
	"github.com/QOSGroup/qstars/baseapp"
	"github.com/QOSGroup/qstars/wire"

	"github.com/QOSGroup/qstars/x/kvstore"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	"io"
)

func NewApp(logger log.Logger, storeTracer io.Writer, rootDir string) abci.Application {
	app:= baseapp.NewAPP(rootDir)
	app.Register(kvstore.NewKVStub())

	app.Start()
	return app.Baseapp
}

func MakeCodec() *wire.Codec {
	cdc := wire.NewCodec()

	return cdc
}