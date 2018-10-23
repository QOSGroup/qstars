package main

import "C"
import (
	"fmt"
	"github.com/QOSGroup/qstars/stub"
)

// ----------------------------------------------------------------------------
// source code for so file generation with command, e.g. go build go build -o awesome.so -buildmode=c-shared awesome.go
// ----------------------------------------------------------------------------

//export AccountCreate
func AccountCreate() *C.char {
	output := stub.AccountCreateStr()
	return C.CString(output)
}

//export QSCKVStoreSet
func QSCKVStoreSet(k, v, privkey, chain *C.char) int {
	output := stub.QSCKVStoreSetPost(C.GoString(k), C.GoString(v), C.GoString(privkey), C.GoString(chain))
	return output
}

//normal function for Golang testing
//func QSCKVStoreSet(k, v, privkey, chain string) int {
//	output := stub.QSCKVStoreSetPost(k, v, privkey, chain)
//	return output
//}

//export QSCKVStoreGet
func QSCKVStoreGet(ul *C.char) *C.char {
	output := stub.QSCKVStoreGetQuery(C.GoString(ul))
	return C.CString(output)
}

//normal function for Golang testing
//func QSCKVStoreGet(ul string) string {
//	output := stub.QSCKVStoreGetQuery(ul)
//	return output
//}

//export QSCQueryAccount
func QSCQueryAccount(ul *C.char) *C.char {
	output := stub.QSCQueryAccountGet(C.GoString(ul))
	return C.CString(output)
}

//normal function for testing
//func QSCQueryAccount(ul string) string {
//	output := stub.QSCQueryAccountGet(ul)
//	return output
//}

//export QSCMintCoin
func QSCMintCoin() {
	fmt.Println("this is QSCMintCoin function interface")
}

//export QSCtransfer
func QSCtransfer(ul, a, privkey, chain, ac, seq, g *C.char) *C.char {
	output := stub.QSCtransferPost(C.GoString(ul), C.GoString(a), C.GoString(privkey), C.GoString(chain), C.GoString(ac), C.GoString(seq), C.GoString(g))
	return C.CString(output)
}

//export QOStoQSCtransfer
func QOStoQSCtransfer() {
	fmt.Println("this is QOStoQSCtransfer function interface")
}

//export QSCtoQOStransfer
func QSCtoQOStransfer() {
	fmt.Println("this is QSCtoQOStransfer function interface")
}

//export AccountRecover
func AccountRecover(mncode *C.char) *C.char {
	output := stub.AccountRecoverStr(mncode)
	return C.CString(output)
}

//export GetIPfromInput
func GetIPfromInput(ip *C.char) {
//	fmt.Println("Please input host including IP and port for initialization on Qstar deamon:")
	stub.GetIPfrom(ip)
}

func main() {
//	AccountCreate()
//	GetIPfromInput("127.0.0.1:1317")
//	out := QSCQueryAccount("address1an4rky8h6c7jgwk92arg2ms3q9ehask87yl4x2")
//	fmt.Println(out)
//	QSCKVStoreSet("2", "Bob", "lEMsVbO4nCbAdQkr9hyTG15IaGvBIq1BFcNt4XeSeF9uo80srafrM25SVpfS1naE8G7MpYhcoQ9Wu1yFIl3ZEw", "test-chain-Ky10Zg")
//	kvout := QSCKVStoreGet("2")
//	fmt.Println(kvout)
	//AccountRecover("clean axis column history legend mosquito worry magic exotic beef layer glue cannon craft worry decide slice soft hockey tennis lottery spatial segment minute")
}
