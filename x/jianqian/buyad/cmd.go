package buyad

import (
	"fmt"
	"github.com/QOSGroup/qbase/account"
	qosaccount "github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/utility"
	"github.com/QOSGroup/qstars/wire"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"time"

	types "github.com/QOSGroup/qstars/types"
)

const (
	flagBuyer       = "buyer"
	flagCoins       = "coins"
	flagArticleHash = "articleHash"
	flagChainid     = "chainid"
)

func getQOSAcc(address []byte, cdc *wire.Codec) (*qosaccount.QOSAccount, error) {
	return config.GetCLIContext().QOSCliContext.GetAccount(address, cdc)
}

func getQSCAcc(address []byte, cdc *wire.Codec) (*qosaccount.QOSAccount, error) {
	return config.GetCLIContext().QSCCliContext.GetAccount(address, cdc)
}

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
			qosacc, err := getQOSAcc(key, cdc)
			if err != nil {
				return err
			}
			qosnonce := int64(qosacc.Nonce)

			qscacc, err := getQSCAcc(key, cdc)
			if err != nil {
				return err
			}
			qscnonce := int64(qscacc.Nonce)

			tx := BuyAd(cdc, chainid, articleHash, coins, buyer, qosnonce, qscnonce)
			result := BuyAdBackground(cdc, tx, time.Second*60)

			log.Printf(result)

			return nil
		},
	}

	return cmd
}
