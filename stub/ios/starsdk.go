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
func QSCKVStoreGet(ul string) string {
	output := stub.QSCKVStoreGetQuery(ul)
	return output
}

//for QSCQueryAccount
func QSCQueryAccount(ul string) string {
	output := stub.QSCQueryAccountGet(ul)
	return output
}


//for QSCtransfer
func QSCtransfer(ul, a, privkey, chain, ac, seq, g string) string {
	output := stub.QSCtransferPost(ul,a,privkey,chain,ac,seq,g)
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

