package main

import (
	"fmt"
	"github.com/QOSGroup/qbase/types"
)

type AOETokens = types.BaseCoins

func main(){


	r0:=types.BaseCoins{&types.BaseCoin{"AOE",types.NewInt(500)}}




	r1:=types.BaseCoins{&types.BaseCoin{"AOE",types.NewInt(1500)}}



	rr:=r0.Plus(r1)


	fmt.Println(rr)
}
