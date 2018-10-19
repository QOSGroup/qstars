package bank

import "github.com/QOSGroup/qbase/types"

import (
	"fmt"
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/qcp"
	"github.com/QOSGroup/qbase/store"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qstars/baseapp"
	"github.com/QOSGroup/qstars/x/kvstore"
	go_amino "github.com/tendermint/go-amino"
)

type KvstoreTx struct {
	baseapp.BaseContract
	Key   []byte
	Value []byte
	Bytes []byte
}

func NewKvstoreTx(key []byte, value []byte) KvstoreTx {
	return KvstoreTx{
		Key:   key,
		Value: value,
	}
}

func (kv KvstoreTx) StartX(base *baseapp.QstarsBaseApp) error{

	var mainStore = store.NewKVStoreKey("kv")
	var kvMapper = kvstore.NewKvMapper(mainStore)
	base.Baseapp.RegisterMapper(kvMapper)

	if err := base.Baseapp.LoadLatestVersion(); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (kv KvstoreTx) RegisterKVCdc(cdc *go_amino.Codec) {

	txs.RegisterCodec(cdc)
	cdc.RegisterConcrete(&kvstore.KvstoreTx{}, "kvstore/main/kvstoretx", nil)
}

func (kv KvstoreTx) ValidateData() bool {
	if len(kv.Key) < 0 {
		return false
	}
	return true
}

func (kv KvstoreTx) Exec(ctx context.Context) (result types.Result, crossTxQcps *txs.TxQcp) {

	logger := ctx.Logger()
	kvMapper := ctx.Mapper(KvMapperName).(*KvMapper)
	qcpMapper := ctx.Mapper(qcp.QcpMapperName).(*qcp.QcpMapper)
	accMapper := ctx.Mapper(account.AccountMapperName).(*account.AccountMapper)

	key := string(kv.Key)
	value := kvMapper.GetKey(key)

	logger.Info("qcpMapper", qcpMapper)
	logger.Info("accMapper", accMapper)

	logger.Info("origin is: ", key, "=", value)

	kvMapper.SaveKV(key, string(kv.Value))

	value = kvMapper.GetKey(key)

	logger.Info("after is: ", key, value)

	clearKey := "lllllll"

	store := kvMapper.GetStore()
	store.Set([]byte(clearKey), []byte("11111"))

	logger.Info("clear value: %s", store.Get([]byte(clearKey)))

	store.Delete([]byte(clearKey))

	return
}

func (kv KvstoreTx) GetSigner() []types.Address {
	return nil
}

func (kv KvstoreTx) CalcGas() types.BigInt {
	return types.ZeroInt()
}

func (kv KvstoreTx) GetGasPayer() types.Address {
	return types.Address{}
}

func (kv KvstoreTx) GetSignData() []byte {
	signData := make([]byte, len(kv.Key)+len(kv.Value)+len(kv.Bytes))
	signData = append(signData, kv.Key...)
	signData = append(signData, kv.Value...)
	signData = append(signData, kv.Bytes...)

	return signData
}
