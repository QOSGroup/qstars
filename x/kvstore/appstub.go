package kvstore

import (
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qstars/baseapp"
	"github.com/QOSGroup/qbase/store"
	"fmt"
	"os"
	go_amino "github.com/tendermint/go-amino"
	"path/filepath"
	dbm "github.com/tendermint/tendermint/libs/db"
)

type KVStub struct{
	baseapp.BaseContract
	KvTx KvstoreTx
}

func NewKVStub()(KVStub){
	return KVStub{}
}

func (kv KVStub) StartX(base *baseapp.QstarsBaseApp) {

	db, err := dbm.NewGoLevelDB("kvstore", filepath.Join(base.RootDir, "data"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	base.Baseapp = baseabci.NewBaseApp("kvstore", base.Logger, db, nil)


	base.Baseapp.RegisterAccountProto(func() account.Account {
		return &account.BaseAccount{}
	})

	var mainStore = store.NewKVStoreKey("kv")
	var kvMapper = NewKvMapper(mainStore)
	base.Baseapp.RegisterMapper(kvMapper)

	if err := base.Baseapp.LoadLatestVersion(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (kv KVStub) RegisterKVCdc(cdc *go_amino.Codec){

	cdc.RegisterConcrete(&kv.KvTx, "kvstore/KvstoreTx", nil)
}