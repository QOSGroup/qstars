package main

import (
	"fmt"
	"github.com/QOSGroup/qstars/slim"
)

func AccountCreate() string {
	output := slim.AccountCreateStr()
	return output
}

func main() {
	output := AccountCreate()
	fmt.Println(output)
}
