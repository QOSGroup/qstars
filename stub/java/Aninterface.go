package main

import "C"
import (
	"fmt"
	"github.com/QOSGroup/qstars/stub"
)

// ----------------------------------------------------------------------------
// source code for so file generation with "go build " command, e.g.go build -o awesome.so -buildmode=c-shared awesome.go
// ----------------------------------------------------------------------------

//export AccountCreate
//func AccountCreate() *C.char {
//	output := stub.AccountCreateStr()
//	return C.CString(output)
//}

func AccountCreate() string {
	output := stub.AccountCreateStr()
	return output
	//stub.AccountCreateStr(w)
}

//export QSCKVStoreSet
//func QSCKVStoreSet(k, v, privkey, chain *C.char) int {
//	output := stub.QSCKVStoreSetPost(C.GoString(k), C.GoString(v), C.GoString(privkey), C.GoString(chain))
//	return output
//}

func QSCKVStoreSet(k, v, privkey, chain string) string {
	output := stub.QSCKVStoreSetPost(k, v, privkey, chain)
	return output
}

//export QSCKVStoreGet
//func QSCKVStoreGet(k *C.char) *C.char {
//	output := stub.QSCKVStoreGetQuery(C.GoString(k))
//	return C.CString(output)
//}

func QSCKVStoreGet(k string) string {
	output := stub.QSCKVStoreGetQuery(k)
	return output
}

//export QSCQueryAccount
//func QSCQueryAccount(addr *C.char) *C.char {
//	output := stub.QSCQueryAccountGet(C.GoString(addr))
//	return C.CString(output)
//}

func QSCQueryAccount(addr string) string {
	output := stub.QSCQueryAccountGet(addr)
	return output
}

//export AccountRecover
//func AccountRecover(mncode *C.char) *C.char {
//	output := stub.AccountRecoverStr(C.GoString(mncode))
//	return C.CString(output)
//}

func AccountRecover(mncode string) string {
	output := stub.AccountRecoverStr(mncode)
	return output
}

//export GetIPfromInput
//func GetIPfromInput(ip *C.char) {
//	stub.GetIPfrom(C.GoString(ip))
//}

func GetIPfromInput(ip string) {
	stub.GetIPfrom(ip)
}

//export PubAddrRetrieval
//func PubAddrRetrieval(priv *C.char) *C.char {
//	output := stub.PubAddrRetrievalStr(C.GoString(priv))
//	return C.CString(output)
//}

func PubAddrRetrieval(priv string) string {
	//	fmt.Println("Please input host including IP and port for initialization on Qstar deamon:")
	output := stub.PubAddrRetrievalStr(priv)
	return output
}

//export QSCtransferSend
//func QSCtransferSend(addrto, coinstr, privkey *C.char) *C.char {
//	output := stub.QSCtransferSendStr(C.GoString(addrto), C.GoString(coinstr), C.GoString(privkey))
//	return C.CString(output)
//}

//for QSCtransferSend
func QSCtransferSend(addrto, coinstr, privkey, chainid string) string {
	output := stub.QSCtransferSendStr(addrto, coinstr, privkey, chainid)
	return output
}

func main() {
	//AccountCreate()
	GetIPfromInput("localhost:1317")
	//acc, err := testQuery("http://localhost:1317", "address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355")
	//fmt.Printf("---acc:%+v, err:%+v\n", acc, err)
	//out := QSCQueryAccount("address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355")
	//fmt.Println(out)
	//QSCKVStoreSet("13", "Melon", "0xa328891040ae9b773bcd30005235f99a8d62df03a89e4f690f9fa03abb1bf22715fc9ca05613f2d8061492e9f8149510b5b67d340d199ff24f34c85dbbbd7e0df780e9a6cc", "test-chain-Ky10Zg")
	kvout := QSCKVStoreGet("13")
	fmt.Println(kvout)
	//AccountRecover("celery quick meat flight garden video adjust like rose process fly leaf series general vast explain rocket rail phrase sing trash drum success cannon")
	//puba := PubAddrRetrieval("0xa3288910400f8f271b2df5df818d267b5d87ea70aa25908748f67de4ed2f3e68b12b07f436483c84704d005d9b8064eb1546c4699d8b386bf285aaf18c8212f85dce28a29e")
	//fmt.Println(puba)
	//transout := QSCtransferSend("address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355", "222qos", "0xa328891040ae9b773bcd30005235f99a8d62df03a89e4f690f9fa03abb1bf22715fc9ca05613f2d8061492e9f8149510b5b67d340d199ff24f34c85dbbbd7e0df780e9a6cc", "qos-testapp")
	//fmt.Println(transout)
}
