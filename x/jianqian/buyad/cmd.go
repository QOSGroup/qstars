package buyad

import (
	"errors"
	"fmt"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/wire"
	"github.com/QOSGroup/qstars/x/jianqian"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

const (
	flagBuyer       = "buyer"
	flagCoins       = "coins"
	flagArticleHash = "articleHash"
	flagChainid     = "chainid"
)

// BuyadCmd will create a send tx and sign it with the given key.
func BuyadCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buyad",
		Short: "buy ad",
	}

	cmd.AddCommand(
		buyadCmd(cdc),
		queryBuyadCmd(cdc),
	)

	return cmd
}

//func getQOSAcc(address []byte, cdc *wire.Codec) (*qosaccount.QOSAccount, error) {
//	return config.GetCLIContext().QOSCliContext.GetAccount(address, cdc)
//}
//
//func getQSCAcc(address []byte, cdc *wire.Codec) (*qosaccount.QOSAccount, error) {
//	return config.GetCLIContext().QSCCliContext.GetAccount(address, cdc)
//}

func buyadCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buyad",
		Short: "buy ad",
		RunE: func(cmd *cobra.Command, args []string) error {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
				}
			}()

			articleHash := viper.GetString(flagArticleHash)

			result := BuyAd(cdc, articleHash)

			log.Printf(result)

			return nil
		},
	}

	cmd.Flags().String(flagArticleHash, "", "articleHash")
	cmd.Flags().String(flagChainid, "", "Chainid")
	cmd.Flags().String(flagBuyer, "", "buyer private key")
	cmd.Flags().String(flagCoins, "", "coins 1QOS")

	return cmd
}

func queryBuyadCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query articleHash",
		Short: "query ad",
		RunE: func(cmd *cobra.Command, args []string) error {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
				}
			}()

			if len(args) < 1 {
				return errors.New("need articleHash")
			}

			articleHash := args[0]

			result, err := jianqian.QueryArticleBuyer(cdc, config.GetCLIContext().QSCCliContext, articleHash)
			if err != nil {
				return err
			}

			log.Printf("%+v", result)

			return nil
		},
	}
	return cmd
}
