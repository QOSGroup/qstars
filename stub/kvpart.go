package stub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func QSCKVStoreSetPost(k, v, privkey, chain string) (result int) {
	payload := map[string]interface{}{"key": k, "value": v, "privatekey": privkey, "chainid": chain}
	jsonpayload, _ := json.Marshal(payload)
	body := bytes.NewBuffer(jsonpayload)
	req, _ := http.NewRequest("POST", KVurl, body)
	req.Header.Set("Content-Type", "application/json")
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
