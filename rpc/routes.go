package rpc

import (
	rpc "github.com/tendermint/tendermint/rpc/lib/server"
)

type (
	ResultTestHealth struct{}
)

func MyHealth() (*ResultTestHealth, error) {
	print("000000000000000000000")
	return &ResultTestHealth{}, nil
}

// TODO: better system than "unsafe" prefix
// NOTE: Amino is registered in rpc/core/types/wire.go.
var Routes = map[string]*rpc.RPCFunc{
	// subscribe/unsubscribe are reserved for websocket events.
	"myhealth": rpc.NewRPCFunc(MyHealth, ""),
}
