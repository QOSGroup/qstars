package main

import "C"
import (
	"fmt"
	"github.com/QOSGroup/qstars/slim"
)

//export QSCKVStoreSet
func QSCKVStoreSet(k, v, privkey, chain *C.char) *C.char {
	output := slim.QSCKVStoreSetPost(C.GoString(k), C.GoString(v), C.GoString(privkey), C.GoString(chain))
	return C.CString(output)
}

//export QSCKVStoreGet
func QSCKVStoreGet(k *C.char) *C.char {
	output := slim.QSCKVStoreGetQuery(C.GoString(k))
	return C.CString(output)
}

func AccountCreate() string {
	output := slim.AccountCreateStr()
	return output
}

func main() {
	output := AccountCreate()
	fmt.Println(output)
}
