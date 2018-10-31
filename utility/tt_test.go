// Copyright 2018 The QOS Authors

package utility

import (
	"encoding/hex"
	"testing"
)

func TestEncodeToString(t *testing.T) {
	s := "rDwWppdGKFCv0wUxFqVID87GI/CFwLbL9p6EM6ug5brPbkXQoZMIH9+Rgi1/vFcNJUHp88fKZDNFdEif8dg73A=="
	t.Log(hex.EncodeToString(Decbase64(s)))
}
