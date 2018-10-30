package utils

import (
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qstars/client/context"
	"github.com/QOSGroup/qstars/wire"
)

// SendTx implements a auxiliary handler that facilitates sending a series of
// messages in a signed transaction given a TxContext and a QueryContext. It
// ensures that the account exists, has a proper number and sequence set. In
// addition, it builds and signs a transaction with the supplied messages.
// Finally, it broadcasts the signed transaction to a node.
func SendTx(cliCtx context.CLIContext, cdc *wire.Codec, txStd *txs.TxStd) (string, error) {

	txBytes, err := cdc.MarshalBinaryBare(txStd)
	if err != nil {
		panic("use cdc encode object fail")
	}

	// broadcast to a Tendermint node
	resJSON, err := cliCtx.EnsureBroadcastTx(txBytes)

	if err != nil {
		return err.Error(), err
	}
	return resJSON.Hash.String(), err
}
