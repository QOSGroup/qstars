package bank

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	starcommon "github.com/QOSGroup/qstars/x/common"
	"github.com/prometheus/common/log"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/common"
	"strconv"
	"github.com/QOSGroup/qstars/config"
)


type WrapperSendTx struct {
	Wrapper *txs.TxStd
}

var _ txs.ITx = (*WrapperSendTx)(nil)

func NewWrapperSendTx(wrapper *txs.TxStd) WrapperSendTx {
	return WrapperSendTx{Wrapper: wrapper}
}

func (tx WrapperSendTx) ValidateData(ctx context.Context) error {

	return nil
}
func GetResultKey(heigth1 string, tx1 string) string{
	qstarskey := "heigth:" + heigth1 + ",hash:" + tx1
	return qstarskey
}
func (tx WrapperSendTx) Exec(ctx context.Context) (result btypes.Result, crossTxQcps *txs.TxQcp) {

	result = btypes.Result{
		Code: btypes.ABCICodeOK,
	}
	fmt.Println("--------------------------------------------")
	fmt.Println("--------------------------------------------")
	fmt.Println("--------------------------------------------")
	cross := txs.TxQcp{}
	crossTxQcps = &cross

	//set for qos result
	kvMapper := ctx.Mapper(starcommon.QSCResultMapperName).(*starcommon.KvMapper)
	heigth1 := strconv.FormatInt(ctx.BlockHeight(), 10)
	tx1 := (common.HexBytes)(tmhash.Sum(ctx.TxBytes()))
	qstarskey := GetResultKey(heigth1,tx1.String())
	log.Info("new request key:"+qstarskey)
	qk := []byte(qstarskey)
	kvMapper.Set(qk, "-1")

	crossTxQcps.TxStd = tx.Wrapper
	crossTxQcps.To = config.GetServerConf().QOSChainName
	crossTxQcps.Extends = qstarskey

	r := btypes.Result{
		Code: btypes.ABCICodeOK,
	}
	return r, &cross
}

func (tx WrapperSendTx) GetSigner() []btypes.Address {
	return tx.Wrapper.ITx.GetSigner()
}

func (tx WrapperSendTx) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx WrapperSendTx) GetGasPayer() btypes.Address {
	return tx.Wrapper.ITx.GetGasPayer()
}

func (tx WrapperSendTx) GetSignData() []byte {
	return tx.Wrapper.ITx.GetSignData()
}
