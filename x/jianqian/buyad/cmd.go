package buyad

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/QOSGroup/qbase/account"
	qosaccount "github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/utility"
	"github.com/QOSGroup/qstars/wire"
	"github.com/QOSGroup/qstars/x/jianqian"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"time"

	"github.com/QOSGroup/qstars/types"
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

func getQOSAcc(address []byte, cdc *wire.Codec) (*qosaccount.QOSAccount, error) {
	return config.GetCLIContext().QOSCliContext.GetAccount(address, cdc)
}

func getQSCAcc(address []byte, cdc *wire.Codec) (*qosaccount.QOSAccount, error) {
	return config.GetCLIContext().QSCCliContext.GetAccount(address, cdc)
}

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
			log.Printf("BuyAd tx:%+v", tx)

			var rb ResultBuy
			if err := json.Unmarshal([]byte(tx), &rb); err != nil {
				return fmt.Errorf("Unmarshal tx error:%s ", err.Error())
			}

			if rb.Code != "0" {
				return fmt.Errorf("InvestAd tx error:%s ", rb.Reason)
			}

			result := BuyAdBackground(cdc, string(rb.Result), time.Second*60)
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
