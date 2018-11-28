package config

import (
	"github.com/QOSGroup/qstars/client/context"
	"github.com/QOSGroup/qstars/wire"
	"github.com/spf13/viper"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

var Clictx QStarsClientContext

type QStarsClientContext struct {
	QOSCliContext *context.CLIContext
	QSCCliContext *context.CLIContext
	Config        *CLIConfig
}

type CLIConfig struct {
	QOSChainID          string `mapstructure:"qos_chain_id"`
	QSCChainID          string `mapstructure:"qsc_chain_id"`
	RootDir             string `mapstructure:"home"`
	QOSNodeURI          string `mapstructure:"qos_node_uri"`
	QSTARSNodeURI       string `mapstructure:"qstars_node_uri"`
	DirectTOQOS         bool   `mapstructure:"direct_to_qos"`
	WaitingForQosResult string `mapstructure:"waiting_for_qos_result"`
}

func GetCLIContext() QStarsClientContext {
	return Clictx
}

func CreateCLIContextTwo(cdc *wire.Codec, cfg *CLIConfig) QStarsClientContext {

	pQOSCliContext := context.NewCLIContext1(cfg.QOSNodeURI).
		WithCodec(cdc).
		WithLogger(os.Stdout)

	pQSCCliContext := context.NewCLIContext1(cfg.QSTARSNodeURI).
		WithCodec(cdc).
		WithLogger(os.Stdout)

	Clictx = QStarsClientContext{
		QOSCliContext: &pQOSCliContext,
		QSCCliContext: &pQSCCliContext,
		Config:        cfg,
	}
	return Clictx
}

// If a new config is created, change some of the default tendermint settings
func InterceptLoadConfig() (conf *CLIConfig, err error) {

	tmpConf := DefaultConfig()
	err = viper.Unmarshal(tmpConf)
	if err != nil {
		panic(err)
	}
	rootDir := tmpConf.RootDir

	configFilePath := filepath.Join(rootDir, "config/config.toml")
	// Intercept only if the file doesn't already exist

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) { //
		// the following parse config is needed to create directories
		conf, _ = ParseConfig()
		conf.QSCChainID = "qstars-test"
		conf.QOSNodeURI = "qos-test"
		conf.QOSNodeURI = "localhost:26657"
		conf.QSTARSNodeURI = "localhost:26657"
		WriteConfigFile(configFilePath, conf)
		// Fall through, just so that its parsed into memory.
	}

	if conf == nil {
		conf, err = ParseConfig()
	}

	return
}

func DefaultConfig() *CLIConfig {
	var result *CLIConfig
	result = &CLIConfig{} //RootDir:"~/.qstarscli"
	if len(result.RootDir) == 0 {
		usr, err := user.Current()
		if nil == err {
			if "windows" == runtime.GOOS {
				result.RootDir = usr.HomeDir + "\\.qstarscli"
			} else {
				result.RootDir = usr.HomeDir + "/.qstarscli"
			}
		}

	}
	return result
}

// ParseConfig retrieves the default environment configuration,
// sets up the Tendermint root and ensures that the root exists
func ParseConfig() (*CLIConfig, error) {
	conf := DefaultConfig()
	err := viper.Unmarshal(conf)
	if err != nil {
		return nil, err
	}
	//conf.SetRoot(conf.RootDir)
	EnsureRoot(conf.RootDir)
	return conf, err
}
