package kvstore

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/baseapp"
	"github.com/QOSGroup/qstars/x/common"
	go_amino "github.com/tendermint/go-amino"
)

type KVStub struct {
	baseapp.BaseXTransaction
	KvTx KvstoreTx
}

func NewKVStub() KVStub {
	return KVStub{}
}

func (kv KVStub) StartX(base *baseapp.QstarsBaseApp) error {

	//var mainStore = store.NewKVStoreKey("kv")
	var kvMapper = common.NewKvMapper(KvMapperName)
	base.Baseapp.RegisterMapper(kvMapper)

	return nil
}
func (kv KVStub) EndBlockNotify(ctx context.Context){

}

func (kv KVStub) RegisterCdc(cdc *go_amino.Codec) {
	cdc.RegisterConcrete(&kv.KvTx, "kvstore/KvstoreTx", nil)
}

func (kv KVStub) ResultNotify(ctx context.Context, txQcpResult interface{}) *types.Result {
	return nil
}
