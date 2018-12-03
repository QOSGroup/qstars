package investad

import (
	"encoding/json"
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
	flagInvestor    = "investor"
	flagCoins       = "coins"
	flagArticleHash = "articleHash"
	flagChainid     = "chainid"
)

// InvestadCmd will create a send tx and sign it with the given key.
func InvestadCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "investad",
		Short: "invest ad",
	}

	cmd.AddCommand(
		investadCmd(cdc),
		queryInvestadCmd(cdc),
	)

	return cmd
}

func investadCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invest",
		Short: "invest ad",
		RunE: func(cmd *cobra.Command, args []string) error {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
				}
			}()

			articleHash := viper.GetString(flagArticleHash)
			chainid := viper.GetString(flagChainid)
			investor := viper.GetString(flagInvestor) //Teddy changes
			coins := viper.GetString(flagCoins)

			//_, addrben32, _ := utility.PubAddrRetrievalFromAmino(investor, cdc)
			//from, err := types.AccAddressFromBech32(addrben32)
			//key := account.AddressStoreKey(from)
			//if err != nil {
			//	return err
			//}
			//acc, err := config.GetCLIContext().QOSCliContext.GetAccount(key, cdc)
			//if err != nil {
			//	return err
			//}
			//nonce := int64(acc.Nonce)
			//nonce++
			nonce := int64(0)

			tx := InvestAd(cdc, chainid, articleHash, coins, investor, nonce)
			fmt.Printf("InvestAd:%s\n", tx)
			var ri ResultInvest
			if err := json.Unmarshal([]byte(tx), &ri); err != nil {
				return fmt.Errorf("Unmarshal tx error:%s ", err.Error())
			}

			if ri.Code != "0" {
				return fmt.Errorf("InvestAd tx error:%s ", ri.Reason)
			}

			result := InvestAdBackground(cdc, string(ri.Result))

			log.Printf(result)

			return nil
		},
	}

	cmd.Flags().String(flagArticleHash, "", "articleHash")
	cmd.Flags().String(flagChainid, "", "Chainid")
	cmd.Flags().String(flagInvestor, "", "investor private key")
	cmd.Flags().String(flagCoins, "", "coins 1AOE")

	return cmd
}

func queryInvestadCmd(cdc *wire.Codec) *cobra.Command {
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

			result, err := jianqian.ListInvestors(config.GetCLIContext().QSCCliContext, cdc, articleHash)
			if err != nil {
				return err
			}

			log.Printf("%+v", result)

			return nil
		},
	}
	return cmd
}
