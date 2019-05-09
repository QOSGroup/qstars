package main

import (
	"encoding/json"
	"fmt"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/star"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"

	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

func main(){
	//addrs:="address1y9r4pjjnvkmpvw46de8tmwunw4nx4qnz2ax5ux1"
	//path := fmt.Sprintf("/store/%s/%s", "aoeaccount", "key")
	//output,err := Query(path,[]byte(addrs))
	//if err!=nil{
	//	fmt.Println(err.Error())
	//}else {
	//
	//	a,_:=json.Marshal(output)
	//	fmt.Println(string(a))
	//}


	str:="[\"a\",\"b\",\"c\",\"d\"]"

	var args []string

	json.Unmarshal([]byte(str),&args)



fmt.Println(args)




}

func Query(path string, key common.HexBytes) (res *types.BaseCoins, err error) {

	cdc := star.MakeCodec()


	RPC := rpcclient.NewHTTP("47.105.52.237:26657", "/websocket")

	opts := rpcclient.ABCIQueryOptions{
		Height: 0,
		Prove:  true,
	}
	result, err := RPC.ABCIQueryWithOptions(path, key, opts)
	if err != nil {
		return res, err
	}
	resp := result.Response
	if !resp.IsOK() {
		return res, errors.New("error query")
	}

	//fmt.Println("key",string(resp.GetKey()))
	//
	//fmt.Println("value",string(resp.GetValue()))

	var basecoin *types.BaseCoins



	//err=json.Unmarshal(resp.Value,&basecoin)
	err = cdc.UnmarshalBinaryBare(resp.Value, &basecoin)

	return basecoin, err
}