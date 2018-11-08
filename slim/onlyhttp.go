package slim

import (
	"bytes"
	"fmt"
	"github.com/QOSGroup/qstars/star"
	"github.com/QOSGroup/qstars/wire"
	"io/ioutil"
	"log"
	"net/http"
)

var cmCdc *wire.Codec

func init() {
	cmCdc = star.MakeCodec()
}

type sendKVReq struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	PrivateKey string `json:"privatekey"`
	ChainID    string `json:"chainid"`
}

var KVurl = "http://localhost:8080"

func QSCKVStoreSetPost(k, v, privkey, chain string) (result string) {
	skr := sendKVReq{}
	skr.Key = k
	skr.Value = v
	skr.PrivateKey = privkey
	skr.ChainID = chain
	payload, _ := cmCdc.MarshalJSON(skr)
	body := bytes.NewBuffer(payload)
	req, _ := http.NewRequest("POST", KVurl, body)
	req.Header.Set("Content-Type", "application/json")
	clt := http.Client{}
	resp, _ := clt.Do(req)
	defer resp.Body.Close()
	rep, _ := ioutil.ReadAll(resp.Body)
	output := string(rep)
	fmt.Println(output)
	return output
}

func QSCKVStoreGetQuery(k string) string {
	kvurl := KVurl + "/" + k
	resp, _ := http.Get(kvurl)
	//	fmt.Println(KVurl)
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}
		defer resp.Body.Close()
		output := string(body)
		return output
	}
	return "nil"
}
