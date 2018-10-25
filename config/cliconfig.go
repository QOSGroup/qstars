package config

import (
	"github.com/QOSGroup/qstars/client/context"
	"github.com/QOSGroup/qstars/wire"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var Clictx QStarsClientContext

type QStarsClientContext struct{
	QOSCliContext *context.CLIContext
	QSCCliContext *context.CLIContext
	Config *CLIConfig
}

type CLIConfig struct {
	ChainID string `mapstructure:"chain_id"`
	RootDir string `mapstructure:"home"`
	QOSNodeURI string `mapstructure:"qos_node_uri"`
	QSTARSNodeURI string `mapstructure:"qstars_node_uri"`
}

func GetCLIContext() QStarsClientContext {
	return Clictx
}

func CreateCLIContextTwo(cdc *wire.Codec, cfg *CLIConfig) QStarsClientContext{

	pQOSCliContext := context.NewCLIContext1(cfg.QOSNodeURI).
		WithCodec(cdc).
		WithLogger(os.Stdout)

	pQSCCliContext := context.NewCLIContext1(cfg.QSTARSNodeURI).
		WithCodec(cdc).
		WithLogger(os.Stdout)

	Clictx = QStarsClientContext{
		QOSCliContext: &pQOSCliContext,
		QSCCliContext: &pQSCCliContext,
		Config: cfg,
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

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// the following parse config is needed to create directories
		conf, _ = ParseConfig()
		conf.ChainID = "abc"
		conf.QOSNodeURI = "localhost:1317"
		conf.QSTARSNodeURI = "localhost:1317"
		WriteConfigFile(configFilePath, conf)
		// Fall through, just so that its parsed into memory.
	}

	if conf == nil {
		conf, err = ParseConfig()
	}

	return
}

func DefaultConfig() *CLIConfig{
	return &CLIConfig{
	}
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

