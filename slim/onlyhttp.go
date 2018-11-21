package slim

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type sendKVReq struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	PrivateKey string `json:"privatekey"`
	ChainID    string `json:"chainid"`
}

// IP initialization
var (
	HostIP     string
	Accounturl string
	KVurl      string
)

func GetIPfrom(host string) {
	HostIP = host
	Accounturl = "http://" + HostIP + "/accounts/"
	KVurl = "http://" + HostIP + "/kv"
}

func init() {
	var h string
	GetIPfrom(h)
}

func QSCKVStoreSetPost(k, v, privkey, chain string) (result string) {
	skr := sendKVReq{}
	skr.Key = k
	skr.Value = v
	skr.PrivateKey = privkey
	skr.ChainID = chain
	payload, _ := Cdc.MarshalJSON(skr)
	body := bytes.NewBuffer(payload)
	req, _ := http.NewRequest("POST", KVurl, body)
	req.Header.Set("Content-Type", "application/json")
	clt := http.Client{}
	resp, _ := clt.Do(req)
	defer resp.Body.Close()
	rep, _ := ioutil.ReadAll(resp.Body)
	output := string(rep)
	//fmt.Println(output)
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
			//log.Fatal(err)
		}
		defer resp.Body.Close()
		output := string(body)
		return output
	}
	return "nil"
}

func QSCQueryAccountGet(addr string) string {
	aurl := Accounturl + addr
	resp, _ := http.Get(aurl)
	var body []byte
	var err error
	if resp.StatusCode == http.StatusOK {
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
	}

	defer resp.Body.Close()
	output := string(body)
	return output
}
