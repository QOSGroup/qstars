package main
import "C"
import (
	"github.com/QOSGroup/qstars/stub/java/quzhuan/jsdk"
)





//export CommonHandler
func CommonHandler(url,chainid,funcName, privatekey, args *C.char) *C.char {
	output := jsdk.CommonHandler(C.GoString(url),C.GoString(chainid),C.GoString(funcName), C.GoString(privatekey), C.GoString(args))
	return C.CString(output)
}

//export ParameterFormat
func ParameterFormat(parameter *C.char) *C.char {
	output := jsdk.ParameterFormat(C.GoString(parameter))
	return C.CString(output)
}


//export CreateAccount
func CreateAccount() *C.char {
	output := jsdk.CreateAccount()
	return C.CString(output)
}



func main(){

}