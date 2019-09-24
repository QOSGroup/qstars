package supply

import (
	"fmt"
	"github.com/QOSGroup/qstars/config"
	sdk "github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/wire"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flag_orderfrom   = "orderfrom"
	flag_orderto     = "orderto"
	flag_orderid     = "orderid"
	flag_orderamount = "orderamount"
)

// SendTxCmd will create a send tx and sign it with the given key.
func NewCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "NewOrder",
		Short: "send order tx",
		RunE: func(cmd *cobra.Command, args []string) error {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
				}
			}()
			orderfrom := viper.GetString(flag_orderfrom)
			orderto := viper.GetString(flag_orderto)
			orderid := viper.GetString(flag_orderid)
			orderamount := viper.GetString(flag_orderamount)
			to, _ := sdk.AccAddressFromBech32(orderto)
			coins, _ := sdk.ParseCoins(orderamount)
			result, err := Send(cdc, orderfrom, to, coins, orderid, NewSendOptions(
				gas(viper.GetInt64("gas")),
				fee(viper.GetString("fee"))))
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println(result)
			}
			return nil
		},
	}

	cmd.Flags().String(flag_orderfrom, "", "NewArticle author address")
	cmd.Flags().String(flag_orderto, "", "NewArticle original address")
	cmd.Flags().String(flag_orderid, "", "NewArticle article hash")
	cmd.Flags().String(flag_orderamount, "", "NewArticle share author ")

	return cmd
}

func QueryArticleCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "QueryOrder",
		Short: "query  order and send tx",
		RunE: func(cmd *cobra.Command, args []string) error {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
				}
			}()
			orderid := viper.GetString(flag_orderid)
			result, err := config.GetCLIContext().QSCCliContext.QueryStore([]byte(orderid), OrderMapperName)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println(string(result))
			return nil
		},
	}

	cmd.Flags().String(flag_orderid, "", "query order")

	return cmd
}
