package utils

import (
	"github.com/QOSGroup/qstars/client/context"
				"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qstars/wire"
)

// SendTx implements a auxiliary handler that facilitates sending a series of
// messages in a signed transaction given a TxContext and a QueryContext. It
// ensures that the account exists, has a proper number and sequence set. In
// addition, it builds and signs a transaction with the supplied messages.
// Finally, it broadcasts the signed transaction to a node.
func SendTx( cliCtx context.CLIContext, cdc *wire.Codec, txStd *txs.TxStd ,priv ed25519.PrivKeyEd25519) (string,error) {
	//if err := cliCtx.EnsureAccountExists(); err != nil {
	//	return err
	//}

	//from, err := cliCtx.GetFromAddress()
	//if err != nil {
	//	return err
	//}

	//// TODO: (ref #1903) Allow for user supplied account number without
	//// automatically doing a manual lookup.
	//if txCtx.AccountNumber == 0 {
	//	accNum, err := cliCtx.GetAccountNumber(from)
	//	if err != nil {
	//		return err
	//	}
	//
	//	txCtx = txCtx.WithAccountNumber(accNum)
	//}
	//
	//// TODO: (ref #1903) Allow for user supplied account sequence without
	//// automatically doing a manual lookup.
	//if txCtx.Sequence == 0 {
	//	accSeq, err := cliCtx.GetAccountSequence(from)
	//	if err != nil {
	//		return err
	//	}
	//
	//	txCtx = txCtx.WithSequence(accSeq)
	//}

	//passphrase, err := keys.GetPassphrase(cliCtx.FromAddressName)
	//if err != nil {
	//	return err
	//}

	// build and sign the transaction
	//txBytes, err := txCtx.BuildAndSign(cliCtx.FromAddressName, "", msgs,priv)
	//if err != nil {
	//	return "",err
	//}



	txBytes, err := cdc.MarshalBinaryBare(txStd)
	if err != nil {
		panic("use cdc encode object fail")
	}

	// broadcast to a Tendermint node
	resJSON,err := cliCtx.EnsureBroadcastTx(txBytes)

	if err != nil {
		return "",err
	}
	return resJSON.Hash.String(),err
}
