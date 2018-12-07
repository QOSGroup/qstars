package slim

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
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
	Shost         string
	Mhost         string
	QSCAccounturl string
	QOSAccounturl string
	Accounturl    string
	KVurl         string
	QResulturl    string
)

//set Block Chain entrance hosts for both Qstars and Qmoon
func SetBlockchainEntrance(qstarshost, qmoonhost string) {
	Shost = qstarshost
	Mhost = qmoonhost
	QSCAccounturl = "http://" + Shost + "/QSCaccounts/"
	QOSAccounturl = "http://" + Shost + "/QOSaccounts/"
	Accounturl = "http://" + Shost + "/accounts/"
	KVurl = "http://" + Shost + "/kv/"
	QResulturl = "http://" + Shost + "/commits/"

}

func init() {
	var sh string
	var mh string
	SetBlockchainEntrance(sh, mh)
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
	kvurl := KVurl + k
	resp, _ := http.Get(kvurl)
	//	fmt.Println(KVurl)
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			//log.Fatal(err)
		}
		defer resp.Body.Close()
		output := string(body)
		return output
	}
	return "nil"
}

func QSCQueryAccountGet(addr string) string {
	aurl := QSCAccounturl + addr
	resp, _ := http.Get(aurl)
	var body []byte
	var err error
	if resp.StatusCode == http.StatusOK {
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
	}

	defer resp.Body.Close()
	output := string(body)
	return output
}

//for QOS account query function
func QOSQueryAccountGet(addr string) string {
	aurl := QOSAccounturl + addr
	resp, _ := http.Get(aurl)
	var body []byte
	var err error
	if resp.StatusCode == http.StatusOK {
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
	}

	defer resp.Body.Close()
	output := string(body)
	return output
}

//for QOSCommitResultCheck Restful interface
func QOSCommitResultCheck(txhash, height string) string {
	qstarskey := "heigth:" + height + ",hash:" + txhash
	qrcurl := QResulturl + qstarskey
	resp, _ := http.Get(qrcurl)
	var body []byte
	var err error
	if resp.StatusCode == http.StatusOK {
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
	}

	defer resp.Body.Close()
	output := fmt.Sprintf("This function has not been realized in QOS yet:%v", string(body))
	return output

}
