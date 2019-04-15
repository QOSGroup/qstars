package recharge

import (
	"errors"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/x/jianqian"
)

type RechargeTx struct {
	Address string
	Tx *jianqian.CoinsTx
}

func (tx RechargeTx) ValidateData(ctx context.Context) error {
	if tx.Tx.Address == nil {
		return errors.New("address must not be empty")
	}
	if tx.Tx.ChangeType == "" || (tx.Tx.ChangeType != jianqian.CHANGE_TYPE_PLUS && tx.Tx.ChangeType != jianqian.CHANGE_TYPE_MINUS) {
		return errors.New("changetype format error")
	}
	if tx.Tx.Amount.IsNil() || tx.Tx.Amount.IsZero() {
		return errors.New("amount format error")
	}
	accMapper := ctx.Mapper(jianqian.AoeAccountMapperName).(*jianqian.AoeAccountMapper)
	//判断余额
	if tx.Tx.ChangeType == jianqian.CHANGE_TYPE_MINUS {
		blance := accMapper.GetBalance(tx.Address, tx.Tx.Cointype)
		if tx.Tx.Amount.GT(blance) {
			//余额不足
			return errors.New("Insufficient balance")
		}
	}
	return nil
}

//执行业务逻辑,
// crossTxQcp: 需要进行跨链处理的TxQcp。
// 业务端实现中crossTxQcp只需包含`to` 和 `txStd`
func (tx RechargeTx) Exec(ctx context.Context) (result types.Result, crossTxQcp *txs.TxQcp) {
	acc := tx.Address
	accMapper := ctx.Mapper(jianqian.AoeAccountMapperName).(*jianqian.AoeAccountMapper)
	//充值
	if tx.Tx.ChangeType == jianqian.CHANGE_TYPE_PLUS {
		accMapper.AddBalance(acc, tx.Tx.Cointype, tx.Tx.Amount)
	} else {
		//提现
		accMapper.SubtractBalance(acc, tx.Tx.Cointype, tx.Tx.Amount)
	}
	result = types.Result{
		Code:  types.CodeOK,
	}
	return
}
func (tx RechargeTx) GetSigner() []types.Address {
	return []types.Address{tx.Tx.Address}

}
func (tx RechargeTx) CalcGas() types.BigInt {
	return types.ZeroInt()
}

func (tx RechargeTx) GetGasPayer() types.Address {
	return tx.Tx.Address

}
func (tx RechargeTx) GetSignData() (ret []byte) {
	ret = append(ret, []byte(tx.Address)...)
	ret = append(ret, tx.Tx.Address.Bytes()...)
	ret = append(ret, types.Int2Byte(tx.Tx.Amount.Int64())...)
	ret = append(ret, []byte(tx.Tx.Cointype)...)
	ret = append(ret, []byte(tx.Tx.ChangeType)...)
	return
}
