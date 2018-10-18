package bank

import (
	"fmt"
	"os"

	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qstars/client/context"
	"github.com/QOSGroup/qstars/client/utils"
	qstarstypes "github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/utility"
	"github.com/QOSGroup/qstars/wire"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"

	qbasetxs "github.com/QOSGroup/qbase/txs"
	qbtype "github.com/QOSGroup/qbase/types"
	qos "github.com/QOSGroup/qos/account"
	txs "github.com/QOSGroup/qos/txs"
	sdk "github.com/QOSGroup/qstars/types"
)

const (
	flagTo     = "to"
	flagAmount = "amount"
	flagFrom   = "from"
)

// SendTxCmd will create a send tx and sign it with the given key.
func SendTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send",
		Short: "Create and sign a send tx",
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout)

			//if err := cliCtx.EnsureAccountExists(); err != nil {
			//	return err
			//}

			cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
			cdc.RegisterConcrete(&ed25519.PubKeyEd25519{}, "ed25519.PubKeyEd25519", nil)
			cdc.RegisterInterface((*account.Account)(nil), nil)
			cdc.RegisterConcrete(&qos.QOSAccount{}, "qbase/account/QOSAccount", nil)

			toStr := viper.GetString(flagTo)

			to, err := sdk.AccAddressFromBech32(toStr)
			if err != nil {
				return err
			}

			//Teddy changes
			fromstr := viper.GetString(flagFrom)

			var priv ed25519.PrivKeyEd25519
			bz := utility.Decbase64(fromstr)
			copy(priv[:], bz)
			_, addrben32 := utility.PubAddrRetrieval(fromstr)

			from, err := sdk.AccAddressFromBech32(addrben32)
			if err != nil {
				return err
			}

			// parse coins trying to be sent
			amount := viper.GetString(flagAmount)
			coins, err := sdk.ParseCoins(amount)
			if err != nil {
				return err
			}

			//from, err := cliCtx.GetFromAddress()
			//if err != nil {
			//	return err
			//}

			account, err := cliCtx.GetAccount(from)
			if err != nil {
				return err
			}
			var coins qstarstypes.Coins
			for _, qsc := range qacc.QscList {
				amount := qsc.Amount
				coins = append(coins, qstarstypes.NewCoin(qsc.Name, qstarstypes.NewInt(amount.Int64())))
			}
			// ensure account has enough coins
			// ensure account has enough coins
			var qcoins qstarstypes.Coins
			for _, qsc := range account.QscList {
				amount := qsc.Amount
				qcoins = append(qcoins, qstarstypes.NewCoin(qsc.Name, qstarstypes.NewInt(amount.Int64())))
			}

			if !qcoins.IsGTE(coins) {
				return errors.Errorf("Address %s doesn't have enough coins to pay for this transaction.", from)
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := BuildMsg(from, to, coins, cdc)

			response, err := utils.SendTx(cliCtx, cdc, msg, priv)
			fmt.Println(response)
			return nil
		},
	}

	cmd.Flags().String(flagTo, "", "Address to send coins")
	cmd.Flags().String(flagAmount, "", "Amount of coins to send")

	return cmd
}

func BuildMsg(from sdk.AccAddress, to sdk.AccAddress, coins sdk.Coins, cdc *wire.Codec) *qbasetxs.TxStd {

	tx := txs.TxTransform{}
	receiver := txs.AddrTrans{}
	receiver.Amount = qbtype.NewInt(coins[0].Amount.Int64())
	receiver.QscName = coins[0].Denom

	sender := txs.AddrTrans{}
	sender.Amount = qbtype.NewInt(coins[0].Amount.Int64())
	sender.QscName = coins[0].Denom

	tx.Receivers[0] = receiver
	tx.Senders[0] = sender

	stdTx := qbasetxs.TxStd{}
	//	stdTx.ITx = tx
	stdTx.ChainID = "chainid"

	return &stdTx
}

// MsgSend - high level transaction of the coin module
type MsgSend struct {
	Inputs  []Input  `json:"inputs"`
	Outputs []Output `json:"outputs"`
}

// NewMsgSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgSend(in []Input, out []Output) MsgSend {
	return MsgSend{Inputs: in, Outputs: out}
}
