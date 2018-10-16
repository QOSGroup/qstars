package rest

import (
	"io/ioutil"
	"net/http"

	"github.com/QOSGroup/qstars/client/context"
	"github.com/QOSGroup/qstars/crypto/keys"
	sdk "github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/wire"
	"github.com/QOSGroup/qstars/x/bank"

	"fmt"
	"github.com/QOSGroup/qstars/utility"
	authcmd "github.com/QOSGroup/qstars/x/auth/client/cli"
	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase) {
	r.HandleFunc("/accounts/{address}/send", SendRequestHandlerFn(cdc, kb, cliCtx)).Methods("POST")
}

type sendBody struct {
	// fees is not used currently
	// Fees             sdk.Coin  `json="fees"`
	Amount           string `json:"amount"`
	PirvateKey string    `json:"privatekey"`
	ChainID          string    `json:"chain_id"`
	AccountNumber    int64     `json:"account_number"`
	Sequence         int64     `json:"sequence"`
	Gas              int64     `json:"gas"`
}

//send --from=GEPPkslt1Duwnb4B4W8OT1h311LYpo9GuJygHCE6mhH6iq1A17jIzMEzf6NiXUi6iGjDyoj9/GAhzSeyZqIzWg== --amount=3QSC1 --to=cosmosaccaddr120ws5500u0q8q75k70uetqp2xnysus5t4x9ug9 --sequence=1 --chain-id=test-chain-AE4XQo
//a:="{\"amount\":\"3QSC1\",\"privatekey\":\"GEPPkslt1Duwnb4B4W8OT1h311LYpo9GuJygHCE6mhH6iq1A17jIzMEzf6NiXUi6iGjDyoj9\",\"chain_id\":\"test-chain-AE4XQo\",\"account_number\":\"1\",\"sequence\":\"1\",\"gas\":\"1\"}"

var msgCdc = wire.NewCodec()

func init() {
	bank.RegisterWire(msgCdc)
}

// SendRequestHandlerFn - http request handler to send coins to a address
func SendRequestHandlerFn(cdc *wire.Codec, kb keys.Keybase, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// collect data
		vars := mux.Vars(r)
		bech32addr := vars["address"]

		cliCtx.AccDecoder = authcmd.GetAccountDecoder(cdc)

		to, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		var m sendBody
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		err = msgCdc.UnmarshalJSON(body, &m)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		fromstr := m.PirvateKey

		//---------------------------------------------------------
		var priv ed25519.PrivKeyEd25519
		bz := utility.Decbase64(fromstr)
		copy(priv[:],bz)
		//Teddy changes
		_, addrben32 := utility.PubAddrRetrieval(fromstr)

		from, err := sdk.AccAddressFromBech32(addrben32)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		account, err := cliCtx.GetAccount(from)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		amount:= m.Amount
		coins, err := sdk.ParseCoins(amount)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// ensure account has enough coins
		if !account.GetCoins().IsGTE(coins) {
			w.WriteHeader(http.StatusInsufficientStorage)
			w.Write([]byte("Not enough money"))
			return
		}

		fmt.Println(to.String())

		//msg := client.BuildMsg(from, to, coins)

		//response,err := utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg},priv)

		//if err == nil{
		//	output, err := wire.MarshalJSONIndent(cdc, response)
		//	if err != nil {
		//		w.WriteHeader(http.StatusInternalServerError)
		//		w.Write([]byte(err.Error()))
		//		return
		//	}
		//	w.Write(output)
		//}else{
		//	w.WriteHeader(http.StatusInternalServerError)
		//	errstr := err.Error()
		//	w.Write([]byte(errstr))
		//}
	}
}
