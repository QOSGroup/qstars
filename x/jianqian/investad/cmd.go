package investad

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
	flagInvester    = "investor"
	flagCoins       = "coins"
	flagArticleHash = "articleHash"
	flagChainid     = "chainid"
)

// InvestadCmd will create a send tx and sign it with the given key.
func InvestadCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "investad",
		Short: "invest ad",
		RunE: func(cmd *cobra.Command, args []string) error {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
				}
			}()

			articleHash := viper.GetString(flagArticleHash)
			chainid := viper.GetString(flagChainid)
			investor := viper.GetString(flagInvester) //Teddy changes
			coins := viper.GetString(flagCoins)

			_, addrben32, _ := utility.PubAddrRetrievalFromAmino(investor, cdc)
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

			tx := InvestAd(cdc, chainid, articleHash, coins, investor, nonce)
			result := InvestAdBackground(cdc, tx)

			log.Printf(result)

			return nil
		},
	}

	return cmd
}
