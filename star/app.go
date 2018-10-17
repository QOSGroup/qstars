package star

import (
	"github.com/QOSGroup/qstars/baseapp"
	"github.com/QOSGroup/qstars/wire"

	"github.com/QOSGroup/qstars/x/kvstore"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	"io"
	"os"
	dbm "github.com/tendermint/tendermint/libs/db"
)

func NewApp(log.Logger, dbm.DB, io.Writer) abci.Application {
	//cfg := ctx.Config
	//rootDir := cfg.RootDir
	rootDir := os.ExpandEnv("$HOME/.qstarsd")
	app := baseapp.NewAPP(rootDir)
	app.Register(kvstore.NewKVStub())

	app.Start()
	return app.Baseapp
}

func MakeCodec() *wire.Codec {
	cdc := wire.NewCodec()

	return cdc
}
