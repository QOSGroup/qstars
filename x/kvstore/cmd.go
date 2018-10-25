package kvstore

import (
	"github.com/QOSGroup/qstars/wire"

	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagKey        = "key"
	flagValue      = "value"
	flagPrivateKey = "private"
	chainIdFlag    = "chain-id"
	sequenceFlag   = "sequence"
)

// SendTxCmd will create a send tx and sign it with the given key.
func SendKVCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kvset",
		Short: "Create and sign a send set kv tx",
		RunE: func(cmd *cobra.Command, args []string) error {

			privatekey := viper.GetString(flagPrivateKey)
			key := viper.GetString(flagKey)
			value := viper.GetString(flagValue)

			opts, err := NewSendKVOption(
				SendKVOptionChainID(viper.GetString(chainIdFlag)),
				SendKVOptionSequence(viper.GetString(sequenceFlag)),
			)
			if err != nil {
				return err
			}

			result, err := SendKV(cdc, privatekey, key, value, opts)
			if err != nil {
				return err
			}

			output, err := wire.MarshalJSONIndent(cdc, result)
			if err != nil {
				return err
			}

			fmt.Println(string(output))

			return nil
		},
	}

	cmd.Flags().String(flagKey, "", "Key")
	cmd.Flags().String(flagValue, "", "Value")
	cmd.Flags().String(flagPrivateKey, "", "Private key")

	return cmd
}

// SendTxCmd will create a send tx and sign it with the given key.
func GetKVCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kvget",
		Short: "Create and sign a send set kv tx",
		RunE: func(cmd *cobra.Command, args []string) error {

			key := viper.GetString(flagKey)

			result, err := GetKV( cdc, key, nil)
			if err != nil {
				return err
			}

			output, err := wire.MarshalJSONIndent(cdc, result)
			if err != nil {
				return err
			}

			fmt.Println(string(output))

			return err
		},
	}

	cmd.Flags().String(flagKey, "", "Key")
	cmd.Flags().String(flagValue, "", "Value")
	cmd.Flags().String(flagPrivateKey, "", "Private key")

	return cmd
}
