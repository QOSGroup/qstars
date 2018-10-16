package rest

import (
	"fmt"
	"github.com/QOSGroup/qstars/client/context"
	"github.com/QOSGroup/qstars/wire"
	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/libs/common"
	"io/ioutil"
	"net/http"
)

type SendKVBody struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	PrivateKey string `json:"privatekey"`
	ChainID    string `json:"chainid"`
}

var msgCdc = wire.NewCodec()

func init() {

}

// register REST routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *wire.Codec, storeName string) {
	r.HandleFunc(
		"/kv/{key}",
		QueryKRequestHandlerFnGet(storeName, cdc, cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/kv",
		QueryKVRequestHandlerFnSet(storeName, cdc, cliCtx),
	).Methods("POST")
}

// query accountREST Handler
func QueryKRequestHandlerFnGet(
	storeName string, cdc *wire.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]

		cliCtx := context.NewCLIContext().
			WithCodec(cdc)

		res, err := cliCtx.QueryStore(common.HexBytes(key), storeName)
		//res, err := clictx.Query(path)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("couldn't get. Error: %s", err.Error())))
			return
		} else if len(res) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("couldn't get. Error: response length is zero.")))
			return
		}

		w.Write([]byte(res))
	}
}

// query accountREST Handler
func QueryKVRequestHandlerFnSet(
	storeName string, cdc *wire.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var m SendKVBody
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		fmt.Println(string(body))
		m.Key = "key"
		m.Value = "value"
		o, err := msgCdc.MarshalJSON(m)
		fmt.Println(string(o))
		err = msgCdc.UnmarshalJSON(body, &m)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		//var priv ed25519.PrivKeyEd25519
		//bz := utility.Decbase64(m.PrivateKey)
		//copy(priv[:],bz)

		//_, addrben32 := utility.PubAddrRetrieval(m.PrivateKey)

		// build and sign the transaction, then broadcast to Tendermint
		//msg := BuildMsg(m.Key, m.Value, addrben32)

		//cliCtx := context.NewCLIContext().
		//	WithCodec(cdc)
		//txCtx := authctx.NewTxContextFromCLI().WithCodec(cdc)
		//txCtx.ChainID = m.ChainID
		//
		//response, err := utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg},priv)
		//if err != nil {
		//	w.WriteHeader(http.StatusBadRequest)
		//	w.Write([]byte(err.Error()))
		//	return
		//}

		//w.Write([]byte(response))
	}
}
