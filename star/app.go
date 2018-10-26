package star

import (
	"github.com/QOSGroup/qbase/example/basecoin/types"
	"github.com/QOSGroup/qstars/baseapp"
	"github.com/QOSGroup/qstars/wire"
	"github.com/QOSGroup/qstars/x/bank"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/QOSGroup/qbase/txs"


	"io"
	"os"

	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qstars/x/kvstore"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func NewApp(log.Logger, dbm.DB, io.Writer) abci.Application {
	//cfg := ctx.Config
	//rootDir := cfg.RootDir
	rootDir := os.ExpandEnv("$HOME/.qstarsd")
	app := baseapp.NewAPP(rootDir)
	app.Register(kvstore.NewKVStub())
	app.Register(bank.NewBankStub())

	app.Start()
	return app.Baseapp
}

func MakeCodec() *wire.Codec {
	cdc := wire.NewCodec()

	cdc.RegisterConcrete(&types.AppAccount{}, "basecoin/AppAccount", nil)

	cryptoAmino.RegisterAmino(cdc)

	cdc.RegisterInterface((*account.Account)(nil), nil)

	cdc.RegisterConcrete(&bank.SendTx{}, "basecoin/SendTx",nil)
	kvstore.NewKVStub().RegisterKVCdc(cdc)

	txs.RegisterCodec(cdc)

	//cdc.RegisterConcrete(&qosacc.QOSAccount{}, "qbase/account/QOSAccount", nil)

	return cdc
}
