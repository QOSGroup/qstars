package buyad

import (
	"fmt"
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/utility"
	"github.com/QOSGroup/qstars/wire"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"

	types "github.com/QOSGroup/qstars/types"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
				}
			}()

			articleHash := viper.GetString(flagArticleHash)
			chainid := viper.GetString(flagChainid)
			buyer := viper.GetString(flagBuyer) //Teddy changes
			coins := viper.GetString(flagCoins)

			_, addrben32, _ := utility.PubAddrRetrievalFromAmino(buyer, cdc)
			from, err := types.AccAddressFromBech32(addrben32)
			key := account.AddressStoreKey(from)
			if err != nil {
				return err
			}
			acc, err := config.GetCLIContext().QOSCliContext.GetAccount(key, cdc)
			if err != nil {
				return err
			}
			nonce := int64(acc.Nonce)
			nonce++

			tx := BuyAd(cdc, chainid, articleHash, coins, buyer, nonce)
			result := BuyAdBackground(cdc, tx)

			log.Printf(result)

			return nil
		},
	}

	return cmd
}
