package kvstore

import (
	"os"

	"github.com/QOSGroup/qstars/client/context"
	"github.com/QOSGroup/qstars/client/utils"
	"github.com/QOSGroup/qstars/wire"

	"github.com/QOSGroup/qstars/utility"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"fmt"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"

)

const (
	flagKey        = "key"
	flagValue      = "value"
	flagPrivateKey = "private"
)

// SendTxCmd will create a send tx and sign it with the given key.
func SendKVCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kvset",
		Short: "Create and sign a send set kv tx",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout)

			privatekey := viper.GetString(flagPrivateKey)

			// parse coins trying to be sent
			key := viper.GetString(flagKey)

			value := viper.GetString(flagValue)

			//get addr from private key
			var priv ed25519.PrivKeyEd25519
			bz := utility.Decbase64(privatekey)
			copy(priv[:], bz)
			//_, addrben32 := utility.PubAddrRetrieval(privatekey)

			txStd := wrapToStdTx(key, value,"chainid")



			response,err := utils.SendTx( cliCtx,cdc,txStd,priv)
			fmt.Println(response)
			return err
		},
	}

	cmd.Flags().String(flagKey, "", "Key")
	cmd.Flags().String(flagValue, "", "Value")
	cmd.Flags().String(flagPrivateKey, "", "Private key")

	return cmd
}

func wrapToStdTx(key string, value string, chainid string) *txs.TxStd {
	kv := NewKvstoreTx([]byte(key), []byte(value))
	return txs.NewTxStd(kv, chainid, types.NewInt(int64(10000)))
}

// SendTxCmd will create a send tx and sign it with the given key.
func GetKVCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kvget",
		Short: "Create and sign a send set kv tx",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout)
			key := viper.GetString(flagKey)
			result, err := cliCtx.QueryKV([]byte(key))

			fmt.Println(string(result))
			return err
		},
	}

	cmd.Flags().String(flagKey, "", "Key")
	cmd.Flags().String(flagValue, "", "Value")
	cmd.Flags().String(flagPrivateKey, "", "Private key")

	return cmd
}
