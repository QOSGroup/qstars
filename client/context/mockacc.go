package context

import (
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func keyPubAddr() (crypto.PrivKey, crypto.PubKey, btypes.Address) {
	key := ed25519.GenPrivKey()
	pub := key.PubKey()
	addr := btypes.Address(pub.Address())
	return key, pub, addr
}

//func genNewAccount() (qosAccount qosacc.QOSAccount) {
//	_, pub, addr := keyPubAddr()
//	coinList := []*types.QSC{
//		types.NewQSC("QSC1", btypes.NewInt(1234)),
//		types.NewQSC("QSC2", btypes.NewInt(5678)),
//	}
//	qosAccount = qosacc.QOSAccount{
//		account.BaseAccount{addr, pub, 0},
//		btypes.NewInt(5380394853),
//		coinList,
//	}
//	return
//}
