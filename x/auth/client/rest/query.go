package rest

import (
	"fmt"
	"net/http"

	"github.com/QOSGroup/qstars/client/context"
	sdk "github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/wire"
	"github.com/QOSGroup/qstars/x/auth"
	authcmd "github.com/QOSGroup/qstars/x/auth/client/cli"

	"github.com/QOSGroup/qstars/stub"
	"github.com/gorilla/mux"
)

// register REST routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *wire.Codec, storeName string) {
	r.HandleFunc(
		"/accounts/{address}",
		QueryAccountRequestHandlerFn(storeName, cdc, authcmd.GetAccountDecoder(cdc), cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/accounts",
		QueryAccountRequestHandlerFn2(storeName, cdc, authcmd.GetAccountDecoder(cdc), cliCtx),
	).Methods("POST")
}

// query accountREST Handler
func QueryAccountRequestHandlerFn2(
	storeName string, cdc *wire.Codec,
	decoder auth.AccountDecoder, cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		output := stub.AccountCreateStr()
		w.Write([]byte(output))
	}
}

// query accountREST Handler
func QueryAccountRequestHandlerFn(
	storeName string, cdc *wire.Codec,
	decoder auth.AccountDecoder, cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32addr := vars["address"]

		address, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		cliCtx := context.NewCLIContext().
			WithCodec(cdc).
			WithAccountDecoder(decoder)

		res, err := cliCtx.QueryQOSAccount(auth.AddressStoreKey(address), cliCtx.AccountStore)

		//res, err := cliCtx.QueryStore(auth.AddressStoreKey(addr), storeName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't query account. Error: %s", err.Error())))
			return
		}

		// the query will return empty if there is no data for this account
		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// decode the value

		account, err := cliCtx.AccDecoder(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't parse query result. Result: %s. Error: %s", res, err.Error())))
			return
		}

		// print out whole account
		output, err := wire.MarshalJSONIndent(cdc, account)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't marshall query result. Error: %s", err.Error())))
			return
		}

		w.Write(output)
	}
}
