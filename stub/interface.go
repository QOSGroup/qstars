package stub

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/QOSGroup/qstars/star"
	"github.com/bartekn/go-bip39"
	"io/ioutil"
	"log"
	"net/http"

	qbasetypes "github.com/QOSGroup/qbase/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/bech32"
)

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

type ResultCreateAccount struct {
	PubKey   string `json:"pubKey"`
	PrivKey  string `json:"privKey"`
	Addr     string `json:"addr"`
	Mnemonic string `json:"mnemonic"`
	Type     string `json:"type"`
}

const (
	// Bech32 prefixes
	Bech32PrefixAccPub = "cosmosaccpub"
	AccountResultType  = "local"
)

func AccountCreate() *ResultCreateAccount {
	cdc := star.MakeCodec()
	acc := InitKeys(cdc)[2]
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
	bech32Addr, _ := bech32.ConvertAndEncode(qbasetypes.PREF_ADD, addr.Bytes())
	//	privkeybase64 := base64.StdEncoding.EncodeToString(key[:])
	//change privkey output to hex string format
	//privkeyhex := hex.EncodeToString(key.Bytes())
	privkeyhex := hex.EncodeToString(acc.PrivKey.Bytes())

	//Type field for future use
	Type := AccountResultType

	result := &ResultCreateAccount{}
	result.PubKey = bech32Pub
	result.PrivKey = privkeyhex
	result.Addr = bech32Addr
	result.Mnemonic = mnemonic
	result.Type = Type

	return result
}

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
	bech32Addr, _ := bech32.ConvertAndEncode(qbasetypes.PREF_ADD, addr.Bytes())
	//	privkeybase64 := base64.StdEncoding.EncodeToString(key[:])
	//change privkey output to hex string format
	privkeyhex := hex.EncodeToString(key[:])

	Type := AccountResultType
	result := &ResultCreateAccount{}
	result.PubKey = bech32Pub
	result.PrivKey = privkeyhex
	result.Addr = bech32Addr
	result.Mnemonic = mncode
	result.Type = Type

	output, _ := json.Marshal(result)
	//	fmt.Println(string(output))
	return string(output)

}

type PubAddrRetrieval struct {
	PubKey string `json:"pubKey"`
	Addr   string `json:"addr"`
}

func PubAddrRetrievalStr(s string) string {
	//the privkey output was in hex string format, decode it with the same decoding
	bz, _ := hex.DecodeString(s)
	var key ed25519.PrivKeyEd25519
	copy(key[:], bz)
	pub := key.PubKey().Bytes()
	addr := key.PubKey().Address()
	bech32Pub, _ := bech32.ConvertAndEncode(Bech32PrefixAccPub, pub)
	bech32Addr, _ := bech32.ConvertAndEncode(qbasetypes.PREF_ADD, addr.Bytes())

	result := &PubAddrRetrieval{}
	result.PubKey = bech32Pub
	result.Addr = bech32Addr
	output, _ := json.Marshal(result)
	//	fmt.Println(string(output))
	return string(output)
}

//func genStdSendTx(cdc *amino.Codec, sender qbasetypes.Address, receiver qbasetypes.Address, coin qbasetypes.BaseCoin,
//	priKey ed25519.PrivKeyEd25519, nonce int64, chainid string) *txs.TxStd {
//	sendTx := bank.NewSendTx(sender, receiver, coin)
//	gas := qbasetypes.NewInt(int64(0))
//	tx := txs.NewTxStd(&sendTx, chainid, gas)
//	//priHex, _ := hex.DecodeString(senderPriHex[2:])
//	//var priKey ed25519.PrivKeyEd25519
//	//cdc.MustUnmarshalBinaryBare(priHex, &priKey)
//	signature, _ := tx.SignTx(priKey, nonce)
//	tx.Signature = []txs.Signature{txs.Signature{
//		Pubkey:    priKey.PubKey(),
//		Signature: signature,
//		Nonce:     nonce,
//	}}
//	return tx
//}
//

////only need the following arguments, it`s enough!

//func QSCtransferSend(addrto, amount, privkey, chain string) string {
//	//make codec
//	cdc := star.MakeCodec()
//	bank.RegisterWire(cdc)
//	//generate the receiver address, i.e. "addrto" with the following format
//	to, err := types.AccAddressFromBech32(addrto)
//	if err != nil {
//		return "nil"
//	}
//
//	//generate the sender address, i.e. the "from" part as the input with privkey in hex string format

