package jsdk

import (
	"encoding/json"
	"github.com/QOSGroup/qstars/slim"
	"github.com/QOSGroup/qstars/star"
	"github.com/QOSGroup/qstars/x/quzhuan/common"
)



var CDC = star.MakeCodec()


func init(){

}


func CommonHandler(url,chainid,funcName, privateKey, args string) string {
	result := common.CommHandler(CDC, url,chainid,funcName, privateKey, args)
	return result
}



func ParameterFormat(parameter string) string {
	strs:=[]string{parameter}
	result,_:=json.Marshal(strs)
	return string(result)
}



func CreateAccount() string {
	info:=slim.AccountCreate("")
	result,_:=json.Marshal(info)
	return string(result)
}


