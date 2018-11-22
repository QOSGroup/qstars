package main

import "C"
import (
	"github.com/QOSGroup/qstars/slim"
)

//export AccountCreate
func AccountCreate() *C.char {
	output := slim.AccountCreateStr()
	return C.CString(output)
}

//func AccountCreate() string {
//	output := slim.AccountCreateStr()
//	return output
//}

//export QSCKVStoreSet
func QSCKVStoreSet(k, v, privkey, chain *C.char) *C.char {
	output := slim.QSCKVStoreSetPost(C.GoString(k), C.GoString(v), C.GoString(privkey), C.GoString(chain))
	return C.CString(output)
}

//for QSCKVStoreset
//func QSCKVStoreSet(k, v, privkey, chain string) string {
//	output := slim.QSCKVStoreSetPost(k, v, privkey, chain)
//	return output
//}

//export QSCKVStoreGet
func QSCKVStoreGet(k *C.char) *C.char {
	output := slim.QSCKVStoreGetQuery(C.GoString(k))
	return C.CString(output)
}

//for QSCKVStoreGet
//func QSCKVStoreGet(k string) string {
//	output := slim.QSCKVStoreGetQuery(k)
//	return output
//}

//export QSCQueryAccount
func QSCQueryAccount(addr *C.char) *C.char {
	output := slim.QSCQueryAccountGet(C.GoString(addr))
	return C.CString(output)
}

//for QSCQueryAccount
//func QSCQueryAccount(addr string) string {
//	output := slim.QSCQueryAccountGet(addr)
//	return output
//}

//export AccountRecover
func AccountRecover(mncode *C.char) *C.char {
	output := slim.AccountRecoverStr(C.GoString(mncode))
	return C.CString(output)
}

//for AccountRecovery
//func AccountRecover(mncode string) string {
//	output := slim.AccountRecoverStr(mncode)
//	return output
//}

//export GetIPfromInput
func GetIPfromInput(ip *C.char) {
	slim.GetIPfrom(C.GoString(ip))
}

//for IP input
//func GetIPfromInput(ip string) {
//	//	fmt.Println("Please input host including IP and port for initialization on Qstar deamon:")
//	slim.GetIPfrom(ip)
//}

//export PubAddrRetrieval
func PubAddrRetrieval(priv *C.char) *C.char {
	output := slim.PubAddrRetrievalStr(C.GoString(priv))
	return C.CString(output)
}

//for PubAddrRetrieval
//func PubAddrRetrieval(priv string) string {
//	//	fmt.Println("Please input host including IP and port for initialization on Qstar deamon:")
//	output := slim.PubAddrRetrievalStr(priv)
//	return output
//}

//export QSCtransferSend
func QSCtransferSend(addrto, coinstr, privkey, chainid *C.char) *C.char {
	output := slim.QSCtransferSendStr(C.GoString(addrto), C.GoString(coinstr), C.GoString(privkey), C.GoString(chainid))
	return C.CString(output)
}

//for QSCtransferSend
//func QSCtransferSend(addrto, coinstr, privkey, chainid string) string {
//	output := slim.QSCtransferSendStr(addrto, coinstr, privkey, chainid)
//	return output
//}

func main() {
	//GetIPfromInput("192.168.1.23:1317")
	////output := AccountCreate()
	////fmt.Println(output)
	////
	////out := QSCQueryAccount("address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355")
	////fmt.Println(out)
	//out := QSCKVStoreSet("14", "Merfer", "0xa328891040ae9b773bcd30005235f99a8d62df03a89e4f690f9fa03abb1bf22715fc9ca05613f2d8061492e9f8149510b5b67d340d199ff24f34c85dbbbd7e0df780e9a6cc", "test-chain-HfiBIx")
	//fmt.Println(out)
	//kvout := QSCKVStoreGet("13")
	//fmt.Println(kvout)
	//AccountRecover("vague success fresh check remove banner music snap jelly medal bring mix eagle seat cash off winter mean comic turn always teach tiny wagon")
	//puba := PubAddrRetrieval("oyiJEECum3c7zTAAUjX5mo1i3wOonk9pD5+gOrsb8icV/JygVhPy2AYUkun4FJUQtbZ9NA0Zn/JPNMhdu71+DfeA6abM")
	//fmt.Println(puba)
	//transoutb64 := QSCtransferSend("address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355", "2qos", "oyiJEECum3c7zTAAUjX5mo1i3wOonk9pD5+gOrsb8icV/JygVhPy2AYUkun4FJUQtbZ9NA0Zn/JPNMhdu71+DfeA6abM", "qos-testapp")
	//fmt.Println(transoutb64)
}
