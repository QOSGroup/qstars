package main

import (
	"encoding/json"
	"fmt"
	"github.com/QOSGroup/qbase/types"

	"github.com/QOSGroup/qbase/server"
	"github.com/QOSGroup/qstars/baseapp"
	"github.com/QOSGroup/qstars/slim"
	"github.com/QOSGroup/qstars/star"
	"github.com/QOSGroup/qstars/wire"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/cli"
	qosdinit "github.com/QOSGroup/qos/cmd/qosd/init"
	go_amino "github.com/tendermint/go-amino"
	tmtypes "github.com/tendermint/tendermint/types"
	sdk "github.com/QOSGroup/qstars/types"
	"os"
)

func main() {

	//logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "main")

	//db, err := dbm.NewGoLevelDB("qstarsd", filepath.Join(rootDir, "data"))
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}

	cdc := star.MakeCodec()
	baseapp.InitApp()
	ctx := baseapp.GetServerContext().ServerContext
	//viper.SetDefault("pruning", "nothing")

	rootCmd := &cobra.Command{
		Use:               "start",
		Short:             "qstars Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}
	rootDir := os.ExpandEnv("$HOME/.qstarsd")
	rootCmd.AddCommand(server.InitCmd(ctx, cdc, genBaseCoindGenesisDoc, rootDir))

	rootCmd.AddCommand(qosdinit.AddGenesisAccount(cdc))
	rootCmd.AddCommand(qosdinit.AddGenesisValidator(cdc))

	server.AddCommands(ctx, cdc, rootCmd, star.NewApp)

	// prepare and add flags

	executor := cli.PrepareBaseCmd(rootCmd, "QSC", rootDir)

	err := executor.Execute()
	if err != nil {
		// Note: Handle with #870
		panic(err)
	}

	return
}

func genBaseCoindGenesisDoc(ctx *server.Context, cdc *go_amino.Codec, chainID string, nodeValidatorPubKey crypto.PubKey) (tmtypes.GenesisDoc, error) {

	validator := tmtypes.GenesisValidator{
		PubKey: nodeValidatorPubKey,
		Power:  10,
	}

	//addr, _, err := types.GenerateCoinKey(cdc, types.DefaultCLIHome)
	//if err != nil {
	//	return tmtypes.GenesisDoc{}, err
	//}

	acc := slim.AccountCreate("")

	output, err := wire.MarshalJSONIndent(cdc, acc)
	if err != nil {
		return tmtypes.GenesisDoc{},err
	}

	fmt.Println(string(output))
	addr,_:=sdk.AccAddressFromBech32(acc.Addr)

	appState, err := BaseCoinQOSAppGenState(cdc, addr)
	if err != nil {
		return tmtypes.GenesisDoc{}, err
	}

	return tmtypes.GenesisDoc{
		ChainID:    chainID,
		Validators: []tmtypes.GenesisValidator{validator},
		AppState:   appState,
	}, nil

}

func BaseCoinQOSAppGenState(cdc *go_amino.Codec, addr types.Address) (appState json.RawMessage, err error) {

	appState = json.RawMessage(fmt.Sprintf(`{
		"qcps":[{
			"name": "qos",
			"chain_id": "qos",
			"pub_key":{
        		"type": "tendermint/PubKeyEd25519",
        		"value": "ish2+qpPsoHxf7m+uwi8FOAWw6iMaDZgLKl1la4yMAs="
			}
		}],
  		"accounts": [{
    		"address": "%s",
    		"coins": [
      			{
        			"coin_name":"qstar",
        			"amount":"100000000"
      			}
			]
  		}]
	}`, addr))
	return
}

