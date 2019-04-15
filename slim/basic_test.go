package slim

import "testing"

func TestAccountCreate(t *testing.T) {
	password := "qstars"
	output := AccountCreateStr(password)
	t.Log(output)
}

func TestAccountRecoverStr(t *testing.T) {
	mncode := "oyster leave weird tiger road rose anger garden planet price small rain cradle rhythm wine spider manual wave plastic solar spray battle parent match"
	password := "qstars"
	output := AccountRecoverStr(mncode, password)
	t.Log(output)
}

func TestPubAddrRetrievalStr(t *testing.T) {
	s := "Me7Ts4jHgytbyq5ctu/lzgT630jKVxGt/Yb14HvHU4zfRHBf1KkdI7uo8H/vWQRFqJVLZeo5MEjUjDL08BWcPw=="
	output := PubAddrRetrievalStr(s)
	t.Log(output)
}

func TestAccountCreateFromSeed(t *testing.T) {
	mncode := "oyster leave weird tiger road rose anger garden planet price small rain cradle rhythm wine spider manual wave plastic solar spray battle parent match"
	output := AccountCreateFromSeed(mncode)
	t.Log(output)
}
