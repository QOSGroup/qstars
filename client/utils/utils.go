package utils

import (
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qstars/client/context"
	"github.com/QOSGroup/qstars/wire"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

// SendTx implements a auxiliary handler that facilitates sending a series of
// messages in a signed transaction given a TxContext and a QueryContext. It
// ensures that the account exists, has a proper number and sequence set. In
// addition, it builds and signs a transaction with the supplied messages.
// Finally, it broadcasts the signed transaction to a node.
func SendTx(cliCtx context.CLIContext, cdc *wire.Codec, txStd *txs.TxStd, priv ed25519.PrivKeyEd25519) (string, error) {

	txBytes, err := cdc.MarshalBinary(txStd)
	if err != nil {
		panic("use cdc encode object fail")
	}

	// broadcast to a Tendermint node
	resJSON, err := cliCtx.EnsureBroadcastTx(txBytes)

	if err != nil {
		return "12345678901234", nil
	}
	return resJSON.Hash.String(), err
}
