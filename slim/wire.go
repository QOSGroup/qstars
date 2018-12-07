package slim

import (
	"github.com/QOSGroup/qstars/slim/funcInlocal/ed25519local"
	"github.com/tendermint/go-amino"
)

//
//import (
//	"github.com/QOSGroup/qbase/account"
//	"github.com/QOSGroup/qbase/baseabci"
//	"github.com/QOSGroup/qbase/context"
//	"github.com/QOSGroup/qbase/mapper"
//	"github.com/QOSGroup/qbase/txs"
//	"github.com/tendermint/go-amino"
//	"github.com/tendermint/tendermint/crypto"
//	"github.com/tendermint/tendermint/crypto/ed25519"
//	cmn "github.com/tendermint/tendermint/libs/common"
//	"log"
//	"reflect"
//)
//
type Codec = amino.Codec

func NewCodec() *Codec {
	cdc := amino.NewCodec()
	return cdc
}

var Cdc *Codec

func init() {
	cdc := NewCodec()
	RegisterAmino(cdc)
	RegisterCodec(cdc)
	Cdc = cdc.Seal()
}

// RegisterAmino registers all crypto related types in the given (amino) codec.
func RegisterAmino(cdc *amino.Codec) {
	cdc.RegisterInterface((*ed25519local.PubKey)(nil), nil)
	cdc.RegisterConcrete(ed25519local.PubKeyEd25519{},
		ed25519local.Ed25519PubKeyAminoRoute, nil)

	cdc.RegisterInterface((*ed25519local.PrivKey)(nil), nil)
	cdc.RegisterConcrete(ed25519local.PrivKeyEd25519{},
		ed25519local.Ed25519PrivKeyAminoRoute, nil)
}

func RegisterCodec(cdc *amino.Codec) {
	cdc.RegisterConcrete(&Signature{}, "qbase/txs/signature", nil)
	cdc.RegisterConcrete(&TxStd{}, "qbase/txs/stdtx", nil)
	cdc.RegisterInterface((*ITx)(nil), nil)
	cdc.RegisterConcrete(&TxTransfer{}, "qos/txs/TxTransfer", nil)
	cdc.RegisterConcrete(&QOSAccount{}, "qbase/account/QOSAccount", nil)
	cdc.RegisterConcrete(&BaseAccount{}, "qbase/account/BaseAccount", nil)
	cdc.RegisterConcrete(&InvestTx{}, "qstars/InvestTx", nil)
}

// amino codec to marshal/unmarshal
//var typeRegistry = make(map[string]reflect.Type)
//var Cdc *amino.Codec
//
//type ABCICodeType uint32
//type Tags cmn.KVPairs
//
//type QstarsBaseApp struct {
//	Transactions    BaseXTransaction
//	Baseapp         *baseabci.BaseApp
//	TransactionList []BaseXTransaction
//	Logger          log.Logger
//	RootDir         string
//}
//type BaseXTransaction interface {
//	mapper.IMapper
//	RegisterCdc(cdc *amino.Codec)
//	StartX(base *QstarsBaseApp) error
//}
//type Result struct {
//
//	// Code is the response code, is stored back on the chain.
//	Code ABCICodeType
//
//	// Data is any data returned from the app.
//	Data []byte
//
//	// Log is just debug information. NOTE: nondeterministic.
//	Log string
//
//	// GasWanted is the maximum units of work we allow this tx to perform.
//	GasWanted int64
//
//	// GasUsed is the amount of gas actually consumed. NOTE: unimplemented
//	GasUsed int64
//
//	// Tx fee amount and denom.
//	FeeAmount int64
//	FeeDenom  string
//
//	// Tags are used for transaction indexing and pubsub.
//	Tags Tags
//}
//
//func MakeCodec() *amino.Codec {
//	cdc := MakeQBaseCodec()
//	for k, _ := range typeRegistry {
//		txs, err := newStruct(k)
//		if err == false {
//			panic("reflect transaction is error.")
//		}
//		t := txs.(BaseXTransaction)
//		t.RegisterCdc(cdc)
//	}
//	//kvstore.NewKVStub().RegisterKVCdc(cdc)
//	//bank.NewBankStub().RegisterKVCdc(cdc)
//	return cdc
//}
//
////
//func newStruct(name string) (interface{}, bool) {
//	elem, ok := typeRegistry[name]
//	if !ok {
//		return nil, false
//	}
//	return reflect.New(elem).Elem().Interface(), true
//}
//
//func init() {
//	Cdc = MakeCodec()
//}
//
//func MakeQBaseCodec() *amino.Codec {
//
//	var cdc = amino.NewCodec()
//	//RegisterAmino(cdc)
//	RegisterCodec(cdc)
//
//	return cdc
//}
//
//func RegisterCodec(cdc *amino.Codec) {
//	//txs.RegisterCodec(cdc)
//	//account.RegisterCodec(cdc)
//}

//func RegisterAmino(cdc *amino.Codec) {
//	// These are all written here instead of
//	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
//	cdc.RegisterConcrete(ed25519.PubKeyEd25519{},
//		"tendermint/PubKeyEd25519", nil)
//
//	cdc.RegisterInterface((*crypto.PrivKey)(nil), nil)
//	cdc.RegisterConcrete(ed25519.PrivKeyEd25519{},
//		"tendermint/PrivKeyEd25519", nil)
//
//}
