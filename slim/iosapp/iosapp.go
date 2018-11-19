package iosapp

import (
	"github.com/QOSGroup/qstars/slim"
)

//for Account Create
func AccountCreate() string {
	output := slim.AccountCreateStr()
	return output
}

//for AccountRecovery
func AccountRecover(mncode string) string {
	output := slim.AccountRecoverStr(mncode)
	return output
}

//for PubAddrRetrieval
func PubAddrRetrieval(priv string) string {
	//	fmt.Println("Please input host including IP and port for initialization on Qstar deamon:")
	output := slim.PubAddrRetrievalStr(priv)
	return output
}

//func QSCtransfer() {
//	output := slim.QSCtransferSendStr("address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355", "22qos", "0xa328891040ae9b773bcd30005235f99a8d62df03a89e4f690f9fa03abb1bf22715fc9ca05613f2d8061492e9f8149510b5b67d340d199ff24f34c85dbbbd7e0df780e9a6cc", "qos-testapp")
//	fmt.Println(output)
//}
