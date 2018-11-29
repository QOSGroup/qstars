package slim

import "testing"

func TestAccountCreate(t *testing.T) {
	output := AccountCreateStr()
	t.Log(output)
}

func TestAccountRecoverStr(t *testing.T) {
	mncode := "ginger heavy absorb annual act two open negative science elite possible blur quantum deer start shove width vacant power tomato nut absurd family rocket"
	output := AccountRecoverStr(mncode)
	t.Log(output)
}

func TestPubAddrRetrievalStr(t *testing.T) {
	s := "4UWd8jmHYjSye6rnhbGHqmZj9Mr8vb3BaGbstyERpjPL683/nAd1x/piZI6/2JFTXb0Grf17HKgCsBVFRjbtCg=="
	output := PubAddrRetrievalStr(s)
	t.Log(output)
}
