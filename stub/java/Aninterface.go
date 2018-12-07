package main

import "C"
import (
	"github.com/QOSGroup/qstars/stub/java/jsdk"
)

// ----------------------------------------------------------------------------
// source code for so file generation with "go build " command, e.g.go build -o awesome.so -buildmode=c-shared awesome.go
// ----------------------------------------------------------------------------

//export InitJNI
func InitJNI() {
	jsdk.InitJNI()
}

//export SendByJNI
//func SendByJNI(fromStr, toStr1, coinstr *C.char) *C.char {
//	output := jsdk.SendByJNI(C.GoString(fromStr), C.GoString(toStr1), C.GoString(coinstr))
//	return C.CString(output)
//}

//for SendByJNI
func SendByJNI(fromStr, toStr1, coinstr string) string {
	output := jsdk.SendByJNI(fromStr, toStr1, coinstr)
	return output
}

//export DispatchCoins
//func DispatchCoins(addrs, coins, causestrings, causecodes, gas *C.char) *C.char {
//	output := jsdk.DispatchCoins(C.GoString(addrs), C.GoString(coins), C.GoString(causestrings), C.GoString(causecodes), C.GoString(gas))
//	return C.CString(output)
//}
//for DispatchCoins
func DispatchCoins(addrs, coins, causecodes, causestrings, gas string) string {
	output := jsdk.DispatchCoins(addrs, coins, causecodes, causestrings, gas)
	return output
}

//export NewArticle
//func NewArticle(authorAddress, originAuthor, articleHash, shareAuthor, shareOriginAuthor, shareCommunity, shareInvestor, endInvestDate, endBuyDate *C.char) *C.char {
//	output := jsdk.NewArticle(C.GoString(authorAddress), C.GoString(originAuthor), C.GoString(articleHash), C.GoString(shareAuthor), C.GoString(shareOriginAuthor), C.GoString(shareCommunity), C.GoString(shareInvestor), C.GoString(endInvestDate), C.GoString(endBuyDate))
//	return C.CString(output)
//}

//for NewArticle
func NewArticle(authorAddress, originAuthor, articleHash, shareAuthor, shareOriginAuthor, shareCommunity, shareInvestor, endInvestDate, endBuyDate string) string {
	output := jsdk.NewArticle(authorAddress, originAuthor, articleHash, shareAuthor, shareOriginAuthor, shareCommunity, shareInvestor, endInvestDate, endBuyDate)
	return output
}

//export InvestAd
//func InvestAd(chainId, articleHash, coins, privatekey *C.char, nonce int64) *C.char {
//	output := jsdk.InvestAd(C.GoString(chainId), C.GoString(articleHash), C.GoString(coins), C.GoString(privatekey),nonce)
//	return C.CString(output)
//}

//for InvestAd
func InvestAdBackground(txb string) string {
	output := jsdk.InvestAdBackground(txb)
	return output
}

//export BuyAd
//func BuyAd(chainId, articleHash, coins, privatekey *C.char, nonce int64) *C.char {
//	output := jsdk.InvestAd(C.GoString(chainId), C.GoString(articleHash), C.GoString(coins), C.GoString(privatekey),nonce)
//	return C.CString(output)
//}

//for BuyAd
func BuyAd(chainId, articleHash, coins, privatekey string, nonce int64) string {
	output := jsdk.BuyAd(chainId, articleHash, coins, privatekey, nonce)
	return output
}

//for investAdbckaground testing
//type ResultInvest struct {
//	Code   string          `json:"code"`
//	Reason string          `json:"reason,omitempty"`
//	Result json.RawMessage `json:"result,omitempty"`
//}

func main() {
	//InitJNI()
	////send --from=rpt3O80wAFI1+ZqNYt8DqJ5PaQ+foDq7G/InFfycoFYT8tgGFJLp+BSVELW2fTQNGZ/yTzTIXbu9fg33gOmmzA== --to=address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355 --amount=2qos
	////r := SendByJNI("Ey+2bNFF2gTUV6skSBgRy3rZwo9nS4Dw0l2WpLrhVvV8MuMRbjN4tUK8orHiJgHTR+enkxyXcA8giVrsrIRM4Q==", "address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355", "2qos")
	////fmt.Println(r)
	//
	////disout := DispatchCoins("a12adc23|18671eab2", "1aoe|10aoe|30aoe", "1|2|1", "signin,invited,abc", "0QOS")
	////fmt.Println(disout)
	//newout := NewArticle("address13mjc3n3xxj73dhkju9a0dfr4lrfvv3whxqg0dy", "address1zsqzn6wdecyar6c6nzem3e8qss2ws95csr8d0r", "a123", "20", "20", "20", "20", "20", "3")
	//fmt.Println(newout)

	//slim.SetBlockchainEntrance("192.168.1.23:1317", "forQmoonAddr")
	//ad := slim.JQInvestAd("qos-testapp", "qstars-test", "abcd", "1AOE", "Ey+2bNFF2gTUV6skSBgRy3rZwo9nS4Dw0l2WpLrhVvV8MuMRbjN4tUK8orHiJgHTR+enkxyXcA8giVrsrIRM4Q==")
	//var ri ResultInvest
	//err := json.Unmarshal([]byte(ad), &ri)
	//fmt.Printf("error is:%s\n ", err)
	//Adout := InvestAdBackground(string(ri.Result))
	//fmt.Println(Adout)
}
