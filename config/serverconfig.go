package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
)

/*
qstars configuration file
QStarsPrivateKey qstars privatekey
QStarsTransactions TBD
*/
type ServerConf struct {
	QStarsPrivateKey   string `toml:"QStarsPrivateKey"`
	QStarsTransactions string `toml:"QStarsTransactions"`
	QOSChainName       string `toml:"QOSChainName"`
	RootDir            string `toml:"RootDir"`
	Community          string `mapstructure:"community"`
}

var serverconfiguration *ServerConf = nil

func Init(filename string, rootDir string) (p *ServerConf, err error) {
	if serverconfiguration == nil {
		sconf, err := readConf(filename)
		if len(sconf.RootDir) == 0 {
			sconf.RootDir = rootDir
		}
		serverconfiguration = sconf
		return sconf, err
	} else {
		return serverconfiguration, nil
	}
}

func GetServerConf() *ServerConf {
	return serverconfiguration
}
func readConf(fname string) (p *ServerConf, err error) {
	var (
		fp       *os.File
		fcontent []byte
	)
	p = new(ServerConf) // &Person{}
	if fp, err = os.Open(fname); err != nil {
		fmt.Println("open error ", err)
		return
	}

	if fcontent, err = ioutil.ReadAll(fp); err != nil {
		fmt.Println("ReadAll error ", err)
		return
	}

	if err = toml.Unmarshal(fcontent, p); err != nil {
		fmt.Println("toml.Unmarshal error ", err)
		return
	}
	return
}
