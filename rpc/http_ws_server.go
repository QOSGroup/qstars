package rpc

import (
	"github.com/tendermint/go-amino"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/rpc/lib/server"
	"github.com/tendermint/tendermint/types"
	"net"
	"net/http"
)

func StartHTTPWS(config *cfg.Config, logger log.Logger) (listener net.Listener, err error) {

	coreCodec := amino.NewCodec()
	//ctypes.RegisterAmino(coreCodec)

	eventBus := types.NewEventBus()
	eventBus.SetLogger(logger.With("module", "events"))

	listenAddrs := splitAndTrimEmpty(config.RPC.ListenAddress, ",", " ")
	// we may expose the rpc over both a unix and tcp socket
	listeners := make([]net.Listener, len(listenAddrs))
	for i, listenAddr := range listenAddrs {
		mux := http.NewServeMux()
		rpcLogger := logger.With("module", "rpc-server")

		wm := rpcserver.NewWebsocketManager(Routes, coreCodec, rpcserver.EventSubscriber(eventBus))

		//wm.SetLogger(rpcLogger.With("protocol", "websocket"))
		mux.HandleFunc("/websocket", wm.WebsocketHandler)
		rpcserver.RegisterRPCFuncs(mux, Routes, coreCodec, rpcLogger)
		listener, err := rpcserver.StartHTTPServer(
			listenAddr,
			mux,
			rpcLogger,
			rpcserver.Config{MaxOpenConnections: config.RPC.MaxOpenConnections},
		)
		if err != nil {

			return nil, err
		}
		listeners[i] = listener
	}
	return listener, nil
}
