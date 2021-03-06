package lcd

import (
	"net/http"
	"os"

	"github.com/QOSGroup/qstars/client"

	"github.com/QOSGroup/qstars/wire"
	"github.com/QOSGroup/qstars/x/auth"
	"github.com/QOSGroup/qstars/x/bank"
	"github.com/QOSGroup/qstars/x/kvstore"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	tmserver "github.com/tendermint/tendermint/rpc/lib/server"
	rpcserver "github.com/tendermint/tendermint/rpc/lib/server"
)

// ServeCommand will generate a long-running rest server
// (aka Light Client Daemon) that exposes functionality similar
// to the cli, but over rest
func ServeCommand(cdc *wire.Codec) *cobra.Command {
	flagListenAddr := "laddr"
	flagCORS := "cors"
	flagMaxOpenConnections := "max-open"

	cmd := &cobra.Command{
		Use:   "rest-server",
		Short: "Start LCD (light-client daemon), a local REST server",
		RunE: func(cmd *cobra.Command, args []string) error {
			listenAddr := viper.GetString(flagListenAddr)
			handler := createHandler(cdc)
			logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "rest-server")
			maxOpen := viper.GetInt(flagMaxOpenConnections)
			listener, err := rpcserver.Listen(
				listenAddr,
				&rpcserver.Config{MaxOpenConnections: maxOpen},
			)
			err = tmserver.StartHTTPServer(
				listener, handler, logger,&rpcserver.Config{MaxOpenConnections: maxOpen},)
			if err != nil {
				return err
			}

			logger.Info("REST server started")

			// wait forever and cleanup
			cmn.TrapSignal(logger,func() {
				err := listener.Close()
				logger.Error("error closing listener", "err", err)
			})

			return nil
		},
	}

	cmd.Flags().String(flagListenAddr, "tcp://localhost:1317", "The address for the server to listen on")
	cmd.Flags().String(flagCORS, "", "Set the domains that can make CORS requests (* for all)")
	cmd.Flags().String(client.FlagChainID, "", "The chain ID to connect to")
	cmd.Flags().String(client.FlagNode, "tcp://localhost:26657", "Address of the node to connect to")
	cmd.Flags().Bool(client.FlagTrustNode, true, "Don't verify proofs for responses")
	cmd.Flags().Int(flagMaxOpenConnections, 1000, "The number of maximum open connections")

	return cmd
}

func createHandler(cdc *wire.Codec) http.Handler {
	r := mux.NewRouter()


	CLIVersionRegisterRoutes( cdc,r)
	NodeVersionRegisterRoutes( cdc,r)

	auth.RegisterRoutes(cdc, r)
	bank.RegisterRoutes(cdc, r)
	kvstore.RegisterRoutes( cdc,r)

	return r
}

type HTTPHandler interface {
	RouterRegister(r *mux.Route)
	ParseRequest(req *http.Request) (interface{}, error)
}