package slim

import "testing"

func TestAccountCreate(t *testing.T) {
	output := AccountCreateStr()
	t.Log(output)
}

func TestAccountRecoverStr(t *testing.T) {
	mncode := "connect piano supply broken stand answer awesome solid concert quit glad prize fiction wreck dilemma element leisure hedgehog hedgehog speak decade someone disagree near"
	output := AccountRecoverStr(mncode)
	t.Log(output)
}

func TestPubAddrRetrievalStr(t *testing.T) {
	s := "oyiJEEBTGUOQ9UL6pcVKTckHm+8QaJlyb4cARs8aVfU6JLdI1eSO0+YDW4M+uYW+u/5ZFp6ybEEdYoOF7Kz37Hwr8eYd"
	output := PubAddrRetrievalStr(s)
	t.Log(output)
}
