package auth

import (
	"github.com/QOSGroup/qstars/slim"
	"net/http"

	"github.com/QOSGroup/qstars/client/lcd/lib"
	"github.com/QOSGroup/qstars/wire"
	"github.com/gorilla/mux"
)

// RegisterRoutes register REST routes
func RegisterRoutes(cdc *wire.Codec, r *mux.Router) {
	r.HandleFunc(
		"/QOSaccounts/{address}",
		func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			bech32addr := vars["address"]
			result, err := QueryAccount(cdc, bech32addr)
			lib.HttpResponseWrapper(w, cdc, result, err)
		}).Methods("GET")

	r.HandleFunc(
		"/QSCaccounts/{address}",
		func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			bech32addr := vars["address"]
			result, err := QSCQueryAccount(cdc, bech32addr)
			lib.HttpResponseWrapper(w, cdc, result, err)
		}).Methods("GET")

	r.HandleFunc(
		"/accounts",
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			password := r.FormValue("password")
			acc := slim.AccountCreate(password)
			lib.HttpResponseWrapper(w, cdc, acc, nil)
		}).Methods("POST")
}
