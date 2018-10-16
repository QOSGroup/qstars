package rpc

import (
	"testing"
)

/*
ws = new WebSocket("ws://127.0.0.1:8080/echo");
	ws.onopen = function(evt) {

	}
	ws.onmessage = function(evt) {

	}
	ws.onerror = function(evt) {

	}
	document.getElementById("send").onclick = function(evt) {
	if (!ws) {
		return false
	}
	ws.send(input.value);
	return;
	};

*/

func TestRunWSServer(t *testing.T) {
	//var testconfig cfg.Config
	//testconfig.RPC = new(cfg.RPCConfig)
	//testconfig.RPC.ListenAddress = "tcp://127.0.0.1:23232"
	//testconfig.RPC.GRPCMaxOpenConnections = 10
	//testconfig.RPC.GRPCListenAddress = "tcp://127.0.0.1:23232"
	//
	//logger := log.TestingLogger()
	//_, err := StartHTTPWS(&testconfig, logger)
	//
	//addr := flag.String("addr", "localhost:23232", "http service address")
	//WSCall("ws", "websocket", addr)
	//
	//if err != nil {
	//	fmt.Errorf("%v", err)
	//	return
	//}
	//for {
	//	runtime.Gosched()
	//}
	//WhatMe()
}
