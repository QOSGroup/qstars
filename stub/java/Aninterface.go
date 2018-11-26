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
////for DispatchCoins
//func DispatchCoins(addrs, coins, causecodes, causestrings, gas string) string {
//	output := jsdk.DispatchCoins(addrs, coins, causecodes, causestrings, gas)
//	return output
//}

//export NewArticle
//func NewArticle(authorAddress, originAuthor, articleHash, shareAuthor, shareOriginAuthor, shareCommunity, shareInvestor, endInvestDate, endBuyDate *C.char) *C.char {
//	output := jsdk.NewArticle(C.GoString(authorAddress), C.GoString(originAuthor), C.GoString(articleHash), C.GoString(shareAuthor), C.GoString(shareOriginAuthor), C.GoString(shareCommunity), C.GoString(shareInvestor), C.GoString(endInvestDate), C.GoString(endBuyDate))
//	return C.CString(output)
//}

//for NewArticle
//func NewArticle(authorAddress, originAuthor, articleHash, shareAuthor, shareOriginAuthor, shareCommunity, shareInvestor, endInvestDate, endBuyDate string) string {
//	output := jsdk.NewArticle(authorAddress, originAuthor, articleHash, shareAuthor, shareOriginAuthor, shareCommunity, shareInvestor, endInvestDate, endBuyDate)
//	return output
//}

func main() {
	//InitJNI()
	////send --from=rpt3O80wAFI1+ZqNYt8DqJ5PaQ+foDq7G/InFfycoFYT8tgGFJLp+BSVELW2fTQNGZ/yTzTIXbu9fg33gOmmzA== --to=address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355 --amount=2qos
	//r := SendByJNI("rpt3O80wAFI1+ZqNYt8DqJ5PaQ+foDq7G/InFfycoFYT8tgGFJLp+BSVELW2fTQNGZ/yTzTIXbu9fg33gOmmzA==", "address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355", "2qos")
	//fmt.Println(r)

	//q := jsdk.QueryAccbyJNI("address1k0m8ucnqug974maa6g36zw7g2wvfd4sug6uxay")
	//fmt.Println(q)
	//disout := DispatchCoins("a12adc23|18671eab2", "1aoe|10aoe|30aoe", "1|2|1", "signin,invited,abc", "0QOS")
	//fmt.Println(disout)
	//newout := NewArticle("a12f87bc", "a12f87bc", "a12f87bc2", "20", "20", "20", "20", "20days", "20days")
	//fmt.Println(newout)
}
