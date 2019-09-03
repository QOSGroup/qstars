package common

import (
	"encoding/json"
	"fmt"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qstars/baseapp"
	qstarstypes "github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/utility"
	"github.com/QOSGroup/qstars/wire"
	"github.com/QOSGroup/qstars/x/bank"
	"github.com/QOSGroup/qstars/x/jianqian/advertisers"
	"github.com/QOSGroup/qstars/x/jianqian/article"
	"github.com/QOSGroup/qstars/x/jianqian/auction"
	"github.com/QOSGroup/qstars/x/jianqian/buyad"
	"github.com/QOSGroup/qstars/x/jianqian/coins"
	"github.com/QOSGroup/qstars/x/jianqian/comm"
	"github.com/QOSGroup/qstars/x/jianqian/investad"
	"github.com/QOSGroup/qstars/x/jianqian/recharge"
	"github.com/QOSGroup/qstars/x/kvstore"
	"reflect"

	"testing"
)

/**
	init a qstar chain instance
   Because golang doesn't support reflect a name of structure to struct instance
   developer has to modify source code below to add their defined transaction
   there is only one place to register developer's transaction
*/

var _ baseapp.BaseXTransaction = bank.BankStub{}
var _ baseapp.BaseXTransaction = kvstore.KVStub{}

var _ baseapp.BaseXTransaction = investad.InvestadStub{}
var _ baseapp.BaseXTransaction = buyad.BuyadStub{}

var _ baseapp.BaseXTransaction = coins.CoinsStub{}
var _ baseapp.BaseXTransaction = article.AricleStub{}
var _ baseapp.BaseXTransaction = advertisers.AdvertisersStub{}
var _ baseapp.BaseXTransaction = auction.AuctionStub{}
var _ baseapp.BaseXTransaction = recharge.RechargeStub{}

var _ baseapp.BaseXTransaction = comm.JianQianStub{}
var _ baseapp.BaseXTransaction = QuZhuanStub{}




func init() {
	registerType((*bank.BankStub)(nil))
	registerType((*kvstore.KVStub)(nil))
	registerType((*investad.InvestadStub)(nil))
	registerType((*buyad.BuyadStub)(nil))
	registerType((*coins.CoinsStub)(nil))
	registerType((*article.AricleStub)(nil))
	registerType((*advertisers.AdvertisersStub)(nil))
	registerType((*auction.AuctionStub)(nil))
	registerType((*recharge.RechargeStub)(nil))
	registerType((*comm.JianQianStub)(nil))
	registerType((*QuZhuanStub)(nil))


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

func MakeCodec() *wire.Codec {
	cdc := baseabci.MakeQBaseCodec()
	for k, _ := range typeRegistry {
		txs, err := newStruct(k)
		if err == false {
			panic("reflect transaction is error.")
		}
		t := txs.(baseapp.BaseXTransaction)
		t.RegisterCdc(cdc)
	}
	//kvstore.NewKVStub().RegisterKVCdc(cdc)
	//bank.NewBankStub().RegisterKVCdc(cdc)
	return cdc
}

func TestCommHandler(m *testing.T) {
	cdc := MakeCodec()

	url:="localhost:26657"
	chainid:="test-chain-GEbNwW"
	funcNanme:="scenesReward"
	privateKey:="hATXd/o2bP1ZHkrb3JgYPeMA4EyYT3jKcMoft4mNVJsh5/u3xJUnYnAddKpeuDWsEbzuG16qOJms/h1IZMoLhw=="
	parameter:=ttt("{\"ScenesId\":\"001\",\"Rewards\":[{\"UserId\":\"13388888888\",\"Amount\":\"888\"},{\"UserId\":\"13399999999\",\"Amount\":\"999\"}]}")
	aaaa:=CommHandler(cdc,url,chainid,funcNanme,privateKey,parameter)


	//BroadcastTransferTxToQSC(cdc,aaaa,url)





	fmt.Println(aaaa)
	//
	//m.Run()
}


func ttt(parameter string) string {
	strs:=[]string{parameter}
	result,_:=json.Marshal(strs)
	return string(result)
}




func TestRpcQueryAccount(m *testing.T){
	cdc := MakeCodec()
	_, addrben32, _ := utility.PubAddrRetrievalFromAmino("hATXd/o2bP1ZHkrb3JgYPeMA4EyYT3jKcMoft4mNVJsh5/u3xJUnYnAddKpeuDWsEbzuG16qOJms/h1IZMoLhw==", cdc)
	from, _ := qstarstypes.AccAddressFromBech32(addrben32)

	account,err:=RpcQueryAccount(cdc,"localhost:26657",from)

	if err!=nil{
		fmt.Println(err)
	}

	if account==nil{
		fmt.Println("account is nul")
	}else {

		fmt.Println(account.Nonce)
	}

}