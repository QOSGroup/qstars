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
