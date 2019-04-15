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

//export DispatchCoins
func DispatchCoins(addrs, coins, causestrings, causecodes, gas *C.char) *C.char {
	output := jsdk.DispatchCoins(C.GoString(addrs), C.GoString(coins), C.GoString(causestrings), C.GoString(causecodes), C.GoString(gas))
	return C.CString(output)
}

//export NewArticle
func NewArticle(authorAddress, articleType, articleHash, shareAuthor, shareOriginAuthor, shareCommunity, shareInvestor, endInvestDate, endBuyDate, cointype *C.char) *C.char {
	output := jsdk.NewArticle(C.GoString(authorAddress), C.GoString(articleType), C.GoString(articleHash), C.GoString(shareAuthor), C.GoString(shareOriginAuthor), C.GoString(shareCommunity), C.GoString(shareInvestor), C.GoString(endInvestDate), C.GoString(endBuyDate), C.GoString(cointype))
	return C.CString(output)
}

//export AcutionAdBackground
func AcutionAdBackground(txb *C.char) *C.char {
	output := jsdk.AcutionAdBackground(C.GoString(txb))
	return C.CString(output)
}

//export AcutionAd
func AcutionAd(articleHash, address, coinsType, coinAmount, qscnonce *C.char) string {
	output := jsdk.AcutionAd(C.GoString(articleHash), C.GoString(address), C.GoString(coinsType), C.GoString(coinAmount), C.GoString(qscnonce))
	return output
}

//export QueryMaxAcution
func QueryMaxAcution(txb *C.char) *C.char {
	output := jsdk.QueryMaxAcution(C.GoString(txb))
	return C.CString(output)
}

//export QueryAllAcution
func QueryAllAcution(txb *C.char) *C.char {
	output := jsdk.QueryAllAcution(C.GoString(txb))
	return C.CString(output)
}

//export InvestAdBackground
func InvestAdBackground(txb *C.char) *C.char {
	output := jsdk.InvestAdBackground(C.GoString(txb))
	return C.CString(output)
}

//export Distribution
func Distribution(articleHash *C.char) *C.char {
	output := jsdk.Distribution(C.GoString(articleHash))
	return C.CString(output)
}

//export RetrieveInvestors
func RetrieveInvestors(articleHash *C.char) *C.char {
	output := jsdk.RetrieveInvestors(C.GoString(articleHash))
	return C.CString(output)
}

//export QueryArticle
func QueryArticle(articleHash *C.char) *C.char {
	output := jsdk.QueryArticle(C.GoString(articleHash))
	return C.CString(output)
}

//export QueryBlance
func QueryBlance(txHash *C.char) *C.char {
	output := jsdk.QueryBlance(C.GoString(txHash))
	return C.CString(output)
}

//export AdvertisersTrue
func AdvertisersTrue(privatekey, coinsType, coinAmount, qscnonce string) string {
	output := jsdk.AdvertisersTrue(C.GoString(privatekey), C.GoString(coinsType), C.GoString(coinAmount), C.GoString(qscnonce))
	return output
}

//export AdvertisersFalse
func AdvertisersFalse(privatekey, coinsType, coinAmount, qscnonce string) string {
	output := jsdk.AdvertisersFalse(C.GoString(privatekey), C.GoString(coinsType), C.GoString(coinAmount), C.GoString(qscnonce))
	return output
}

//export Recharge
func Recharge(privatekey, address, coinsType, coinAmount, qscnonce string) string {
	output := jsdk.Recharge(C.GoString(privatekey), C.GoString(address), C.GoString(coinsType), C.GoString(coinAmount), C.GoString(qscnonce))
	return C.CString(output)
}

//export Extract
func Extract(privatekey, address, coinsType, coinAmount, qscnonce *C.char) string {
	output := jsdk.Extract(C.GoString(privatekey), C.GoString(address), C.GoString(coinsType), C.GoString(coinAmount), C.GoString(qscnonce))
	return C.CString(output)
}

//export QSCCommitResultCheck
func QSCCommitResultCheck(txhash, height *C.char) *C.char {
	output := jsdk.QSCCommitResultCheck(C.GoString(txhash), C.GoString(height))
	return C.CString(output)
}

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

	//slim.SetBlockchainEntrance("192.168.1.223:1317", "192.168.1.223:9527")
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

	//trc := slim.TransferRecordsQuery("address1l7d3dc26adk9gwzp777s3a9p5tprn7m43p99cg", "AOE")
	//fmt.Println(trc)
}
