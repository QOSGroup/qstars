package kvstore

import (
	"github.com/QOSGroup/qstars/baseapp"
	"github.com/QOSGroup/qbase/store"
	"fmt"
	"os"
	go_amino "github.com/tendermint/go-amino"
)

type KVStub struct{
	baseapp.BaseContract
	KvTx KvstoreTx
}

func NewKVStub()(KVStub){
	return KVStub{}
}
func (kv KVStub) StartX(base *baseapp.QstarsBaseApp) {

	var mainStore = store.NewKVStoreKey("kv")
	var kvMapper = NewKvMapper(mainStore)
	base.Baseapp.RegisterSeedMapper(kvMapper)

	if err := base.Baseapp.LoadLatestVersion(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (kv KVStub) RegisterKVCdc(cdc *go_amino.Codec){

	cdc.RegisterConcrete(&kv.KvTx, "kvstore/KvstoreTx", nil)
}