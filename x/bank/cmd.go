package bank

import (
	"fmt"

	"github.com/QOSGroup/qstars/wire"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	qbasetxs "github.com/QOSGroup/qbase/txs"
	qbtype "github.com/QOSGroup/qbase/types"
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
			toStr := viper.GetString(flagTo)
			to, err := sdk.AccAddressFromBech32(toStr)
			if err != nil {
				return err
			}

			fromstr := viper.GetString(flagFrom) //Teddy changes

			amount := viper.GetString(flagAmount)
			// parse coins trying to be sent
			coins, err := sdk.ParseCoins(amount)
			if err != nil {
				return err
			}

			result, err := Send(cdc, fromstr, to, coins, NewSendOptions(
				gas(viper.GetInt64("gas")),
				fee(viper.GetString("fee"))))
			if err != nil {
				return err
			}

			output, err := wire.MarshalJSONIndent(cdc, result)
			if err != nil {
				return err
			}

			fmt.Println(string(output))
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

	tx.Receivers = append(tx.Receivers,receiver)
	tx.Senders = append(tx.Senders,sender)

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
