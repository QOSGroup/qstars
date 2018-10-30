package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"github.com/BurntSushi/toml"
)

// person
type ServerConf struct {
	QStarsPrivateKey   string `toml:"QStarsPrivateKey"`
}


func ReadConf(fname string) (p *ServerConf, err error) {
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
