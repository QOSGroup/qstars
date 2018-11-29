package slim

import (
	"encoding/json"
	"fmt"
	"github.com/QOSGroup/qstars/slim/funcInlocal/bech32local"
	"github.com/QOSGroup/qstars/slim/funcInlocal/bip39local"
	"github.com/QOSGroup/qstars/slim/funcInlocal/ed25519local"
	"github.com/QOSGroup/qstars/slim/funcInlocal/respwrap"
	"log"
)

const PREF_ADD = "address"

type ResultCreateAccount struct {
	PubKey   string `json:"pubKey"`
	PrivKey  string `json:"privKey"`
	Addr     string `json:"addr"`
	Mnemonic string `json:"mnemonic"`
	Type     string `json:"type"`
}

type PrivkeyAmino struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

const (
	// Bech32 prefixes
	Bech32PrefixAccPub = "cosmosaccpub"
	AccountResultType  = "local"
)

func AccountCreate() *ResultCreateAccount {
	entropy, _ := bip39local.NewEntropy(256)
	mnemonic, _ := bip39local.NewMnemonic(entropy)
	seedo := bip39local.NewSeed(mnemonic, "qstars")

	key := ed25519local.GenPrivKeyFromSecret(seedo)
	pub := key.PubKey().Bytes()
	addr := key.PubKey().Address()
	bech32Pub, _ := bech32local.ConvertAndEncode(Bech32PrefixAccPub, pub)
	bech32Addr, _ := bech32local.ConvertAndEncode(PREF_ADD, addr.Bytes())

	privkeyAmino, _ := Cdc.MarshalJSON(key)
	var privkeyAminoStc PrivkeyAmino
	err := json.Unmarshal(privkeyAmino, &privkeyAminoStc)
	if err != nil {
		log.Fatalln(err.Error())
	}
	privkeyAminoStr := string(privkeyAminoStc.Value)

	Type := AccountResultType

	result := &ResultCreateAccount{}
	result.PubKey = bech32Pub
	result.PrivKey = privkeyAminoStr
	result.Addr = bech32Addr
	result.Mnemonic = mnemonic
	result.Type = Type

	return result
}

//convert the output to json string format
func AccountCreateStr() string {
	acc := AccountCreate()
	result, _ := respwrap.ResponseWrapper(Cdc, acc, nil)
	out := string(result)

	return out
}

func AccountRecoverStr(mncode string) string {
	seed := bip39local.NewSeed(mncode, "qstars")
	key := ed25519local.GenPrivKeyFromSecret(seed)
	pub := key.PubKey().Bytes()
	addr := key.PubKey().Address()
	bech32Pub, _ := bech32local.ConvertAndEncode("cosmosaccpub", pub)
	bech32Addr, _ := bech32local.ConvertAndEncode(PREF_ADD, addr.Bytes())

	privkeyAmino, _ := Cdc.MarshalJSON(key)
	var privkeyAminoStc PrivkeyAmino
	err := json.Unmarshal(privkeyAmino, &privkeyAminoStc)
	if err != nil {
		log.Fatalln(err.Error())
	}
	privkeyAminoStr := string(privkeyAminoStc.Value)

	Type := AccountResultType
	result := &ResultCreateAccount{}
	result.PubKey = bech32Pub
	result.PrivKey = privkeyAminoStr
	result.Addr = bech32Addr
	result.Mnemonic = mncode
	result.Type = Type

	resp, _ := respwrap.ResponseWrapper(Cdc, result, nil)
	out := string(resp)
	return out
}

type PubAddrRetrieval struct {
	PubKey string `json:"pubKey"`
	Addr   string `json:"addr"`
}

func PubAddrRetrievalStr(s string) string {
	//change the private unmarshal format according to the other pack
	ts := "{\"type\": \"tendermint/PrivKeyEd25519\",\"value\": \"" + s + "\"}"
	var key ed25519local.PrivKeyEd25519

	err := Cdc.UnmarshalJSON([]byte(ts), &key)
	if err != nil {
		fmt.Println(err)
	}
	pub := key.PubKey().Bytes()
	addr := key.PubKey().Address()
	bech32Pub, _ := bech32local.ConvertAndEncode(Bech32PrefixAccPub, pub)
	bech32Addr, _ := bech32local.ConvertAndEncode(PREF_ADD, addr.Bytes())

	result := &PubAddrRetrieval{}
	result.PubKey = bech32Pub
	result.Addr = bech32Addr

	resp, _ := respwrap.ResponseWrapper(Cdc, result, nil)
	out := string(resp)
	return out
}
