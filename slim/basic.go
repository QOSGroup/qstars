package slim

import (
	"encoding/hex"
	"github.com/QOSGroup/qstars/client/lcd/lib"
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
	//from the QOS generation methods
	//	cdc := star.MakeCodec()
	//	acc := InitKeys(cdc)[2]
	entropy, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	seedo := bip39.NewSeed(mnemonic, "qstars")
	//seedh := hex.EncodeToString(seedo)

	key := ed25519.GenPrivKeyFromSecret(seedo)
	pub := key.PubKey().Bytes()
	addr := key.PubKey().Address()
	bech32Pub, _ := bech32.ConvertAndEncode(Bech32PrefixAccPub, pub)
	bech32Addr, _ := bech32.ConvertAndEncode(PREF_ADD, addr.Bytes())
	privkeyhex := "0x" + hex.EncodeToString(key.Bytes())

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

//convert the output to json string format
func AccountCreateStr() string {
	acc := AccountCreate()
	result, _ := lib.ResponseWrapper(cmCdc, acc, nil)
	out := string(result)
	//fmt.Println(out)
	return out
}
