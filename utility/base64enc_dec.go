package utility

import "C"
import (
	"encoding/base64"
	"fmt"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/bech32"
)

var cdc = amino.NewCodec()

const (
	// expected address length
	//AddrLen = 20

	// Bech32 prefixes
	Bech32PrefixAccAddr = "cosmosaccaddr"
	Bech32PrefixAccPub  = "cosmosaccpub"
)

func Encbase64(input []byte) string {
	return base64.StdEncoding.EncodeToString(input[:])
}

func Decbase64(input string) []byte {
	bz, _ := base64.StdEncoding.DecodeString(input)
	return bz

}

func PubAddrRetrieval(s string) (string, string) {
	bz := Decbase64(s)
	var key ed25519.PrivKeyEd25519
	copy(key[:], bz)
	pub := key.PubKey().Bytes()
	addr := key.PubKey().Address()
	bech32Pub, _ := bech32.ConvertAndEncode(Bech32PrefixAccPub, pub)
	bech32Addr, _ := bech32.ConvertAndEncode(Bech32PrefixAccAddr, addr.Bytes())
	fmt.Println(bech32Pub)
	fmt.Println(bech32Addr)
	return bech32Pub, bech32Addr
}

//func main() {
//	s := "9Rg9mNEXVh9aUsxJ74Ogqe8O6wrBw8EeMhyK/GgHcfUsGprPgC7YXH6YEwGM+eXmc7oV1ci7ivlxo7k6amd3Lg=="
//	PubAddrRetrieval(s)
//	bz := Decbase64(s)
//	s1 := Encbase64(bz)
//	fmt.Printf("%x\n",bz)
//	fmt.Println(s1)
//}
