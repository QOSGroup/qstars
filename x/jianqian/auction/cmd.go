package auction

import (
	"encoding/json"
	"fmt"
	"github.com/QOSGroup/qbase/account"
	qosaccount "github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/utility"
	"github.com/QOSGroup/qstars/types"

	"github.com/QOSGroup/qstars/wire"
	"github.com/QOSGroup/qstars/x/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"time"
)

const (
	flag_privatekey   = "privatekey"
	flag_otherAddress = "otheraddress"
	flag_articleHash  = "articleHash"
	flag_cointype     = "coinType"
	flag_amount       = "amount"
)

func getQOSAcc(address []byte, cdc *wire.Codec) (*qosaccount.QOSAccount, error) {
	return config.GetCLIContext().QOSCliContext.GetAccount(address, cdc)
}

func getQSCAcc(address []byte, cdc *wire.Codec) (*qosaccount.QOSAccount, error) {
	return config.GetCLIContext().QSCCliContext.GetAccount(address, cdc)
}

// SendTxCmd will create a send tx and sign it with the given key.
func NewAuctionCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "NewAuction",
		Short: "add Auction and send tx",
		RunE: func(cmd *cobra.Command, args []string) error {
			articleHash := viper.GetString(flag_articleHash)
			private := viper.GetString(flag_privatekey)
			otherAddres := viper.GetString(flag_otherAddress) //Teddy changes
			cointype := viper.GetString(flag_cointype)
			amount := viper.GetString(flag_amount)

			fmt.Println(articleHash,private,otherAddres,cointype,amount)
			_, addrben32, _ := utility.PubAddrRetrievalFromAmino(private, cdc)
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
			tx := AcutionAd(cdc, articleHash, private, otherAddres, cointype, amount, qosnonce, qscnonce)
			log.Printf("NewAuction tx:%+v", tx)

			var rb common.Result
			if err := json.Unmarshal([]byte(tx), &rb); err != nil {
				return fmt.Errorf("Unmarshal tx error:%s ", err.Error())
			}

			if rb.Code != "0" {
				return fmt.Errorf("InvestAd tx error:%s ", rb.Reason)
			}
			result := AcutionAdBackground(cdc, string(rb.Result), time.Second*60)
			log.Printf(result)

			return nil
		},
	}

	cmd.Flags().String(flag_privatekey, "", "NewAuction private key")
	cmd.Flags().String(flag_otherAddress, "", "NewAuction other address")
	cmd.Flags().String(flag_articleHash, "", "NewAuction article hash")
	cmd.Flags().String(flag_cointype, "", "NewAuction coin type")
	cmd.Flags().String(flag_amount, "", "NewAuction  amount")


	return cmd
}

func QueryMaxAcutionCMD(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "QueryMaxAcution",
		Short: "query  max acution",
		RunE: func(cmd *cobra.Command, args []string) error {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
				}
			}()

			articleHash := viper.GetString(flag_articleHash)

			result:= QueryMaxAcution(cdc,articleHash)

			fmt.Println(result)

			return nil
		},
	}


	cmd.Flags().String(flag_articleHash, "", "query article hash")

	return cmd
}
func QueryAllAcutionCMD(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "QueryAllAcution",
		Short: "query  all acution",
		RunE: func(cmd *cobra.Command, args []string) error {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
				}
			}()

			articleHash := viper.GetString(flag_articleHash)

			result:= QueryAllAcution(cdc,articleHash)

			fmt.Println(result)

			return nil
		},
	}


	cmd.Flags().String(flag_articleHash, "", "query article hash")

	return cmd
}