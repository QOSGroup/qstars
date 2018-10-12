package rpc

import (
	"fmt"
	cfg "github.com/tendermint/tendermint/config"
	context "golang.org/x/net/context"
	"golang.org/x/net/netutil"
	"google.golang.org/grpc"
	"net"
	"strings"
)

func StartAPPGRPC(config *cfg.Config) ([]net.Listener, error) {

	listenAddrs := splitAndTrimEmpty(config.RPC.ListenAddress, ",", " ")
	listeners := make([]net.Listener, len(listenAddrs))

	grpcListenAddr := config.RPC.GRPCListenAddress
	if grpcListenAddr != "" {
		listener, err := StartGRPCServer(
			grpcListenAddr,
			Config{
				MaxOpenConnections: config.RPC.GRPCMaxOpenConnections,
			},
		)
		if err != nil {
			return nil, err
		}
		listeners = append(listeners, listener)
	}
	return listeners, nil

}

// splitAndTrimEmpty slices s into all subslices separated by sep and returns a
// slice of the string s with all leading and trailing Unicode code points
// contained in cutset removed. If sep is empty, SplitAndTrim splits after each
// UTF-8 sequence. First part is equivalent to strings.SplitN with a count of
// -1.  also filter out empty strings, only return non-empty strings.

func splitAndTrimEmpty(s, sep, cutset string) []string {
	if s == "" {
		return []string{}
	}

	spl := strings.Split(s, sep)
	nonEmptyStrings := make([]string, 0, len(spl))
	for i := 0; i < len(spl); i++ {
		element := strings.Trim(spl[i], cutset)
		if element != "" {
			nonEmptyStrings = append(nonEmptyStrings, element)
		}
	}
	return nonEmptyStrings
}

// Config is an gRPC server configuration.
type Config struct {
	MaxOpenConnections int
}

type GreeterClientAPI struct {
}

var _ QSCAppInterfaceServer = (*GreeterClientAPI)(nil)

// set key-value
func (a GreeterClientAPI) QSCKVStoreSet(ctx context.Context, in *KVSetRequest) (*GeneralReply, error) {
	var reply GeneralReply
	return &reply, nil
}

// get key-value
func (a GreeterClientAPI) QSCKVStoreGet(ctx context.Context, in *KVSetRequest) (*KVGetReply, error) {
	var reply KVGetReply
	return &reply, nil

}

// query account
func (a GreeterClientAPI) QSCQueryAccount(ctx context.Context, in *QueryRequest) (*QueryRequest, error) {
	var reply QueryRequest
	return &reply, nil

}

// mint coin
func (a GreeterClientAPI) QSCMintCoin(ctx context.Context, in *MintRequest) (*GeneralReply, error) {
	print("hello: " + in.Addr)
	var reply GeneralReply
	return &reply, nil

}

// QSC transfer
func (a GreeterClientAPI) QSCtransfer(ctx context.Context, in *CoinTransferRequest) (*GeneralReply, error) {
	var reply GeneralReply
	return &reply, nil

}

// QOS to QSC transfer
func (a GreeterClientAPI) QOStoQSCtransfer(ctx context.Context, in *CoinTransferRequest) (*GeneralReply, error) {
	var reply GeneralReply
	return &reply, nil

}

// QSC to QOS transfer
func (a GreeterClientAPI) QSCtoQOStransfer(ctx context.Context, in *CoinTransferRequest) (*GeneralReply, error) {
	var reply GeneralReply
	return &reply, nil

}

// StartGRPCServer starts a new gRPC BroadcastAPIServer, listening on
// protoAddr, in a goroutine. Returns a listener and an error, if it fails to
// parse an address.
func StartGRPCServer(protoAddr string, config Config) (net.Listener, error) {
	parts := strings.SplitN(protoAddr, "://", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid listen address for grpc server (did you forget a tcp:// prefix?) : %s", protoAddr)
	}
	proto, addr := parts[0], parts[1]
	ln, err := net.Listen(proto, addr)
	if err != nil {
		return nil, err
	}
	if config.MaxOpenConnections > 0 {
		ln = netutil.LimitListener(ln, config.MaxOpenConnections)
	}

	grpcServer := grpc.NewServer()
	RegisterBroadcastAPIServer(grpcServer, &GreeterClientAPI{})
	go grpcServer.Serve(ln) // nolint: errcheck

	return ln, nil
}

func RegisterBroadcastAPIServer(s *grpc.Server, srv QSCAppInterfaceServer) {
	s.RegisterService(&_QSCAppInterface_serviceDesc, srv)
}

func WhatMe() {

}
