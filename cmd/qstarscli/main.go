package main

import (
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qstars/client"
	"github.com/QOSGroup/qstars/client/lcd"
	"github.com/QOSGroup/qstars/star"
	authcmd "github.com/QOSGroup/qstars/x/auth/client/cli"
	bankcmd "github.com/QOSGroup/qstars/x/bank"
	"github.com/QOSGroup/qstars/x/kvstore"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
	"os"
)

// rootCmd is the entry point for this binary
var (
	rootCmd = &cobra.Command{
		Use:   "basecli",
		Short: "Basecoin light-client",
	}
)

func main() {
	// disable sorting
	cobra.EnableCommandSorting = false

	// get the codec
	cdc := star.MakeCodec()

	// TODO: Setup keybase, viper object, etc. to be passed into
	// the below functions and eliminate global vars, like we do
	// with the cdc.


	// add query/post commands (custom to binary)
	rootCmd.AddCommand(
		client.GetCommands(

			authcmd.GetAccountCmd("acc", cdc, authcmd.GetAccountDecoder(cdc)),
			authcmd.CreateAccountCmd(cdc),
		)...)
	//
	rootCmd.AddCommand(
		client.PostCommands(
			bankcmd.SendTxCmd(cdc),
		)...)

	//

	txs.RegisterCodec(cdc)
	kvstore.NewKVStub().RegisterKVCdc(cdc)
	rootCmd.AddCommand(
		client.PostCommands(
			kvstore.SendKVCmd(cdc),
			kvstore.GetKVCmd(cdc),
		)...)

	// add proxy, version and key info
	rootCmd.AddCommand(
		client.LineBreak,
		lcd.ServeCommand(cdc),
	//	keys.Commands(),
	//	client.LineBreak,
	//	version.VersionCmd,
	)

	// prepare and add flags
	executor := cli.PrepareMainCmd(rootCmd, "BC", os.ExpandEnv("$HOME/.basecli"))
	err := executor.Execute()
	if err != nil {
		// Note: Handle with #870
		panic(err)
	}
}
