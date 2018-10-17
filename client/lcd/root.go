package lcd

import (
	"net/http"
	"os"

	client "github.com/QOSGroup/qstars/client"
	"github.com/QOSGroup/qstars/client/context"
	"github.com/QOSGroup/qstars/wire"
	auth "github.com/QOSGroup/qstars/x/auth"
	bank "github.com/QOSGroup/qstars/x/bank/client/rest"
	kvstore "github.com/QOSGroup/qstars/x/kvstore/rest"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	tmserver "github.com/tendermint/tendermint/rpc/lib/server"
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

			listener, err := tmserver.StartHTTPServer(
				listenAddr, handler, logger,
				tmserver.Config{MaxOpenConnections: maxOpen},
			)
			if err != nil {
				return err
			}

			logger.Info("REST server started")

			// wait forever and cleanup
			cmn.TrapSignal(func() {
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
	cmd.Flags().Int(flagMaxOpenConnections, 1000, "The number of maximum open connections")

	return cmd
}

func createHandler(cdc *wire.Codec) http.Handler {
	r := mux.NewRouter()

	cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout)

	CLIVersionRegisterRoutes(cliCtx, r)
	NodeVersionRegisterRoutes(cliCtx, r)

	auth.RegisterRoutes(cdc, r)
	bank.RegisterRoutes(cliCtx, r, cdc, nil)
	kvstore.RegisterRoutes(cliCtx, r, cdc, "main")
	//ibc.RegisterRoutes(cliCtx, r, cdc, kb)
	//stake.RegisterRoutes(cliCtx, r, cdc, kb)
	//slashing.RegisterRoutes(cliCtx, r, cdc, kb)
	//gov.RegisterRoutes(cliCtx, r, cdc)

	return r
}

type HTTPHandler interface {
	RouterRegister(r *mux.Route)
	ParseRequest(req *http.Request) (interface{}, error)
}
