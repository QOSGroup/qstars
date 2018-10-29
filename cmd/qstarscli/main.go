package main

import (
	"os"

	"github.com/QOSGroup/qbase/version"
	"github.com/QOSGroup/qstars/client"
	"github.com/QOSGroup/qstars/client/lcd"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/star"
	"github.com/QOSGroup/qstars/x/bank"
	authcmd "github.com/QOSGroup/qstars/x/auth"
	bankcmd "github.com/QOSGroup/qstars/x/bank"
	"github.com/QOSGroup/qstars/x/kvstore"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
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

	bank.NewBankStub().RegisterKVCdc(cdc)
	kvstore.NewKVStub().RegisterKVCdc(cdc)

	rootCmd := &cobra.Command{
		Use:   "cmd",
		Short: "qstars Command Line Interface(command)",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Name() == version.VersionCmd.Name() {
				return nil
			}
			cfg, err := config.InterceptLoadConfig()
			if err != nil {
				return err
			}
			config.CreateCLIContextTwo(cdc,cfg)
			return nil
		},
	}

	// TODO: Setup keybase, viper object, etc. to be passed into
	// the below functions and eliminate global vars, like we do
	// with the cdc.

	// add query/post commands (custom to binary)
	rootCmd.AddCommand(
		client.GetCommands(

			authcmd.GetAccountCmd(cdc),
			authcmd.CreateAccountCmd(cdc),
		)...)
	//
	rootCmd.AddCommand(
		client.PostCommands(
			bankcmd.SendTxCmd(cdc),
		)...)

	//


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
	executor := cli.PrepareBaseCmd(rootCmd, "BC", os.ExpandEnv("$HOME/.qstarscli"))



	err := executor.Execute()
	if err != nil {
		// Note: Handle with #870
		panic(err)
	}
}

