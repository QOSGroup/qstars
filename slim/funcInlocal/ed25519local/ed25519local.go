package ed25519local

import (
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"github.com/tendermint/ed25519/edwards25519"
	"github.com/tendermint/go-amino"
)

const (
	Ed25519PrivKeyAminoRoute = "tendermint/PrivKeyEd25519"
	Ed25519PubKeyAminoRoute  = "tendermint/PubKeyEd25519"
	// Size of an Edwards25519 signature. Namely the size of a compressed
	// Edwards25519 point, and a field element. Both of which are 32 bytes.
	SignatureSize  = 64
	PrivateKeySize = 64
	PublicKeySize  = 32
)

var cdc = amino.NewCodec()

func init() {
	cdc.RegisterInterface((*PubKey)(nil), nil)
	cdc.RegisterConcrete(PubKeyEd25519{},
		Ed25519PubKeyAminoRoute, nil)

	cdc.RegisterInterface((*PrivKey)(nil), nil)
	cdc.RegisterConcrete(PrivKeyEd25519{},
		Ed25519PrivKeyAminoRoute, nil)
}

var _ PrivKey = PrivKeyEd25519{}

type PrivKey interface {
	Bytes() []byte
	Sign(msg []byte) ([]byte, error)
	PubKey() PubKey
}

type Address = HexBytes

type HexBytes []byte

func (bz HexBytes) Bytes() []byte {
	return bz
}

type PubKey interface {
	Address() Address
	Bytes() []byte
	VerifyBytes(msg []byte, sig []byte) bool
}

// PrivKeyEd25519 implements crypto.PrivKey.
type PrivKeyEd25519 [64]byte

// Bytes marshals the privkey using amino encoding.
func (privKey PrivKeyEd25519) Bytes() []byte {
	return cdc.MustMarshalBinaryBare(privKey)
}

// Sign produces a signature on the provided message.
func (privKey PrivKeyEd25519) Sign(msg []byte) ([]byte, error) {
	privKeyBytes := [64]byte(privKey)
	signatureBytes := Sign(&privKeyBytes, msg)
	return signatureBytes[:], nil
}

// PubKey gets the corresponding public key from the private key.
func (privKey PrivKeyEd25519) PubKey() PubKey {
	privKeyBytes := [64]byte(privKey)
	initialized := false
	// If the latter 32 bytes of the privkey are all zero, compute the pubkey
	// otherwise privkey is initialized and we can use the cached value inside
	// of the private key.
	for _, v := range privKeyBytes[32:] {
		if v != 0 {
			initialized = true
			break
		}
	}
	if initialized {
		var pubkeyBytes [32]byte
		copy(pubkeyBytes[:], privKeyBytes[32:])
		return PubKeyEd25519(pubkeyBytes)
	}

	pubBytes := *MakePublicKey(&privKeyBytes)
	return PubKeyEd25519(pubBytes)
}

// PubKeyEd25519 implements crypto.PubKey for the Ed25519 signature scheme.
type PubKeyEd25519 [32]byte

// Address is the SHA256-20 of the raw pubkey bytes.
func (pubKey PubKeyEd25519) Address() Address {
	return Address(Sum(pubKey[:]))
}

// Bytes marshals the PubKey using amino encoding.
func (pubKey PubKeyEd25519) Bytes() []byte {
	bz, err := cdc.MarshalBinaryBare(pubKey)
	if err != nil {
		panic(err)
	}
	return bz
}

func (pubKey PubKeyEd25519) VerifyBytes(msg []byte, sig_ []byte) bool {
	// make sure we use the same algorithm to sign
	if len(sig_) != SignatureSize {
		return false
	}
	sig := new([SignatureSize]byte)
	copy(sig[:], sig_)
	pubKeyBytes := [32]byte(pubKey)
	return Verify(&pubKeyBytes, msg, sig)
}

