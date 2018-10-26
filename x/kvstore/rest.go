package kvstore

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/QOSGroup/qstars/client/context"
	"github.com/QOSGroup/qstars/client/lcd/lib"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/wire"
	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/libs/common"
)

var msgCdc = wire.NewCodec()

func init() {

}

// register REST routes
func RegisterRoutes(r *mux.Router, cdc *wire.Codec, storeName string) {
	r.HandleFunc(
		"/kv/{key}",
		func(w http.ResponseWriter, r *http.Request) {
			skr, err := NewGetKVReq(r)
			if err != nil {
				lib.HttpResponseWrapper(w, cdc, nil, err)
				return
			}
			result, err := skr.GetKV(storeName, cdc)
			lib.HttpResponseWrapper(w, cdc, result, err)
		}).Methods("GET")

	r.HandleFunc("/kv", func(w http.ResponseWriter, r *http.Request) {
		skr, err := NewSendKVReq(r)
		if err != nil {
			lib.HttpResponseWrapper(w, cdc, nil, err)
			return
		}
		result, err := skr.SendKV(storeName, cdc)
		lib.HttpResponseWrapper(w, cdc, result, err)
	}).Methods("POST")
}

type sendKVReq struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	PrivateKey string `json:"privatekey"`
	ChainID    string `json:"chainid"`
}

func NewSendKVReq(r *http.Request) (*sendKVReq, error) {
	skr := &sendKVReq{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	err = msgCdc.UnmarshalJSON(body, skr)
	if err != nil {
		return nil, err
	}

	return skr, nil
}

func (skr *sendKVReq) SendKV(storeName string, cdc *wire.Codec) (*ResultSendKV, error) {
	opts, err := NewSendKVOption(
		SendKVOptionChainID(skr.ChainID),
	)
	if err != nil {
		return nil, err
	}

	result, err := SendKV(cdc, skr.PrivateKey, skr.Key, skr.Value, opts)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type getKVReq struct {
	Key string `url:"key"`
}

func NewGetKVReq(r *http.Request) (*getKVReq, error) {
	skr := &getKVReq{}
	vars := mux.Vars(r)
	skr.Key = vars["key"]

	return skr, nil
}

func (skr *getKVReq) GetKV(storeName string, cdc *wire.Codec) (*ResultGetKV, error) {
	result, err := GetKV(cdc, skr.Key, nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// query accountREST Handler
func QueryKRequestHandlerFnGet(
	storeName string, cdc *wire.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]

		cliCtx := config.GetCLIContext().QSCCliContext

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
