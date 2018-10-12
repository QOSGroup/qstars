package rpc

import (
	cmn "github.com/tendermint/tendermint/libs/common"
	"google.golang.org/grpc"
	"net"
	"time"
)

// StartGRPCClient dials the gRPC server using protoAddr and returns a new
// BroadcastAPIClient.
func StartGRPCClient(protoAddr string) QSCAppInterfaceClient {
	conn, err := grpc.Dial(protoAddr, grpc.WithInsecure(), grpc.WithDialer(dialerFunc))
	if err != nil {
		panic(err)
	}
	return NewQSCAppInterfaceClient(conn)
}

func dialerFunc(addr string, timeout time.Duration) (net.Conn, error) {
	return cmn.Connect(addr)
}