// GenPrivKeyFromSecret hashes the secret with SHA2, and uses
// that 32 byte output to create the private key.
// NOTE: secret should be the output of a KDF like bcrypt,
// if it's derived from user input.
func GenPrivKeyFromSecret(secret []byte) PrivKeyEd25519 {
	privKey32 := Sha256(secret) // Not Ripemd160 because we want 32 bytes.
	privKey := new([64]byte)
	copy(privKey[:32], privKey32)
	// ed25519.MakePublicKey(privKey) alters the last 32 bytes of privKey.
	// It places the pubkey in the last 32 bytes of privKey, and returns the
	// public key.
	MakePublicKey(privKey)
	return PrivKeyEd25519(*privKey)
}

// MakePublicKey makes a publicKey from the first half of privateKey.
func MakePublicKey(privateKey *[64]byte) (publicKey *[32]byte) {
	publicKey = new([32]byte)

	h := sha512.New()
	h.Write(privateKey[:32])
	digest := h.Sum(nil)

	digest[0] &= 248
	digest[31] &= 127
	digest[31] |= 64

	var A edwards25519.ExtendedGroupElement
	var hBytes [32]byte
	copy(hBytes[:], digest)
	edwards25519.GeScalarMultBase(&A, &hBytes)
	A.ToBytes(publicKey)

	copy(privateKey[32:], publicKey[:])
	return
}

func Sha256(bytes []byte) []byte {
	hasher := sha256.New()
	hasher.Write(bytes)
	return hasher.Sum(nil)
}

// Sum returns the first 20 bytes of SHA256 of the bz.
func Sum(bz []byte) []byte {
	hash := sha256.Sum256(bz)
	return hash[:20]
}

// Sign signs the message with privateKey and returns a signature.
func Sign(privateKey *[PrivateKeySize]byte, message []byte) *[SignatureSize]byte {
	h := sha512.New()
	h.Write(privateKey[:32])

	var digest1, messageDigest, hramDigest [64]byte
	var expandedSecretKey [32]byte
	h.Sum(digest1[:0])
	copy(expandedSecretKey[:], digest1[:])
	expandedSecretKey[0] &= 248
	expandedSecretKey[31] &= 63
	expandedSecretKey[31] |= 64

	h.Reset()
	h.Write(digest1[32:])
	h.Write(message)
	h.Sum(messageDigest[:0])

	var messageDigestReduced [32]byte
	edwards25519.ScReduce(&messageDigestReduced, &messageDigest)
	var R edwards25519.ExtendedGroupElement
	edwards25519.GeScalarMultBase(&R, &messageDigestReduced)

	var encodedR [32]byte
	R.ToBytes(&encodedR)

	h.Reset()
	h.Write(encodedR[:])
	h.Write(privateKey[32:])
	h.Write(message)
	h.Sum(hramDigest[:0])
	var hramDigestReduced [32]byte
	edwards25519.ScReduce(&hramDigestReduced, &hramDigest)

	var s [32]byte
	edwards25519.ScMulAdd(&s, &hramDigestReduced, &expandedSecretKey, &messageDigestReduced)

	signature := new([64]byte)
	copy(signature[:], encodedR[:])
	copy(signature[32:], s[:])
	return signature
}

// Verify returns true iff sig is a valid signature of message by publicKey.
func Verify(publicKey *[PublicKeySize]byte, message []byte, sig *[SignatureSize]byte) bool {
	if sig[63]&224 != 0 {
		return false
	}

	var A edwards25519.ExtendedGroupElement
	if !A.FromBytes(publicKey) {
		return false
	}
	edwards25519.FeNeg(&A.X, &A.X)
	edwards25519.FeNeg(&A.T, &A.T)

	h := sha512.New()
	h.Write(sig[:32])
	h.Write(publicKey[:])
	h.Write(message)
	var digest [64]byte
	h.Sum(digest[:0])

	var hReduced [32]byte
	edwards25519.ScReduce(&hReduced, &digest)

	var R edwards25519.ProjectiveGroupElement
	var b [32]byte
	copy(b[:], sig[32:])
	edwards25519.GeDoubleScalarMultVartime(&R, &hReduced, &A, &b)

	var checkR [32]byte
	R.ToBytes(&checkR)
	return subtle.ConstantTimeCompare(sig[:32], checkR[:]) == 1
}
