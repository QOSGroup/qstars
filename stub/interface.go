package stub

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/bartekn/go-bip39"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/QOSGroup/qbase/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/bech32"
)

// IP initialization
var (
	HostIP string
	Accounturl string
	KVurl string
)
func GetIPfrom(host string) {
	HostIP = host
	Accounturl = "http://" + HostIP + "/accounts/"
	KVurl      = "http://" + HostIP + "/kv"
}

func init() {
	var h string
	GetIPfrom(h)
}

type ResultCreateAccount struct {
	PubKey   string `json:"pubKey"`
	PrivKey  string `json:"privKey"`
	Addr     string `json:"addr"`
	Mnemonic string `json:"mnemonic"`
	Type string `json:"type"`
}

const (
	// Bech32 prefixes
	Bech32PrefixAccPub = "cosmosaccpub"
	AccountResultType = "local"
)

func AccountCreate() *ResultCreateAccount {
//	fmt.Printf("Please write down your mnemonic words BELOW for further account recovery purpose:\n")
	entropy, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	//fmt.Println(mnemonic)
	seedo := bip39.NewSeed(mnemonic, "qstars")
	//seedh := hex.EncodeToString(seedo)

	key := ed25519.GenPrivKeyFromSecret(seedo)
	pub := key.PubKey().Bytes()
	addr := key.PubKey().Address()
	bech32Pub, _ := bech32.ConvertAndEncode(Bech32PrefixAccPub, pub)
	bech32Addr, _ := bech32.ConvertAndEncode(types.PREF_ADD, addr.Bytes())
	privkeybase64 := base64.StdEncoding.EncodeToString(key[:])

	//Type field for future use
	Type := AccountResultType

	result := &ResultCreateAccount{}
	result.PubKey = bech32Pub
	result.PrivKey = privkeybase64
	result.Addr = bech32Addr
	result.Mnemonic = mnemonic
	result.Type = Type

	return result
}
//convert the output to json string format
func AccountCreateStr() string {
	acc := AccountCreate()
	output, _ := json.Marshal(acc)
	out := string(output)
	fmt.Println(out)
	return out
}

func QSCQueryAccountGet(addr string) string {
	aurl := Accounturl + addr
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

func QSCKVStoreGetQuery(k string) string {
	kvurl := KVurl + "/" + k
	resp, _ := http.Get(kvurl)
//	fmt.Println(KVurl)
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

//input spelling adjusted to human readable
func QSCtransferPost(addr, amount, privkey, chain, accountnumber, seq, gas string) string {
	aurl := Accounturl + addr + "/send"
	payload := map[string]interface{}{"amount": amount, "privatekey": privkey, "chain_id": chain, "account_number": accountnumber, "sequence": seq, "gas": gas}
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

func AccountRecoverStr(mncode string) string {
	seed := bip39.NewSeed(mncode, "qstars")
	key := ed25519.GenPrivKeyFromSecret(seed)
	pub := key.PubKey().Bytes()
	addr := key.PubKey().Address()
	bech32Pub, _ := bech32.ConvertAndEncode("cosmosaccpub", pub)
	bech32Addr, _ := bech32.ConvertAndEncode(types.PREF_ADD, addr.Bytes())
	privkeybase64 := base64.StdEncoding.EncodeToString(key[:])

	Type := AccountResultType
	result := &ResultCreateAccount{}
	result.PubKey = bech32Pub
	result.PrivKey = privkeybase64
	result.Addr = bech32Addr
	result.Mnemonic = mncode
	result.Type = Type

	output, _ := json.Marshal(result)
//	fmt.Println(string(output))
	return string(output)

}

type PubAddrRetrieval struct {
	PubKey   string `json:"pubKey"`
	Addr     string `json:"addr"`
}

func PubAddrRetrievalStr(s string) string {
	bz,_ :=base64.StdEncoding.DecodeString(s)
	var key ed25519.PrivKeyEd25519
	copy(key[:], bz)
	pub := key.PubKey().Bytes()
	addr := key.PubKey().Address()
	bech32Pub, _ := bech32.ConvertAndEncode(Bech32PrefixAccPub, pub)
	bech32Addr, _ := bech32.ConvertAndEncode(types.PREF_ADD, addr.Bytes())

	result := &PubAddrRetrieval{}
	result.PubKey = bech32Pub
	result.Addr = bech32Addr
	output, _ := json.Marshal(result)
//	fmt.Println(string(output))
	return string(output)
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
