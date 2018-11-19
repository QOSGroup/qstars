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
	s := "oyiJEECFjzYxyq47liWtzFyxEldTauWJaUBAfcmchfIDwEIcVlnlgXcC9Ev+S7Jc7L0l4AI9E1E1hTHH3fp+V0yGUoLp"
	output := PubAddrRetrievalStr(s)
	t.Log(output)
}
