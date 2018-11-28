package tx

import (
	"fmt"
	"github.com/QOSGroup/qstars/config"
	"github.com/spf13/viper"
	"path/filepath"
)

func GetConfig()*Config{
	if Conf==nil{
		InitConfig()
	}
	return Conf
}

type Config struct {
	Community string    `mapstructure:"community"`
	Authormock string   `mapstructure:"authormock"`
	Adbuyermock string  `mapstructure:"adbuyermock"`
	Banker string       `mapstructure:"banker"`
	Dappowner string    `mapstructure:"dappowner"`
}


var (
	Conf     *Config
)


func InitConfig() (err error) {
	rootDir := 	config.GetCLIContext().Config.RootDir
	configFilePath := filepath.Join(rootDir, "config/config.toml")

	Conf = NewConfig()
	viper.SetConfigName("config") // 配置文件的名字
	viper.SetConfigType("toml") // 配置文件的类型
	viper.AddConfigPath(configFilePath) // 配置文件的路径
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(&Conf); err != nil {
		panic(fmt.Errorf("unable to decode into struct：  %s \n", err))
	}

	fmt.Println(Conf.Banker)
	return nil
}

func NewConfig() *Config {
	return &Config{}
}
