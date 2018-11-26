package slim

import "testing"

func TestAccountCreate(t *testing.T) {
	output := AccountCreateStr()
	t.Log(output)
}

func TestAccountRecoverStr(t *testing.T) {
	mncode := "ethics erase secret rail frost talent lady load involve dream cushion lava void ten penalty off better include path produce they amateur avoid theory"
	output := AccountRecoverStr(mncode)
	t.Log(output)
}

func TestPubAddrRetrievalStr(t *testing.T) {
	s := "rpt3O80wAFI1+ZqNYt8DqJ5PaQ+foDq7G/InFfycoFYT8tgGFJLp+BSVELW2fTQNGZ/yTzTIXbu9fg33gOmmzA=="
	output := PubAddrRetrievalStr(s)
	t.Log(output)
}
