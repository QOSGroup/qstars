package slim

import (
	"encoding/base64"
	"github.com/QOSGroup/qstars/slim/funcInlocal/respwrap"
	"github.com/bartekn/go-bip39"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/bech32"
)

const PREF_ADD = "address"

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
	entropy, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	seedo := bip39.NewSeed(mnemonic, "qstars")
	//seedh := hex.EncodeToString(seedo)

	key := ed25519.GenPrivKeyFromSecret(seedo)
	pub := key.PubKey().Bytes()
	addr := key.PubKey().Address()
	bech32Pub, _ := bech32.ConvertAndEncode(Bech32PrefixAccPub, pub)
	bech32Addr, _ := bech32.ConvertAndEncode(PREF_ADD, addr.Bytes())

	privkeybase64 := base64.StdEncoding.EncodeToString(key.Bytes())
	//privkeyhex := "0x" + hex.EncodeToString(key.Bytes())

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
	result, _ := respwrap.ResponseWrapper(Cdc, acc, nil)
	out := string(result)
	//fmt.Println(out)
	return out
}

func AccountRecoverStr(mncode string) string {
	seed := bip39.NewSeed(mncode, "qstars")
	key := ed25519.GenPrivKeyFromSecret(seed)
	pub := key.PubKey().Bytes()
	addr := key.PubKey().Address()
	bech32Pub, _ := bech32.ConvertAndEncode("cosmosaccpub", pub)
	bech32Addr, _ := bech32.ConvertAndEncode(PREF_ADD, addr.Bytes())
	privkeybase64 := base64.StdEncoding.EncodeToString(key.Bytes())
	//change privkey output to hex string format
	//privkeyhex := "0x" + hex.EncodeToString(key.Bytes())

	Type := AccountResultType
	result := &ResultCreateAccount{}
	result.PubKey = bech32Pub
	result.PrivKey = privkeybase64
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
	//the privkey output was in hex string format, decode it with the same decoding
	//bz, _ := hex.DecodeString(s[2:])
	bz, _ := base64.StdEncoding.DecodeString(s)
	var key ed25519.PrivKeyEd25519
	Cdc.MustUnmarshalBinaryBare(bz, &key)
	pub := key.PubKey().Bytes()
	addr := key.PubKey().Address()
	bech32Pub, _ := bech32.ConvertAndEncode(Bech32PrefixAccPub, pub)
	bech32Addr, _ := bech32.ConvertAndEncode(PREF_ADD, addr.Bytes())
	//privkeybase64 := base64.StdEncoding.EncodeToString(key.Bytes())

	result := &PubAddrRetrieval{}
	result.PubKey = bech32Pub
	result.Addr = bech32Addr
	//result.PrivKey = privkeybase64

	resp, _ := respwrap.ResponseWrapper(Cdc, result, nil)
	out := string(resp)
	return out
}
