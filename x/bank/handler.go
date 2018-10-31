package bank

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/common"
	starcommon "github.com/QOSGroup/qstars/x/common"
	"strconv"
)
const QSCResultMapperName  =  "qstarsResult"
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

	//set for qos result
	kvMapper := ctx.Mapper(QSCResultMapperName).(*starcommon.KvMapper)
	heigth1 := strconv.FormatInt(ctx.BlockHeight(),10)
	tx1 := (common.HexBytes)(tmhash.Sum(ctx.TxBytes()))
	qstarskey := "heigth:" + heigth1 + ",hash:" + tx1.String()
	fmt.Println(qstarskey)
	qk := []byte(qstarskey)
	kvMapper.Set(qk,"-1")

	crossTxQcps.TxStd = tx.Wrapper
	crossTxQcps.To = "main-chain"

	r := btypes.Result{
		Code:btypes.ABCICodeOK,
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
