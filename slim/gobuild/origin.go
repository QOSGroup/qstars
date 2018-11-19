package main

import (
	"fmt"
	"github.com/QOSGroup/qstars/slim"
)

func AccountCreate() string {
	output := slim.AccountCreateStr()
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

//for AccountRecovery
func AccountRecover(mncode string) string {
	output := slim.AccountRecoverStr(mncode)
	return output
}

//for IP input
func GetIPfromInput(ip string) {
	//	fmt.Println("Please input host including IP and port for initialization on Qstar deamon:")
	slim.GetIPfrom(ip)
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

func main() {
	//GetIPfromInput("localhost:1317")
	output := AccountCreate()
	fmt.Println(output)
	//acc, err := testQuery("http://localhost:1317", "address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355")
	//fmt.Printf("---acc:%+v, err:%+v\n", acc, err)
	//out := QSCQueryAccount("address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355")
	//fmt.Println(out)
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
