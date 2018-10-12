package bank

import (
				"github.com/QOSGroup/qstars/wire"
					"github.com/spf13/cobra"
			)

const (
	flagTo     = "to"
	flagAmount = "amount"
	flagFrom = "from"
)

// SendTxCmd will create a send tx and sign it with the given key.
func SendTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send",
		Short: "Create and sign a send tx",
		RunE: func(cmd *cobra.Command, args []string) error {
			//txCtx := authctx.NewTxContextFromCLI().WithCodec(cdc)
			//cliCtx := context.NewCLIContext().
			//	WithCodec(cdc).
			//	WithLogger(os.Stdout).
			//	WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			//
			////if err := cliCtx.EnsureAccountExists(); err != nil {
			////	return err
			////}
			//
			//toStr := viper.GetString(flagTo)
			//
			//to, err := sdk.AccAddressFromBech32(toStr)
			//if err != nil {
			//	return err
			//}
			//
			////Teddy changes
			//fromstr := viper.GetString(flagFrom)
			//
			//var priv ed25519.PrivKeyEd25519
			//bz := utility.Decbase64(fromstr)
			//copy(priv[:],bz)
			//_, addrben32 := utility.PubAddrRetrieval(fromstr)
			//
			//from, err := sdk.AccAddressFromBech32(addrben32)
			//if err != nil {
			//	return err
			//}
			//
			//// parse coins trying to be sent
			//amount := viper.GetString(flagAmount)
			//coins, err := sdk.ParseCoins(amount)
			//if err != nil {
			//	return err
			//}
			//
			////from, err := cliCtx.GetFromAddress()
			////if err != nil {
			////	return err
			////}
			//
			//account, err := cliCtx.GetAccount(from)
			//if err != nil {
			//	return err
			//}
			//
			//// ensure account has enough coins
			//if !account.GetCoins().IsGTE(coins) {
			//	return errors.Errorf("Address %s doesn't have enough coins to pay for this transaction.", from)
			//}
			//
			//// build and sign the transaction, then broadcast to Tendermint
			////msg := client.BuildMsg(from, to, coins)
			//
			////response,err := utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg}, priv)
			//fmt.Println("")
			return nil
		},
	}

	cmd.Flags().String(flagTo, "", "Address to send coins")
	cmd.Flags().String(flagAmount, "", "Amount of coins to send")

	return cmd
}