//	_, addrben32, priv := utility.PubAddrRetrieval(privkey, cdc)
//	from, err := types.AccAddressFromBech32(addrben32)
//	//	key := account.AddressStoreKey(from)
//
//	//coins for transfer
//	accountinfo := QSCQueryAccountGet(addrben32)
//	var cc qbasetypes.BaseCoin
//	var qcoins types.Coins
//	for _, qsc := range account.Coins {
//		amount := qsc.Amount
//		qcoins = append(qcoins, types.NewCoin(qsc.Name, types.NewInt(amount.Int64())))
//	}
//
//
//	//chainID with context setting
///*	var chainID string
//	directTOQOS := config.GetCLIContext().Config.DirectTOQOS
//	var cliCtx context.CLIContext
//	if directTOQOS==true{
//		cliCtx = *config.GetCLIContext().QOSCliContext
//		chainID = config.GetCLIContext().Config.QOSChainID
//	}else {
//		cliCtx = *config.GetCLIContext().QSCCliContext
//		chainID = config.GetCLIContext().Config.ChainID
//	}
//
//
//	account, err := config.GetCLIContext().QOSCliContext.GetAccount(key,cdc)
//	if err != nil {
//		return nil, err
//	}
//
//	var cc qbasetypes.BaseCoin
//	var qcoins types.Coins
//	for _, qsc := range account.Coins {
//		amount := qsc.Amount
//		qcoins = append(qcoins, types.NewCoin(qsc.Name, types.NewInt(amount.Int64())))
//
//		//TODO-------------------------
//		if !amount.IsZero() {
//			mount := qbasetypes.NewInt(100)
//			cc = qbasetypes.BaseCoin{
//				Name:   qsc.Name,
//				Amount: mount,
//			}
//		}
//	}
//
//	var coins types.Coins
//	if !qcoins.IsGTE(coins) {
//		return nil, errors.Errorf("Address %s doesn't have enough coins to pay for this transaction.", from)
//	}
//
//	var nn int64
//	if directTOQOS==true {
//		nn = int64(account.Nonce)
//	}else {
//		qscaccount, err := config.GetCLIContext().QSCCliContext.GetAccount(key,cdc)
//		if err != nil{
//			if err.Error()=="Account is not exsit." {
//				nn = int64(0)
//			}else{
//				return nil,err
//			}
//		}else{
//			nn = int64(qscaccount.Nonce)
//		}
//	}
//*/
//
//
//	//http part for restruction
///*
//	payload := map[string]interface{}{"amount": amount, "privatekey": privkey, "chain_id": chain, "account_number": accountnumber, "sequence": seq, "gas": gas}
//	jsonpayload, _ := json.Marshal(payload)
//	data := bytes.NewBuffer(jsonpayload)
//*/
//	msg := genStdSendTx(cdc,from,to,cc,priv,nn,chain)

//	//response, err := utils.SendTx(cliCtx,cdc,msg,priv)
//	//result := &bank.SendResult{}
//	//result.Hash = response
//
//	jasonpayload,_ := json.Marshal(msg)
//	data := bytes.NewBuffer(jasonpayload)
//	aurl := Accounturl + addrto + "/send"
//	req, _ := http.NewRequest("POST", aurl, data)
//	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//	clt := http.Client{}
//	resp, _ := clt.Do(req)
//	defer resp.Body.Close()
//	bresp, err := ioutil.ReadAll(resp.Body)
//	var body []byte
//	n := len(bresp)
//	for i := 0; i < n; i++ {
//		body = append(body, bresp[i])
//	}
//	if err != nil {
//		fmt.Println(err)
//		log.Fatal(err)
//	}
//	defer resp.Body.Close()
//	output := string(body)
//	//	fmt.Println(output)
//	return output
//}

//
///*
//func AccountCreatePostGet(ul,input string) string {
//	payload := map[string]interface{}{"name":input}
//	jsonpayload, _ := json.Marshal(payload)
//	body := bytes.NewBuffer(jsonpayload)
//	//	fmt.Println(body)
//	req, _ := http.NewRequest("POST", ul, body)
//	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//	clt := http.Client{}
//	resp, _ := clt.Do(req)
//	defer resp.Body.Close()
//	//	fmt.Println(resp.StatusCode)
//	if resp.StatusCode == http.StatusOK {
//		bresp, err := ioutil.ReadAll(resp.Body)
//		var body []byte
//		n := len(bresp)
//		for i := 0; i < n; i++ {
//			body = append(body, bresp[i])
//		}
//		if err != nil {
//			fmt.Println(err)
//			log.Fatal(err)
//		}
//		output := string(body)
//		return output
//	}
//	return "nil"
//}
//*/
