package bank

import (
	"fmt"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/wire"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"

	sdk "github.com/QOSGroup/qstars/types"
	qbasetypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/types"
)

const (
	flagTo      = "to"
	flagFromAmount  = "fromamount"
	flagToAmount  = "toamount"
	flagFrom    = "from"
	flagCommand = "command"
	flagChainid = "chainid"
)

// SendTxCmd will create a send tx and sign it with the given key.
func SendTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send",
		Short: "Create and sign a send tx",
		RunE: func(cmd *cobra.Command, args []string) error {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
				}
			}()
			toStr := viper.GetString(flagTo)
			fromstr := viper.GetString(flagFrom) //Teddy changes

			fromamount := viper.GetString(flagFromAmount)
			toamount := viper.GetString(flagToAmount)

			directTOQOS := config.GetCLIContext().Config.DirectTOQOS

			a,b,c,d := formatInput(cdc,toStr,fromstr,fromamount,toamount)


			if directTOQOS == true {
				result, err := MultiSendDirect(cdc,a,b,c,d);
				if err != nil {
					fmt.Println(err)
					return err
				}
				output, err := wire.MarshalJSONIndent(cdc, result)
				if err != nil {
					return err
				}

				fmt.Println(string(output))
				return nil
			}else {
				result, err := MultiSendViaQStars(cdc,a,b,c,d);
				if err != nil {
					fmt.Println(err)
					return err
				}
				output, err := wire.MarshalJSONIndent(cdc, result)
				if err != nil {
					return err
				}

				fmt.Println(string(output))
				return nil
			}
			//result, err := Send(cdc, fromstr, to, coins, NewSendOptions(
			//	gas(viper.GetInt64("gas")),
			//	fee(viper.GetString("fee"))))
			//if err != nil {
			//	return err
			//}


		},
	}

	cmd.Flags().String(flagTo, "", "Addresses to send coins")
	cmd.Flags().String(flagFromAmount, "", "Amount of coins to send")
	cmd.Flags().String(flagToAmount, "", "Amount of coins to receive")

	return cmd
}

// ApproveCmd will create a approve tx and sign it with the given key.
func ApproveCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve",
		Short: "create increase decrease use and cancel approve",
		RunE: func(cmd *cobra.Command, args []string) error {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
				}
			}()
			toStr := viper.GetString(flagTo)
			fromstr := viper.GetString(flagFrom) //Teddy changes
			amount := viper.GetString(flagFromAmount)
			// parse coins trying to be sent
			coins, err := sdk.ParseCoins(amount)
			chainid := viper.GetString(flagChainid)
			if err != nil {
				return err
			}

			command := viper.GetString(flagCommand)
			result, err := Approve(cdc, command, fromstr, toStr, coins, chainid, NewSendOptions(
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
	cmd.Flags().String(flagFromAmount, "", "Amount of coins to send")
	cmd.Flags().String(flagToAmount, "", "Amount of coins to receive")
	cmd.Flags().String(flagCommand, "", "client command, for approve: create,increase,decrease,use,cancel")

	return cmd
}

//func BuildMsg(from qbtype.Address, to qbtype.Address, coins sdk.Coins, cdc *wire.Codec) *qbasetxs.TxStd {
//
//	tx := txs.TxTransform{}
//	receiver := txs.AddrTrans{}
//	receiver.Amount = qbtype.NewInt(coins[0].Amount.Int64())
//	receiver.QscName = coins[0].Denom
//
//	sender := txs.AddrTrans{}
//	sender.Amount = qbtype.NewInt(coins[0].Amount.Int64())
//	sender.QscName = coins[0].Denom
//
//	tx.Receivers = append(tx.Receivers,receiver)
//	tx.Senders = append(tx.Senders,sender)
//
//	stdTx := qbasetxs.TxStd{}
//	//	stdTx.ITx = tx
//	stdTx.ChainID = "chainid"
//
//	return &stdTx
//}

// MsgSend - high level transaction of the coin module
type MsgSend struct {
	Inputs  []Input  `json:"inputs"`
	Outputs []Output `json:"outputs"`
}

// NewMsgSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgSend(in []Input, out []Output) MsgSend {
	return MsgSend{Inputs: in, Outputs: out}
}

func formatInput(cdc *wire.Codec, toaddressesStr string,fromaddressesStr string,fromamountstr string,toamountstr string) ([]string,[]qbasetypes.Address,[]types.Coins,[]types.Coins){

	tostrs := []qbasetypes.Address {}

	toaddresses := strings.Split(toaddressesStr,";")
	fmt.Println("to address:",toaddresses)
	for i:=0;i< len(toaddresses);i++{
		to, err := sdk.AccAddressFromBech32(toaddresses[i])
		if err != nil {
			fmt.Println(err)
		}
		tostrs = append(tostrs, to)
	}

	fromstrs := []string{}

	fromaddresses := strings.Split(fromaddressesStr,";")
	fmt.Println("private key:",fromaddresses)
	for i:=0;i< len(fromaddresses);i++{
		fromstrs = append(fromstrs, fromaddresses[i])
	}

	scoins :=[]types.Coins {}

	// parse coins trying to be sent
	fromamounts := strings.Split(fromamountstr,";")
	fmt.Println("from mounts:",fromamounts)
	for i:=0;i< len(fromamounts);i++ {
		secointmp, err := sdk.ParseCoins(fromamounts[i])
		if err != nil {
			fmt.Println(err)
		}
		scoins = append(scoins, secointmp)
	}


	rcoins :=[]types.Coins {}

	toamounts := strings.Split(toamountstr,";")
	fmt.Println("to mounts:",toamounts)
	for i:=0;i< len(toamounts);i++ {
		// parse coins trying to be sent
		recointmp, err := sdk.ParseCoins(toamounts[i])
		if err != nil {
			fmt.Println(err)
		}
		rcoins = append(rcoins, recointmp)
	}
	return fromstrs, tostrs, scoins, rcoins

}