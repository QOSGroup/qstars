package slim

import "testing"

func TestAccountCreate(t *testing.T) {
	output := AccountCreateStr()
	t.Log(output)
}

func TestAccountRecoverStr(t *testing.T) {
	mncode := "walk board image invite route below uncle toss celery negative hamster ordinary pink swift round student tip dragon rich such video sheriff just pull"
	output := AccountRecoverStr(mncode)
	t.Log(output)
}

func TestPubAddrRetrievalStr(t *testing.T) {
	s := "Me7Ts4jHgytbyq5ctu/lzgT630jKVxGt/Yb14HvHU4zfRHBf1KkdI7uo8H/vWQRFqJVLZeo5MEjUjDL08BWcPw=="
	output := PubAddrRetrievalStr(s)
	t.Log(output)
}
