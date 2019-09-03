package recharge

import (
	"fmt"
	//"github.com/QOSGroup/qbase/account"
	qosaccount "github.com/QOSGroup/qos/types"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/utility"
	"github.com/QOSGroup/qstars/wire"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"log"
)

const (
	flagAmount     = "amount"
	flagPrivatekey = "privatekey"
	flagCointype   = "cointype"
	flagDeposit    = "deposit"
)

// SendTxCmd will create a send tx and sign it with the given key.
func RechargeCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "Recharge",
		Short: "Recharge or extract",
		RunE: func(cmd *cobra.Command, args []string) error {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
				}
			}()

			amount := viper.GetString(flagAmount)
			privatekey := viper.GetString(flagPrivatekey) //Teddy changes
			cointype := viper.GetString(flagCointype)
			deposit := viper.GetString(flagDeposit)

			_, addrben32, _ := utility.PubAddrRetrievalFromAmino(privatekey, cdc)
			from, _ := types.AccAddressFromBech32(addrben32)
			//key := account.AddressStoreKey(from)
			//if err != nil {
			//	return err
			//}
			//var qscnonce int64 = 0
			//qscacc, err := getQSCAcc(key, cdc)
			//if err != nil {
			//	qscnonce = 0
			//} else {
			//	qscnonce = int64(qscacc.Nonce)
			//}
			result := Recharge(cdc, amount, privatekey,from.String(), cointype, deposit)

			//fmt.Printf("Recharge:%s\n", tx)
			//var ri common.Result
			//if err := json.Unmarshal([]byte(tx), &ri); err != nil {
			//	return fmt.Errorf("Unmarshal tx error:%s ", err.Error())
			//}
			//
			//if ri.Code != "0" {
			//	return fmt.Errorf("InvestAd tx error:%s ", ri.Reason)
			//}
			//
			//result := RechargeBackground(cdc, string(ri.Result), 0)

			log.Printf(result)

			return nil
		},
	}

	cmd.Flags().String(flagAmount, "", "amount")
	cmd.Flags().String(flagPrivatekey, "", "privatekey")
	cmd.Flags().String(flagCointype, "", "cointype")
	cmd.Flags().String(flagDeposit, "", "deposit")

	return cmd
}

func getQSCAcc(address []byte, cdc *wire.Codec) (*qosaccount.QOSAccount, error) {
	return config.GetCLIContext().QSCCliContext.GetAccount(address, cdc)
}
