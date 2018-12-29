package starsdkapp

import (
	"github.com/QOSGroup/qstars/slim"
)

func AccountCreate(password string) string {
	output := slim.AccountCreateStr(password)
	return output
}

//for QSCKVStoreset
func QSCKVStoreSet(k, v, privkey, chain string) string {
	output := slim.QSCKVStoreSetPost(k, v, privkey, chain)
	return output
}

//for QSCKVStoreGet
func QSCKVStoreGet(k string) string {
	output := slim.QSCKVStoreGetQuery(k)
	return output
}

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

//for AccountRecovery
func AccountRecover(mncode, password string) string {
	output := slim.AccountRecoverStr(mncode, password)
	return output
}

//for IP input
func SetBlockchainEntrance(sh, mh string) {
	slim.SetBlockchainEntrance(sh, mh)
}

//for PubAddrRetrieval
func PubAddrRetrieval(priv string) string {
	//	fmt.Println("Please input host including IP and port for initialization on Qstar deamon:")
	output := slim.PubAddrRetrievalStr(priv)
	return output
}

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

func JQInvestAd(QOSchainId, QSCchainId, articleHash, coins, privatekey string) string {
	output := slim.JQInvestAd(QOSchainId, QSCchainId, articleHash, coins, privatekey)
	return output
}

func AesEncrypt(key, plainText string) string {
	output := slim.AesEncrypt(key, plainText)
	return output
}

func AesDecrypt(key, cipherText string) string {
	output := slim.AesDecrypt(key, cipherText)
	return output
}

func TransferRecordsQuery(chainid, addr, cointype, offset, limit string) string {
	output := slim.TransferRecordsQuery(chainid, addr, cointype, offset, limit)
	return output
}

//func QSCtransfer() {
//	output := slim.QSCtransferSendStr("address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355", "22qos", "0xa328891040ae9b773bcd30005235f99a8d62df03a89e4f690f9fa03abb1bf22715fc9ca05613f2d8061492e9f8149510b5b67d340d199ff24f34c85dbbbd7e0df780e9a6cc", "qos-testapp")
//	fmt.Println(output)
//}
