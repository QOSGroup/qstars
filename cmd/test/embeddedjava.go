package main

import (
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/star"
	"github.com/QOSGroup/qstars/wire"
	"github.com/QOSGroup/qstars/x/bank"

	sdk "github.com/QOSGroup/qstars/types"
)
var CDC *wire.Codec
var CONF *config.CLIConfig

func InitJNI(){
	CDC := star.MakeCodec()
	//if cmd.Name() == version.VersionCmd.Name() {
	//}
	CONF, err := config.InterceptLoadConfig()
	if err != nil {
		panic("config is wrong.")
	}
	config.CreateCLIContextTwo(CDC, CONF)
}

func SendByJNI(fromStr string, toStr1 string,coinstr string) string{

	// (*SendResult, error)

	toStr, err := sdk.AccAddressFromBech32(toStr1)
	if err != nil {
		return ""
	}


	// parse coins trying to be sent
	coins, err := sdk.ParseCoins(coinstr)
	if err != nil {
		return ""
	}
	result ,err := bank.Send(CDC, fromStr, toStr, coins , nil)
	output, err := wire.MarshalJSONIndent(CDC,result)
	if err != nil {
		return err.Error()
	}
	return string(output)
}