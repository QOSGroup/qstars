package baseapp

import (
	"fmt"
	"github.com/QOSGroup/qbase/store"
	"github.com/QOSGroup/qstars/wire"
	"testing"
	go_amino "github.com/tendermint/go-amino"
	"github.com/QOSGroup/qbase/mapper"

	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qbase/context"
)

// TODO update
func TestInitCmd(t *testing.T) {
	InitApp()
	cdc := wire.NewCodec()
	app,_:=NewAPP("",cdc)
	mock := new(MockABCI)
	mock.RegisterKVCdc(cdc)
	app.Register(mock)
	app.Start()
}

type MockABCI struct{
	Cdc *go_amino.Codec
}

func (mock *MockABCI) MapperName() string {
	panic("implement me")
}

func (mock *MockABCI )RegisterKVCdc(cdc *go_amino.Codec){
	mock.Cdc = cdc
}

func (mock MockABCI )StartX(base *QstarsBaseApp) error{
	fmt.Println("StartX")
	return nil
}

func (mock MockABCI )Name() string{
	return "mock"
}
func (mock MockABCI )GetStoreKey() store.StoreKey{
	return nil
}
func (mock MockABCI )GetCodec() *go_amino.Codec{
	return mock.Cdc

}
func (mock *MockABCI )SetCodec(cdc *go_amino.Codec){
	mock.Cdc = cdc

}
func (mock MockABCI )Get(key []byte, ptr interface{}) (exsits bool){
	return false
}
func (mock *MockABCI )Set(key []byte, val interface{}){

}
func (mock *MockABCI )SetStore(store store.KVStore){

}
func (mock MockABCI )GetStore() store.KVStore{
	return nil
}
func (mock MockABCI )Copy() mapper.IMapper{
	return nil
}

func (mock *MockABCI )ResultNotify(ctx context.Context, txQcpResult interface{}) *types.Result{
	return &types.Result{}
}

func (mock MockABCI )Del(key []byte){

}