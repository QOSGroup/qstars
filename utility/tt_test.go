// Copyright 2018 The QOS Authors

package utility

import (
	"encoding/hex"
	"testing"
)

func TestEncodeToString(t *testing.T) {
	s := "a328891040ae9b773bcd30005235f99a8d62df03a89e4f690f9fa03abb1bf22715fc9ca05613f2d8061492e9f8149510b5b67d340d199ff24f34c85dbbbd7e0df780e9a6cc"
	t.Log(hex.EncodeToString(Decbase64(s)))
}

func TestEncbase64(t *testing.T) {
	s := "0xa328891040ae9b773bcd30005235f99a8d62df03a89e4f690f9fa03abb1bf22715fc9ca05613f2d8061492e9f8149510b5b67d340d199ff24f34c85dbbbd7e0df780e9a6cc"
	hexbyte, _ := hex.DecodeString(s[2:])
	privkeybase64 := Encbase64(hexbyte)
	t.Log(privkeybase64)
}

func TestDecbase64(t *testing.T) {
	s := "rpt3O80wAFI1+ZqNYt8DqJ5PaQ+foDq7G/InFfycoFYT8tgGFJLp+BSVELW2fTQNGZ/yTzTIXbu9fg33gOmmzA=="
	bz := Decbase64(s)
	t.Log(bz)
}

//func TestPubAddrRetrievalFromAmino(t *testing.T) {
//	s := "oyiJEECum3c7zTAAUjX5mo1i3wOonk9pD5+gOrsb8icV/JygVhPy2AYUkun4FJUQtbZ9NA0Zn/JPNMhdu71+DfeA6abM"
//	Pub, Addr, _ := PubAddrRetrievalFromAmino(s, cdc)
//	t.Log(Pub, Addr)
//
//}
