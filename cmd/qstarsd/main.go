package main

import (
	"github.com/QOSGroup/qstars/server"
	"github.com/spf13/cobra"
	"os"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/QOSGroup/qstars/star"
	)

func main() {

	//logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "main")

	//db, err := dbm.NewGoLevelDB("qstarsd", filepath.Join(rootDir, "data"))
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}

	cdc := star.MakeCodec()
	ctx := server.NewDefaultContext()
	//viper.SetDefault("pruning", "nothing")

	rootCmd := &cobra.Command{
		Use:               "start",
		Short:             "qstars Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	server.AddCommands(ctx, cdc, rootCmd, server.DefaultAppInit,
		server.ConstructAppCreator(star.NewApp, "qstars"))

	// prepare and add flags
	rootDir := os.ExpandEnv("$HOME/.qstarsd")
	executor := cli.PrepareBaseCmd(rootCmd, "QSC", rootDir)

	err := executor.Execute()
	if err != nil {
		// Note: Handle with #870
		panic(err)
	}

	return
}

