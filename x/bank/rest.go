package bank

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/QOSGroup/qstars/client/context"
	"github.com/QOSGroup/qstars/crypto/keys"
	qstarstypes "github.com/QOSGroup/qstars/types"
	sdk "github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/wire"

	"fmt"

	"github.com/QOSGroup/qstars/client/lcd/lib"
	"github.com/QOSGroup/qstars/utility"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase) {
	//r.HandleFunc("/accounts/{address}/send", SendRequestHandlerFn(cdc, kb, cliCtx)).Methods("POST")
	r.HandleFunc("/accounts/{address}/send", func(w http.ResponseWriter, r *http.Request) {
		sb, err := NewSendBody(r)
		if err != nil {
			lib.HttpResponseWrapper(w, cdc, nil, err)
		}

		result, err := sb.Send(cdc, kb, cliCtx)
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

	var m sendBody
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	err = msgCdc.UnmarshalJSON(body, &m)
	if err != nil {
		return nil, err
	}

	return sb, nil
}
func (sb *sendBody) Send(cdc *wire.Codec, kb keys.Keybase, cliCtx context.CLIContext) (*SendResult, error) {
	to, err := sdk.AccAddressFromBech32(sb.address)
	if err != nil {
		return nil, err
	}
	fromstr := sb.PirvateKey
	//---------------------------------------------------------
	var priv ed25519.PrivKeyEd25519
	bz := utility.Decbase64(fromstr)
	copy(priv[:], bz)
	//Teddy changes
	_, addrben32 := utility.PubAddrRetrieval(fromstr)

	from, err := sdk.AccAddressFromBech32(addrben32)
	if err != nil {
		return nil, fmt.Errorf("no auth,%s", err.Error())
	}

	account, err := cliCtx.GetAccount(from)
	if err != nil {
		return nil, fmt.Errorf("no auth,%s", err.Error())
	}

	amount := sb.Amount
	coins, err := sdk.ParseCoins(amount)
	if err != nil {
		return nil, fmt.Errorf("amount不支持,%s", err.Error())
	}

	// ensure account has enough coins
	var qcoins qstarstypes.Coins
	for _, qsc := range account.QscList {
		amount := qsc.Amount
		qcoins = append(qcoins, qstarstypes.NewCoin(qsc.Name, qstarstypes.NewInt(amount.Int64())))
	}

	if !qcoins.IsGTE(coins) {
		return nil, errors.New("Not enough money ")
	}

	result, err := Send(cdc, fromstr, to, coins, NewSendOptions(
		gas(viper.GetInt64("gas")),
		fee(viper.GetString("fee"))))

	return result, err

}
