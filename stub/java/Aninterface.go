package main 
import "C"
import (
	"fmt"
		"github.com/QOSGroup/qstars/stub"
)
// ----------------------------------------------------------------------------
// source code for so file generation with command, e.g. go build go build -o awesome.so -buildmode=c-shared awesome.go
// ----------------------------------------------------------------------------
// init func for the query the IP of the seeds node in qsatrs
//export init
//func init() {
//
//}


//export AccountCreate
func AccountCreate() *C.char {
	output := stub.AccountCreateStr()
	return C.CString(output)
}

//export QSCKVStoreSet
func QSCKVStoreSet(k,v,privkey,chain *C.char) int {
	output := stub.QSCKVStoreSetPost(C.GoString(k),C.GoString(v),C.GoString(privkey),C.GoString(chain))
	return output
}

//export QSCKVStoreGet
func QSCKVStoreGet(ul *C.char) *C.char {
	output := stub.QSCKVStoreGetQuery(C.GoString(ul))
	return C.CString(output)
}

//export QSCQueryAccount
func QSCQueryAccount(ul *C.char) *C.char {
	output := stub.QSCQueryAccountGet(C.GoString(ul))
	return C.CString(output)
}

//export QSCMintCoin
func QSCMintCoin() {
	fmt.Println("this is QSCMintCoin function interface")
}

//export QSCtransfer
func QSCtransfer(ul,a,privkey,chain,ac,seq,g *C.char) *C.char {
	output := stub.QSCtransferPost(C.GoString(ul),C.GoString(a),C.GoString(privkey),C.GoString(chain),C.GoString(ac),C.GoString(seq),C.GoString(g))
	return C.CString(output)
}

//export QOStoQSCtransfer
func QOStoQSCtransfer() {
	fmt.Println("this is QOStoQSCtransfer function interface")
}

//export QSCtoQOStransfer
func QSCtoQOStransfer() {
	fmt.Println("this is QSCtoQOStransfer function interface")
}

func main() {
	//AccountCreate()
//	out := QSCQueryAccount("http://localhost:1317/accounts/cosmosaccaddr1nskqcg35k8du3ydhntkcqjxtk254qv8me943mv")
//	fmt.Println(out)
}

