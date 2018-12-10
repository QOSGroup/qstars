package main

import "C"
import (
	"github.com/QOSGroup/qstars/slim"
)

//export AccountCreate
//func AccountCreate() *C.char {
//	output := slim.AccountCreateStr()
//	return C.CString(output)
//}

func AccountCreate(password string) string {
	output := slim.AccountCreateStr(password)
	return output
}

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
//func QSCQueryAccount(addr *C.char) *C.char {
//	output := slim.QSCQueryAccountGet(C.GoString(addr))
//	return C.CString(output)
//}

//for QSCQueryAccount
func QSCQueryAccount(addr string) string {
	output := slim.QSCQueryAccountGet(addr)
	return output
}

//for QOSQueryAccount
func QOSQueryAccount(addr string) string {
	output := slim.QOSQueryAccountGet(addr)
	return output
}

//export AccountRecover
//func AccountRecover(mncode, password *C.char) *C.char {
//	output := slim.AccountRecoverStr(C.GoString(mncode), C.GoString(password))
//	return C.CString(output)
//}

//for AccountRecovery
func AccountRecover(mncode, password string) string {
	output := slim.AccountRecoverStr(mncode, password)
	return output
}

//export GetIPfromInput
//func SetBlockchainEntrance(sh, mh *C.char) {
//	slim.SetBlockchainEntrance(C.GoString(sh), C.GoString(mh))
//}

//for hosts input
func SetBlockchainEntrance(sh, mh string) {
	slim.SetBlockchainEntrance(sh, mh)
}

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
//func QSCtransferSend(addrto, coinstr, privkey, chainid *C.char) *C.char {
//	output := slim.QSCtransferSendStr(C.GoString(addrto), C.GoString(coinstr), C.GoString(privkey), C.GoString(chainid))
//	return C.CString(output)
//}

//for QSCtransferSend
func QSCtransferSend(addrto, coinstr, privkey, chainid string) string {
	output := slim.QSCtransferSendStr(addrto, coinstr, privkey, chainid)
	return output
}

//for QOSCommitResultCheck
func QOSCommitResultCheck(txhash, height string) string {
	output := slim.QOSCommitResultCheck(txhash, height)
	return output
}

//for InvesAd in mobile app
func JQInvestAd(QOSchainId, QSCchainId, articleHash, coins, privatekey string) string {
	output := slim.JQInvestAd(QOSchainId, QSCchainId, articleHash, coins, privatekey)
	return output
}

func main() {
	//SetBlockchainEntrance("192.168.1.23:1317", "forQmoonAddr")
	//output := AccountCreate("qstars")
	//fmt.Println(output)
	//out := QSCQueryAccount("address13l90zvt26szkrquutwpgj7kef58mgyntfs46l2")
	//fmt.Println(out)
	//out := QSCKVStoreSet("14", "Merfer", "0xa328891040ae9b773bcd30005235f99a8d62df03a89e4f690f9fa03abb1bf22715fc9ca05613f2d8061492e9f8149510b5b67d340d199ff24f34c85dbbbd7e0df780e9a6cc", "test-chain-HfiBIx")
	//fmt.Println(out)
	//kvout := QSCKVStoreGet("13")
	//fmt.Println(kvout)
	//aout := AccountRecover("address rescue flower seven erode trigger panic apart mango put tenant version matrix devote ozone critic damp edge panda tuition view index sound account", "qstars")
	//fmt.Println(aout)
	//puba := PubAddrRetrieval("oyiJEECum3c7zTAAUjX5mo1i3wOonk9pD5+gOrsb8icV/JygVhPy2AYUkun4FJUQtbZ9NA0Zn/JPNMhdu71+DfeA6abM")
	//fmt.Println(puba)
	//transoutb64 := QSCtransferSend("address13l90zvt26szkrquutwpgj7kef58mgyntfs46l2", "2aoe", "Ey+2bNFF2gTUV6skSBgRy3rZwo9nS4Dw0l2WpLrhVvV8MuMRbjN4tUK8orHiJgHTR+enkxyXcA8giVrsrIRM4Q==", "qos-testapp")
	//fmt.Println(transoutb64)
	//queryResult := QOSCommitResultCheck("1915BF14E0583E0F38D695F12EF122D017AAAA86", "341")
	//fmt.Println(queryResult)

	//jqresult := JQInvestAd("qos-testapp", "qstars-test", "abcd", "1AOE", "Ey+2bNFF2gTUV6skSBgRy3rZwo9nS4Dw0l2WpLrhVvV8MuMRbjN4tUK8orHiJgHTR+enkxyXcA8giVrsrIRM4Q==")
	//fmt.Println(jqresult)
}
