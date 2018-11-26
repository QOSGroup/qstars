package jsdk

import (
	"github.com/QOSGroup/qbase/version"
	"github.com/QOSGroup/qstars/client"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/star"
	sdk "github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/wire"
	"github.com/QOSGroup/qstars/x/bank"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
	"os"
)

var CONF *config.CLIConfig
var CDC = star.MakeCodec()

func InitJNI() {
	//CDC := star.MakeCodec()
	rootCmd := &cobra.Command{
		Use:   "cmd",
		Short: "qstars Command Line Interface(command)",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Name() == version.VersionCmd.Name() {
				return nil
			}
			CONF, err1 := config.InterceptLoadConfig()
			if err1 != nil {
				return err1
			}
			config.CreateCLIContextTwo(CDC, CONF)
			return nil
		},
	}
	rootCmd.AddCommand(
		client.PostCommands(MimicCmd(CDC))...)

	var newarg []string
	newarg = append(newarg, os.Args[0])
	newarg = append(newarg, "mimic")
	os.Args = newarg
	executor := cli.PrepareBaseCmd(rootCmd, "BC", os.ExpandEnv("$HOME/.qstarscli"))

	err := executor.Execute()
	if err != nil {
		// Note: Handle with #870
		panic(err)
	}
}

// mimic cli for JNI
func MimicCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mimic",
		Short: "JNI mimic",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	return cmd
}

func SendByJNI(fromStr string, toStr1 string, coinstr string) string {

	// (*SendResult, error)

	toStr, err := sdk.AccAddressFromBech32(toStr1)
	if err != nil {
		return ""
	}

	// parse coins trying to be sent
	coins, err := sdk.ParseCoins(coinstr)
	if err != nil {
		return ""
	}
	result, err := bank.Send(CDC, fromStr, toStr, coins, nil)
	output, err := CDC.MarshalJSON(result)
	if err != nil {
		return err.Error()
	}
	return string(output)
}

//test JNI by QueryAccount function at first
//func QueryAccbyJNI(addr string) string {
//	acc, _ := auth.QueryAccount(CDC, addr)
//	output, _ := wire.MarshalJSONIndent(CDC, acc)
//	return string(output)
//}
