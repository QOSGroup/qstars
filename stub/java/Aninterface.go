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
	GetIPfromInput("192.168.1.23:1317")
	//acc, err := testQuery("http://localhost:1317", "address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355")
	//fmt.Printf("---acc:%+v, err:%+v\n", acc, err)
	out := QSCQueryAccount("address1k0m8ucnqug974maa6g36zw7g2wvfd4sug6uxay")
	fmt.Println(out)
	//QSCKVStoreSet("13", "Melon", "0xa328891040ae9b773bcd30005235f99a8d62df03a89e4f690f9fa03abb1bf22715fc9ca05613f2d8061492e9f8149510b5b67d340d199ff24f34c85dbbbd7e0df780e9a6cc", "test-chain-Ky10Zg")
	//kvout := QSCKVStoreGet("13")
	//fmt.Println(kvout)
	//AccountRecover("vague success fresh check remove banner music snap jelly medal bring mix eagle seat cash off winter mean comic turn always teach tiny wagon")
	//puba := PubAddrRetrieval("oyiJEECum3c7zTAAUjX5mo1i3wOonk9pD5+gOrsb8icV/JygVhPy2AYUkun4FJUQtbZ9NA0Zn/JPNMhdu71+DfeA6abM")
	//fmt.Println(puba)
	//transout := QSCtransferSend("address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355", "22qos", "0xa328891040ae9b773bcd30005235f99a8d62df03a89e4f690f9fa03abb1bf22715fc9ca05613f2d8061492e9f8149510b5b67d340d199ff24f34c85dbbbd7e0df780e9a6cc", "qos-testapp")
	//transoutb64 := QSCtransferSend("address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355", "22qos", "oyiJEECum3c7zTAAUjX5mo1i3wOonk9pD5+gOrsb8icV/JygVhPy2AYUkun4FJUQtbZ9NA0Zn/JPNMhdu71+DfeA6abM", "qos-testapp")
	//fmt.Println(transoutb64)
}
