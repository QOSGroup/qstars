package stub

import (
	"encoding/base64"
	"fmt"
	qbasetypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/client/lcd/lib"
	"github.com/bartekn/go-bip39"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/bech32"
)

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
	bech32Addr, _ := bech32.ConvertAndEncode(qbasetypes.PREF_ADD, addr.Bytes())
	//according to #92, return base64 format output
	privkeybase64 := base64.StdEncoding.EncodeToString(key.Bytes())

	//to be discarded, change privkey output to hex string format, from the QOS mechanism
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
	result, _ := lib.ResponseWrapper(cmCdc, acc, nil)
	out := string(result)
	fmt.Println(out)
	return out
}
func AccountRecoverStr(mncode string) string {
	seed := bip39.NewSeed(mncode, "qstars")
	key := ed25519.GenPrivKeyFromSecret(seed)
	pub := key.PubKey().Bytes()
	addr := key.PubKey().Address()
	bech32Pub, _ := bech32.ConvertAndEncode("cosmosaccpub", pub)
	bech32Addr, _ := bech32.ConvertAndEncode(qbasetypes.PREF_ADD, addr.Bytes())
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

	resp, _ := lib.ResponseWrapper(cmCdc, result, nil)
	out := string(resp)
	fmt.Println(out)
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
	cmCdc.MustUnmarshalBinaryBare(bz, &key)
	pub := key.PubKey().Bytes()
	addr := key.PubKey().Address()
	bech32Pub, _ := bech32.ConvertAndEncode(Bech32PrefixAccPub, pub)
	bech32Addr, _ := bech32.ConvertAndEncode(qbasetypes.PREF_ADD, addr.Bytes())

	result := &PubAddrRetrieval{}
	result.PubKey = bech32Pub
	result.Addr = bech32Addr

	resp, _ := lib.ResponseWrapper(cmCdc, result, nil)
	out := string(resp)
	//	fmt.Println(out)
	return out
}
