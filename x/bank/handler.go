package bank

import (
	"bytes"
	"fmt"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qbase/context"
)

type SendTx struct {
	From btypes.Address  `json:"from"`
	To   btypes.Address  `json:"to"`
	Coin btypes.BaseCoin `json:"coin"`
}

var _ txs.ITx = (*SendTx)(nil)

func NewSendTx(from btypes.Address, to btypes.Address, coin btypes.BaseCoin) SendTx {
	return SendTx{From: from, To: to, Coin: coin}
}

func (tx SendTx) ValidateData(ctx context.Context) bool {
	if len(tx.From) == 0 || len(tx.To) == 0 || btypes.NewInt(0).GT(tx.Coin.Amount) {
		return false
	}
	return true
}


func (tx SendTx) Exec(ctx context.Context) (result btypes.Result, crossTxQcps *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.ABCICodeOK,
	}

	fmt.Println("--------------------------------------------")
	fmt.Println("--------------------------------------------")
	fmt.Println("--------------------------------------------")
	cross := txs.TxQcp{

	}
	crossTxQcps = &cross
	tt := txs.TxStd{

	}

	tt.ITx = tx
	crossTxQcps.TxStd = &tt
	crossTxQcps.To = "basecoin-chain"

	r := btypes.Result{
		Code:0,
	}
	return r, &cross
}

func (tx SendTx) GetSigner() []btypes.Address {
	return []btypes.Address{tx.From}
}

func (tx SendTx) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx SendTx) GetGasPayer() btypes.Address {
	return tx.From
}

func (tx SendTx) GetSignData() []byte {
	var buf bytes.Buffer
	buf.Write(tx.From)
	buf.Write(tx.To)
	buf.Write([]byte(tx.Coin.String()))
	return buf.Bytes()
}
