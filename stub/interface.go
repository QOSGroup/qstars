package stub

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/QOSGroup/qbase/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/bech32"
)

// url for different modules, e.g kv, accounts, .etc
var (
	HostIP     = "http://localhost:1317"
	Accounturl = HostIP + "/accounts/"
	KVurl      = HostIP + "/kv"
)

type ResultCreateAccount struct {
	PubKey  string `json:"pubKey"`
	PrivKey string `json:"privKey"`
	Addr    string `json:"addr"`
	Seed    string `json:"seed"`
}

func AccountCreate() *ResultCreateAccount {
	const (
		// Bech32 prefixes
		Bech32PrefixAccPub = "cosmosaccpub"
	)
	fmt.Printf("++++++++++++++++\n")
	key := ed25519.GenPrivKey()
	pub := key.PubKey().Bytes()
	addr := key.PubKey().Address()
	bech32Pub, _ := bech32.ConvertAndEncode(Bech32PrefixAccPub, pub)
	bech32Addr, _ := bech32.ConvertAndEncode(types.PREF_ADD, addr.Bytes())
	privkeybase64 := base64.StdEncoding.EncodeToString(key[:])

	result := &ResultCreateAccount{}
	result.PubKey = bech32Pub
	result.PrivKey = privkeybase64
	result.Addr = bech32Addr
	result.Seed = ""

	return result
}

func AccountCreateStr() string {
	acc := AccountCreate()
	output := acc.PrivKey + "#" + acc.PrivKey + "#" + acc.Addr
	fmt.Println(output)
	return output
}

func QSCQueryAccountGet(ul string) string {
	aurl := Accounturl + ul
	resp, _ := http.Get(aurl)
	if resp.StatusCode == http.StatusOK {
		bresp, err := ioutil.ReadAll(resp.Body)
		var body []byte
		n := len(bresp)
		for i := 0; i < n; i++ {
			body = append(body, bresp[i])
		}
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

func QSCKVStoreSetPost(k, v, privkey, chain string) (result int) {
	payload := map[string]interface{}{"key": k, "value": v, "privatekey": privkey, "chainid": chain}
	jsonpayload, _ := json.Marshal(payload)
	body := bytes.NewBuffer(jsonpayload)
	req, _ := http.NewRequest("POST", KVurl, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	clt := http.Client{}
	resp, _ := clt.Do(req)
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		result = 1
		return result
	}

	return 0

}

func QSCKVStoreGetQuery(ul string) string {
	kvurl := KVurl + "/" + ul
	resp, _ := http.Get(kvurl)
	if resp.StatusCode == http.StatusOK {
		bresp, err := ioutil.ReadAll(resp.Body)
		var body []byte
		n := len(bresp)
		for i := 0; i < n; i++ {
			body = append(body, bresp[i])
		}
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

func QSCtransferPost(ul, a, privkey, chain, ac, seq, g string) string {
	aurl := Accounturl + ul + "/send"
	payload := map[string]interface{}{"amount": a, "privatekey": privkey, "chain_id": chain, "account_number": ac, "sequence": seq, "gas": g}
	jsonpayload, _ := json.Marshal(payload)
	data := bytes.NewBuffer(jsonpayload)
	req, _ := http.NewRequest("POST", aurl, data)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	clt := http.Client{}
	resp, _ := clt.Do(req)
	defer resp.Body.Close()
	bresp, err := ioutil.ReadAll(resp.Body)
	var body []byte
	n := len(bresp)
	for i := 0; i < n; i++ {
		body = append(body, bresp[i])
	}
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	defer resp.Body.Close()
	output := string(body)
	//	fmt.Println(output)
	return output

}

/*
func AccountCreatePostGet(ul,input string) string {
	payload := map[string]interface{}{"name":input}
	jsonpayload, _ := json.Marshal(payload)
	body := bytes.NewBuffer(jsonpayload)
	//	fmt.Println(body)
	req, _ := http.NewRequest("POST", ul, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	clt := http.Client{}
	resp, _ := clt.Do(req)
	defer resp.Body.Close()
	//	fmt.Println(resp.StatusCode)
	if resp.StatusCode == http.StatusOK {
		bresp, err := ioutil.ReadAll(resp.Body)
		var body []byte
		n := len(bresp)
		for i := 0; i < n; i++ {
			body = append(body, bresp[i])
		}
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}
		output := string(body)
		return output
	}
	return "nil"
}
*/
