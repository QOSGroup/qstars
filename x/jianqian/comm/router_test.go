package comm

import "testing"

func TestGetStruct(t *testing.T) {

	args := []string{"address103cfpkrsw8yu78gem6mg26rkyqsfmnr5sj0hvp", "ATOM", "2", "10000"}
	//args:=[]string{"address103cfpkrsw8yu78gem6mg26rkyqsfmnr5sj0hvp","1000","aaaaa","atom"}

	tx, err := getStruct(ExtractTxFlag, args,nil)

	t.Log(err)

	t.Error(tx)

}
