package bank

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/tendermint/tendermint/types"
	"strconv"
)

type WrapperSendTx struct {
	Wrapper *txs.TxStd
}

var _ txs.ITx = (*WrapperSendTx)(nil)

func NewWrapperSendTx(wrapper *txs.TxStd) WrapperSendTx {
	return WrapperSendTx{Wrapper: wrapper}
}

func (tx WrapperSendTx) ValidateData(ctx context.Context) bool {

	return true
}


func (tx WrapperSendTx) Exec(ctx context.Context) (result btypes.Result, crossTxQcps *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.ABCICodeOK,
	}

	fmt.Println("--------------------------------------------")
	fmt.Println("--------------------------------------------")
	fmt.Println("--------------------------------------------")
	cross := txs.TxQcp{

	}
	crossTxQcps = &cross

	heigth1 := strconv.FormatInt(ctx.BlockHeight(),10)
	tx1 := (types.Tx (ctx.TxBytes())).String()

	fmt.Println(" heigth: " + heigth1 + " hash:" + tx1)

	crossTxQcps.TxStd = tx.Wrapper
	crossTxQcps.To = "main-chain"

	r := btypes.Result{
		Code:0,
	}
	return r, &cross
}

func (tx WrapperSendTx) GetSigner() []btypes.Address {
	return nil
}

func (tx WrapperSendTx) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx WrapperSendTx) GetGasPayer() btypes.Address {
	return nil
}

func (tx WrapperSendTx) GetSignData() []byte {
	return nil
}
