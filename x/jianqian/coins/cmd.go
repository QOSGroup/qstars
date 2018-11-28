package coins

import (
	"fmt"
	"github.com/QOSGroup/qstars/wire"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	)

const (
	flag_address="address"
	flag_acoin="coin"
	flag_causecode="causecode"
	flag_causestrings="causestrings"

)

// SendTxCmd will create a send tx and sign it with the given key.
func DispatchAOECmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "DispatchAoe",
		Short: "Dispatch AOE and send tx",
		RunE: func(cmd *cobra.Command, args []string) error {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
				}
			}()
			toStr := viper.GetString("address")
			coin := viper.GetString("coin")
			causecode:=viper.GetString("causecode")
			causestrings:=viper.GetString("causestrings")

			result := DispatchAOE(cdc,toStr,coin,causecode,causestrings,"0")
			fmt.Println(result)
			return nil
		},
	}

	cmd.Flags().String(flag_address, "", "Address to send coins")
	cmd.Flags().String(flag_acoin, "", "Address to send coins")
	cmd.Flags().String(flag_causecode, "", "Address to send coins")
	cmd.Flags().String(flag_causestrings, "", "Address to send coins")
	return cmd
}