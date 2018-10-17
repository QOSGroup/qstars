package auth

import (
	"net/http"

	"github.com/QOSGroup/qstars/client/lcd/lib"
	"github.com/QOSGroup/qstars/stub"
	"github.com/QOSGroup/qstars/wire"
	"github.com/gorilla/mux"
)

// RegisterRoutes register REST routes
func RegisterRoutes(cdc *wire.Codec, r *mux.Router) {
	r.HandleFunc(
		"/accounts/{address}",
		func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			bech32addr := vars["address"]
			result, err := QueryAccount(cdc, bech32addr)
			lib.HttpResponseWrapper(w, cdc, result, err)
		}).Methods("GET")

	r.HandleFunc(
		"/accounts",
		func(w http.ResponseWriter, r *http.Request) {
			acc := stub.AccountCreate()
			lib.HttpResponseWrapper(w, cdc, acc, nil)
		}).Methods("POST")
}
