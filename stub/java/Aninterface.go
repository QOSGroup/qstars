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
func AccountCreate() *C.char {
	output := stub.AccountCreateStr()
	return C.CString(output)
}

//export QSCKVStoreSet
func QSCKVStoreSet(k, v, privkey, chain *C.char) int {
	output := stub.QSCKVStoreSetPost(C.GoString(k), C.GoString(v), C.GoString(privkey), C.GoString(chain))
	return output
}


//export QSCKVStoreGet
func QSCKVStoreGet(k *C.char) *C.char {
	output := stub.QSCKVStoreGetQuery(C.GoString(k))
	return C.CString(output)
}

//export QSCQueryAccount
func QSCQueryAccount(addr *C.char) *C.char {
	output := stub.QSCQueryAccountGet(C.GoString(addr))
	return C.CString(output)
}


//export QSCMintCoin
func QSCMintCoin() {
	fmt.Println("this is QSCMintCoin function interface")
}

//export QSCtransfer
func QSCtransfer(addr, amount, privkey, chain, accountnumber, seq, gas *C.char) *C.char {
	output := stub.QSCtransferPost(C.GoString(addr), C.GoString(amount), C.GoString(privkey), C.GoString(chain), C.GoString(accountnumber), C.GoString(seq), C.GoString(gas))
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
	output := stub.AccountRecoverStr(C.GoString(mncode))
	return C.CString(output)
}

//export GetIPfromInput
func GetIPfromInput(ip *C.char) {
	stub.GetIPfrom(C.GoString(ip))
}

//export PubAddrRetrieval
func PubAddrRetrieval(priv *C.char) *C.char {
	output := stub.PubAddrRetrievalStr(C.GoString(priv))
	return C.CString(output)
}

func main() {
//	AccountCreate()
//	GetIPfromInput("127.0.0.1:1317")
//	out := QSCQueryAccount("address1an4rky8h6c7jgwk92arg2ms3q9ehask87yl4x2")
//	fmt.Println(out)
//	QSCKVStoreSet("2", "Bob", "lEMsVbO4nCbAdQkr9hyTG15IaGvBIq1BFcNt4XeSeF9uo80srafrM25SVpfS1naE8G7MpYhcoQ9Wu1yFIl3ZEw", "test-chain-Ky10Zg")
//	kvout := QSCKVStoreGet("2")
//	fmt.Println(kvout)
//	AccountRecover("clean axis column history legend mosquito worry magic exotic beef layer glue cannon craft worry decide slice soft hockey tennis lottery spatial segment minute")
//	puba:= PubAddrRetrieval("lEMsVbO4nCbAdQkr9hyTG15IaGvBIq1BFcNt4XeSeF9uo80srafrM25SVpfS1naE8G7MpYhcoQ9Wu1yFIl3ZEw")
//	fmt.Println(puba)
}
