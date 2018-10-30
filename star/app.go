package star

import (
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qstars/baseapp"
	"github.com/QOSGroup/qstars/wire"
	"github.com/QOSGroup/qstars/x/bank"
	"io"
	"os"

	"github.com/QOSGroup/qstars/x/kvstore"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func NewApp(log.Logger, dbm.DB, io.Writer) abci.Application {
	//cfg := ctx.Config
	//rootDir := cfg.RootDir
	rootDir := os.ExpandEnv("$HOME/.qstarsd")
	app,err := baseapp.NewAPP(rootDir,MakeCodec())
	if err != nil{
		return nil
	}
	app.Register(kvstore.NewKVStub())
	app.Register(bank.NewBankStub())

	app.Start()
	return app.Baseapp
}

func MakeCodec() *wire.Codec {

	cdc := baseabci.MakeQBaseCodec()
	//cdc.RegisterConcrete(&bctypes.AppAccount{}, "basecoin/AppAccount", nil)
	return cdc
}
