package main

import (
	"github.com/QOSGroup/qstars/star"
	"github.com/QOSGroup/qstars/x/jianqian/advertisers"
	"github.com/QOSGroup/qstars/x/jianqian/auction"
	"github.com/QOSGroup/qstars/x/jianqian/recharge"

	"github.com/QOSGroup/qstars/x/jianqian/buyad"
	"github.com/QOSGroup/qstars/x/jianqian/investad"

	"github.com/QOSGroup/qstars/x/jianqian/article"
	"github.com/QOSGroup/qstars/x/jianqian/coins"

	qbasecli "github.com/QOSGroup/qbase/client"
	"os"

	"github.com/QOSGroup/qbase/version"
	"github.com/QOSGroup/qstars/client"
	"github.com/QOSGroup/qstars/client/lcd"
	"github.com/QOSGroup/qstars/config"
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

	rootCmd = &cobra.Command{
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
			config.CreateCLIContextTwo(cdc, cfg)
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
			bankcmd.ApproveCmd(cdc),
		)...)

	//
	rootCmd.AddCommand(
		client.PostCommands(
			kvstore.SendKVCmd(cdc),
			kvstore.GetKVCmd(cdc),
		)...)

	//
	rootCmd.AddCommand(
		client.PostCommands(
			coins.DispatchAOECmd(cdc),
			article.NewArticleCmd(cdc),
			article.QueryArticleCmd(cdc),
			coins.QueryBlanceCmd(cdc),
			recharge.RechargeCmd(cdc),
			auction.NewAuctionCmd(cdc),
			advertisers.AdvertisersCmd(cdc),
		)...)

	// add proxy, version and key info
	rootCmd.AddCommand(
		client.LineBreak,
		lcd.ServeCommand(cdc),
	//	keys.Commands(),
	//	client.LineBreak,
	//	version.VersionCmd,
	)

	rootCmd.AddCommand(
		buyad.BuyadCmd(cdc),
		investad.InvestadCmd(cdc),
	)

	// add commands provided by qbase
	qbaseCmd := &cobra.Command{
		Use:   "qbase",
		Short: "qbase commands",
	}

	qbaseCmd.AddCommand(
		qbasecli.QueryCommand(cdc),
		qbasecli.KeysCommand(cdc),
		qbasecli.TendermintCommand(cdc),
		qbasecli.TxCommand(),
	)
	rootCmd.AddCommand(qbaseCmd)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "BC", os.ExpandEnv("$HOME/.qstarscli"))

	err := executor.Execute()
	if err != nil {
		// Note: Handle with #870
		panic(err)
	}
}
