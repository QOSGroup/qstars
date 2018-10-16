package kvstore

import (
	"fmt"
	"github.com/QOSGroup/qbase/store"
	"github.com/QOSGroup/qstars/baseapp"
	go_amino "github.com/tendermint/go-amino"
	"os"
)

type KVStub struct {
	baseapp.BaseContract
	KvTx KvstoreTx
}

func NewKVStub() KVStub {
	return KVStub{}
}

func (kv KVStub) StartX(base *baseapp.QstarsBaseApp) {

	var mainStore = store.NewKVStoreKey("kv")
	var kvMapper = NewKvMapper(mainStore)
	base.Baseapp.RegisterMapper(kvMapper)

	if err := base.Baseapp.LoadLatestVersion(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (kv KVStub) RegisterKVCdc(cdc *go_amino.Codec) {

	cdc.RegisterConcrete(&kv.KvTx, "kvstore/KvstoreTx", nil)
}
