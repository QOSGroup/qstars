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
func SendByJNI(fromStr, toStr1, coinstr *C.char) *C.char {
	output := jsdk.SendByJNI(C.GoString(fromStr), C.GoString(toStr1), C.GoString(coinstr))
	return C.CString(output)
}

//for SendByJNI
//func SendByJNI(fromStr, toStr1, coinstr string) string {
//	output := jsdk.SendByJNI(fromStr, toStr1, coinstr)
//	return output
//}

//export DispatchCoins
func DispatchCoins(addrs, coins, causestrings, causecodes, gas *C.char) *C.char {
	output := jsdk.DispatchCoins(C.GoString(addrs), C.GoString(coins), C.GoString(causestrings), C.GoString(causecodes), C.GoString(gas))
	return C.CString(output)
}

//for DispatchCoins
//func DispatchCoins(addrs, coins, causecodes, causestrings, gas string) string {
//	output := jsdk.DispatchCoins(addrs, coins, causecodes, causestrings, gas)
//	return output
//}

//export NewArticle
func NewArticle(authorAddress, originAuthor, articleHash, shareAuthor, shareOriginAuthor, shareCommunity, shareInvestor, endInvestDate, endBuyDate *C.char) *C.char {
	output := jsdk.NewArticle(C.GoString(authorAddress), C.GoString(originAuthor), C.GoString(articleHash), C.GoString(shareAuthor), C.GoString(shareOriginAuthor), C.GoString(shareCommunity), C.GoString(shareInvestor), C.GoString(endInvestDate), C.GoString(endBuyDate))
	return C.CString(output)
}

//for NewArticle
//func NewArticle(authorAddress, originAuthor, articleHash, shareAuthor, shareOriginAuthor, shareCommunity, shareInvestor, endInvestDate, endBuyDate string) string {
//	output := jsdk.NewArticle(authorAddress, originAuthor, articleHash, shareAuthor, shareOriginAuthor, shareCommunity, shareInvestor, endInvestDate, endBuyDate)
//	return output
//}

//export InvestAdBackground
func InvestAdBackground(txb *C.char) *C.char {
	output := jsdk.InvestAdBackground(C.GoString(txb))
	return C.CString(output)
}

//for InvestAdBackground
//func InvestAdBackground(txb string) string {
//	output := jsdk.InvestAdBackground(txb)
//	return output
//}

//export BuyAd
func BuyAd(articleHash, coins, buyer *C.char) *C.char {
	output := jsdk.BuyAd(C.GoString(articleHash), C.GoString(coins), C.GoString(buyer))
	return C.CString(output)
}

//for BuyAd
//func BuyAd(articleHash, coins, buyer string) string {
//	output := jsdk.BuyAd(articleHash, coins, buyer)
//	return output
//}

//for investAdbckaground testing
//type ResultInvest struct {
//	Code   string          `json:"code"`
//	Reason string          `json:"reason,omitempty"`
//	Result json.RawMessage `json:"result,omitempty"`
//}

//export RetrieveInvestors
func RetrieveInvestors(articleHash *C.char) *C.char {
	output := jsdk.RetrieveInvestors(C.GoString(articleHash))
	return C.CString(output)
}

//for RetrieveInvestors
//func RetrieveInvestors(articleHash string) string {
//	output := jsdk.RetrieveInvestors(articleHash)
//	return output
//}

//export RetrieveBuyer
func RetrieveBuyer(articleHash *C.char) *C.char {
	output := jsdk.RetrieveBuyer(C.GoString(articleHash))
	return C.CString(output)
}

//for RetrieveBuyer
//func RetrieveBuyer(articleHash string) string {
//	output := jsdk.RetrieveBuyer(articleHash)
//	return output
//}

//export QueryArticle
func QueryArticle(articleHash *C.char) *C.char {
	output := jsdk.QueryArticle(C.GoString(articleHash))
	return C.CString(output)
}

//for QueryArticle
//func QueryArticle(articleHash string) string {
//	output := jsdk.QueryArticle(articleHash)
//	return output
//}

//export QueryCoins
func QueryCoins(txHash *C.char) *C.char {
	output := jsdk.QueryCoins(C.GoString(txHash))
	return C.CString(output)
}

//for QueryCoins
//func QueryCoins(articleHash string) string {
//	output := jsdk.QueryCoins(articleHash)
//	return output
//}

func main() {
	//InitJNI()
	//send --from=rpt3O80wAFI1+ZqNYt8DqJ5PaQ+foDq7G/InFfycoFYT8tgGFJLp+BSVELW2fTQNGZ/yTzTIXbu9fg33gOmmzA== --to=address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355 --amount=2qos
	//r := SendByJNI("Ey+2bNFF2gTUV6skSBgRy3rZwo9nS4Dw0l2WpLrhVvV8MuMRbjN4tUK8orHiJgHTR+enkxyXcA8giVrsrIRM4Q==", "address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355", "2qos")
	//fmt.Println(r)
	//
	//disout := DispatchCoins("address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355|address1zsqzn6wdecyar6c6nzem3e8qss2ws95csr8d0r", "500|400", "2|3", "qiandao|shiming", "0QOS")
	//fmt.Println(disout)
	//newout := NewArticle("address13mjc3n3xxj73dhkju9a0dfr4lrfvv3whxqg0dy", "address1zsqzn6wdecyar6c6nzem3e8qss2ws95csr8d0r", "a123", "20", "20", "20", "20", "20", "3")
	//fmt.Println(newout)

	//slim.SetBlockchainEntrance("192.168.1.23:1317", "forQmoonAddr")
	//ad := slim.JQInvestAd("qos-testapp", "qstars-test", "abcd", "1AOE", "Ey+2bNFF2gTUV6skSBgRy3rZwo9nS4Dw0l2WpLrhVvV8MuMRbjN4tUK8orHiJgHTR+enkxyXcA8giVrsrIRM4Q==")
	//var ri ResultInvest
	//err := json.Unmarshal([]byte(ad), &ri)
	//fmt.Printf("error is:%s\n ", err)
	//Adout := InvestAdBackground(string(ri.Result))
	//fmt.Println(Adout)

	//Buad := BuyAd("abcd", "10QOS", "Ey+2bNFF2gTUV6skSBgRy3rZwo9nS4Dw0l2WpLrhVvV8MuMRbjN4tUK8orHiJgHTR+enkxyXcA8giVrsrIRM4Q==")
	//fmt.Println(Buad)

	//reinv := RetrieveInvestors("abcd")
	//fmt.Println(reinv)

	//reby := RetrieveBuyer("abcd")
	//fmt.Println(reby)

	//qa := QueryArticle("adcd")
	//fmt.Println(qa)

	//qc := QueryCoins("adcd")
	//fmt.Println(qc)
}
