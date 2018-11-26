package jsdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var KVurl string

type dispatchCoinsReq struct {
	Addresses    string `json:"address"`
	Coins        string `json:"coins"`
	Causecodes   string `json:"causecodes"`
	Causestrings string `json:"causestrings"`
	Gas          string `json:"gas"`
	ChainID      string `json:"chainid"`
}

func DispatchCoins(addrs, coins, causecodes, causestrings, gas string) (result string) {
	skr := dispatchCoinsReq{}
	skr.Addresses = addrs
	skr.Coins = coins
	skr.Causecodes = causecodes
	skr.Causestrings = causestrings
	skr.Gas = gas
	skr.ChainID = "test-chain-HfiBIx"
	payload, _ := json.Marshal(skr)
	fmt.Println(string(payload))
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

type newArticleReq struct {
	AuthorAddress     string `json:"authorAddress"`
	OriginAuthor      string `json:"originAuthor"`
	ArticleHash       string `json:"articleHash"`
	ShareAuthor       string `json:"shareAuthor"`
	ShareOriginAuthor string `json:"shareOriginAuthor"`
	ShareCommunity    string `json:"shareCommunity"`
	ShareInvestor     string `json:"shareInvestor"`
	EndInvestDate     string `json:"endInvestDate"`
	EndBuyDate        string `json:"endBuyDate"`
}

type SendNA struct {
	Key     string        `json:"key"`
	Value   newArticleReq `json:"value"`
	ChainID string        `json:"chainid"`
}

func NewArticle(authorAddress, originAuthor, articleHash, shareAuthor, shareOriginAuthor, shareCommunity, shareInvestor, endInvestDate, endBuyDate string) string {
	skr := newArticleReq{}
	skr.AuthorAddress = authorAddress
	skr.OriginAuthor = originAuthor
	skr.ArticleHash = articleHash
	skr.ShareAuthor = shareAuthor
	skr.ShareOriginAuthor = shareOriginAuthor
	skr.ShareCommunity = shareCommunity
	skr.ShareInvestor = shareInvestor
	skr.EndInvestDate = endInvestDate
	skr.EndBuyDate = endBuyDate

	send := SendNA{}
	send.ChainID = "test-chain-HfiBIx"
	send.Key = articleHash
	send.Value = skr
	payload, _ := json.Marshal(send)
	fmt.Println(string(payload))
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
