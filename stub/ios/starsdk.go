package starsdk

import "C"
import (
	"github.com/QOSGroup/qstars/stub"
)

//package starsdk is for ios app generation
//for Account Create
func AccountCreate() string {
	output := stub.AccountCreateStr()
	return output
}

//for QSCKVStoreset
func QSCKVStoreSet(k, v, privkey, chain string) int {
	output := stub.QSCKVStoreSetPost(k, v, privkey, chain)
	return output
}

//for QSCKVStoreGet
func QSCKVStoreGet(k string) string {
	output := stub.QSCKVStoreGetQuery(k)
	return output
}

//for QSCQueryAccount
func QSCQueryAccount(addr string) string {
	output := stub.QSCQueryAccountGet(addr)
	return output
}

//for AccountRecovery
func AccountRecover(mncode string) string {
	output := stub.AccountRecoverStr(mncode)
	return output
}

//for IP input
func GetIPfromInput(ip string) {
	//	fmt.Println("Please input host including IP and port for initialization on Qstar deamon:")
	stub.GetIPfrom(ip)
}

//for PubAddrRetrieval
func PubAddrRetrieval(priv string) string {
	//	fmt.Println("Please input host including IP and port for initialization on Qstar deamon:")
	output := stub.PubAddrRetrievalStr(priv)
	return output
}

//for QSCtransferSend
func QSCtransferSend(addrto, coinstr, privkey string) string {
	output := stub.QSCtransferSendStr(addrto, coinstr, privkey)
	return output
}
