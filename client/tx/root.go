package tx

import (
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/QOSGroup/qstars/client/context"
	"github.com/QOSGroup/qstars/wire"
)

// AddCommands adds a number of tx-query related subcommands
func AddCommands(cmd *cobra.Command, cdc *wire.Codec) {
	cmd.AddCommand(
		QueryTxCmd(cdc),
	)
}

// register REST routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *wire.Codec) {
	r.HandleFunc("/txs/{hash}", QueryTxRequestHandlerFn(cdc, cliCtx)).Methods("GET")
	// r.HandleFunc("/txs/sign", SignTxRequstHandler).Methods("POST")
	// r.HandleFunc("/txs/broadcast", BroadcastTxRequestHandler).Methods("POST")
}
