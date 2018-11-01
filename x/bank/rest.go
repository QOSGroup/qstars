package bank

import (
	"io/ioutil"
	"net/http"

	sdk "github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/wire"

	"github.com/QOSGroup/qstars/client/lcd/lib"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cdc *wire.Codec, r *mux.Router) {
	//r.HandleFunc("/accounts/{address}/send", SendRequestHandlerFn(cdc, kb, cliCtx)).Methods("POST")
	r.HandleFunc("/accounts/{address}/send", func(w http.ResponseWriter, r *http.Request) {
		sb, err := NewSendBody(r)
		if err != nil {
			lib.HttpResponseWrapper(w, cdc, nil, err)
			return
		}

		result, err := sb.Send(cdc)
		lib.HttpResponseWrapper(w, cdc, result, err)
	}).Methods("POST")

	r.HandleFunc("/accounts/txSend", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			lib.HttpResponseWrapper(w, cdc, nil, err)
			return
		}

		result, err := TxSend(cdc, body)

		lib.HttpResponseWrapper(w, cdc, result, err)
	}).Methods("POST")
}

type sendBody struct {
	// fees is not used currently
	// Fees             sdk.Coin  `json="fees"`
	address       string `json:"-"`
	Amount        string `json:"amount"`
	PirvateKey    string `json:"privatekey"`
	ChainID       string `json:"chain_id"`
	AccountNumber int64  `json:"account_number"`
	Sequence      int64  `json:"sequence"`
	Gas           int64  `json:"gas"`
}

//send --from=GEPPkslt1Duwnb4B4W8OT1h311LYpo9GuJygHCE6mhH6iq1A17jIzMEzf6NiXUi6iGjDyoj9/GAhzSeyZqIzWg== --amount=3QSC1 --to=cosmosaccaddr120ws5500u0q8q75k70uetqp2xnysus5t4x9ug9 --sequence=1 --chain-id=test-chain-AE4XQo
//a:="{\"amount\":\"3QSC1\",\"privatekey\":\"GEPPkslt1Duwnb4B4W8OT1h311LYpo9GuJygHCE6mhH6iq1A17jIzMEzf6NiXUi6iGjDyoj9\",\"chain_id\":\"test-chain-AE4XQo\",\"account_number\":\"1\",\"sequence\":\"1\",\"gas\":\"1\"}"

func init() {
	RegisterWire(msgCdc)
}

func NewSendBody(r *http.Request) (*sendBody, error) {
	sb := &sendBody{}
	vars := mux.Vars(r)
	sb.address = vars["address"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	err = msgCdc.UnmarshalJSON(body, sb)
	if err != nil {
		return nil, err
	}

	return sb, nil
}
func (sb *sendBody) Send(cdc *wire.Codec) (*SendResult, error) {

	to, err := sdk.AccAddressFromBech32(sb.address)
	if err != nil {
		return nil, err
	}
	fromstr := sb.PirvateKey

	amount := sb.Amount
	// parse coins trying to be sent
	coins, err := sdk.ParseCoins(amount)
	chainid := sb.ChainID
	if err != nil {
		return nil, err
	}

	result, err := Send(cdc, fromstr, to, coins, chainid, NewSendOptions(
		gas(viper.GetInt64("gas")),
		fee(viper.GetString("fee"))))

	return result, err
}
