package kvstore

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/store"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/baseapp"
	go_amino "github.com/tendermint/go-amino"
)

type KVStub struct {
	baseapp.BaseContract
	KvTx KvstoreTx
}

func NewKVStub() KVStub {
	return KVStub{}
}

func (kv KVStub) StartX(base *baseapp.QstarsBaseApp) error{

	var mainStore = store.NewKVStoreKey("kv")
	var kvMapper = NewKvMapper(mainStore)
	base.Baseapp.RegisterMapper(kvMapper)

	if err := base.Baseapp.LoadLatestVersion(); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (kv KVStub) RegisterKVCdc(cdc *go_amino.Codec) {
	cdc.RegisterConcrete(&kv.KvTx, "kvstore/KvstoreTx", nil)
}

func (kv KVStub) ResultNotify(ctx context.Context, txQcpResult interface{}) *types.Result{
	return nil
}
