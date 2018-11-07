package star

import (
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qstars/baseapp"
	"github.com/QOSGroup/qstars/wire"
	"github.com/QOSGroup/qstars/x/bank"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/x/kvstore"
	"io"
	"os"
	"reflect"

	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

/**
	init a qstar chain instance
   Because golang doesn't support reflect a name of structure to struct instance
   developer has to modify source code below to add their defined transaction
   there is only one place to register developer's transaction
 */

var _ baseapp.BaseXTransaction = bank.BankStub{}
var _ baseapp.BaseXTransaction = kvstore.KVStub{}

func init() {
	registerType((*bank.BankStub)(nil))
	registerType((*kvstore.KVStub)(nil))
}

/**
	startup a qstar chain instance
 */
func NewApp(logger log.Logger, db dbm.DB, io io.Writer) (abci.Application) {
	//cfg := ctx.Config
	//rootDir := cfg.RootDir
	rootDir := os.ExpandEnv("$HOME/.qstarsd")
	sconf, err := config.Init(rootDir + "/config/qstarsconf.toml",rootDir)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	app, err := baseapp.NewAPP(sconf, MakeCodec())
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	//TODO
	for k, _ := range typeRegistry {
		txs, _ := newStruct(k)
		t := txs.(baseapp.BaseXTransaction)
		app.Register(t)
	}
	//app.Register(kvstore.NewKVStub())
	//app.Register(bank.NewBankStub())

	err = app.Start()
	if err!=nil{
		logger.Error(err.Error())
		return nil
	}
	return app.Baseapp
}

/*
Both client and server can get a well-setting cdc via the funcation
 */
func MakeCodec() *wire.Codec {
	cdc := baseabci.MakeQBaseCodec()
	for k, _ := range typeRegistry {
		txs, err := newStruct(k)
		if err==false {
			panic("reflect transaction is error.")
		}
		t := txs.(baseapp.BaseXTransaction)
		t.RegisterCdc(cdc)
	}
	//kvstore.NewKVStub().RegisterKVCdc(cdc)
	//bank.NewBankStub().RegisterKVCdc(cdc)
	return cdc
}

//---------------------------------------------------------------------------
var typeRegistry = make(map[string]reflect.Type)

func registerType(elem interface{}) {
	t := reflect.TypeOf(elem).Elem()
	typeRegistry[t.Name()] = t
}

func newStruct(name string) (interface{}, bool) {
	elem, ok := typeRegistry[name]
	if !ok {
		return nil, false
	}
	return reflect.New(elem).Elem().Interface(), true
}
