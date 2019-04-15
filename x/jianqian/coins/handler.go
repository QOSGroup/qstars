package coins

import (
	"errors"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/x/jianqian"
)

//转账 支持一转多
type CoinsTx struct {
	From     types.Address
	CoinType string
	To       []Recipient
}

type Recipient struct {
	Address string
	Amount  types.BigInt
}

func (tx CoinsTx) ValidateData(ctx context.Context) error {
	if tx.From == nil {
		return errors.New("address must not be empty")
	}
	if tx.CoinType == "" {
		return errors.New("CoinType format error")
	}
	totle := types.ZeroInt()
	for _, v := range tx.To {
		if v.Amount.IsNil() || v.Amount.IsZero() {
			return errors.New("amount must GT 0")
		}
		totle = totle.Add(v.Amount)
	}
	accMapper := ctx.Mapper(jianqian.AoeAccountMapperName).(*jianqian.AoeAccountMapper)
	//判断余额
	blance := accMapper.GetBalance(tx.From.String(), tx.CoinType)
	if totle.GT(blance) {
		//余额不足
		return errors.New(tx.From.String()+" 余额不足")
	}
	return nil
}

//执行业务逻辑,
// crossTxQcp: 需要进行跨链处理的TxQcp。
// 业务端实现中crossTxQcp只需包含`to` 和 `txStd`
func (tx CoinsTx) Exec(ctx context.Context) (result types.Result, crossTxQcp *txs.TxQcp) {
	acc := tx.From.String()
	accMapper := ctx.Mapper(jianqian.AoeAccountMapperName).(*jianqian.AoeAccountMapper)
	totle := types.ZeroInt()
	for _, v := range tx.To {
		totle = totle.Add(v.Amount)
		accMapper.AddBalance(v.Address, tx.CoinType, v.Amount)
	}
	accMapper.SubtractBalance(acc, tx.CoinType, totle)
	result = types.Result{
		Code: types.CodeOK,
	}
	return
}
func (tx CoinsTx) GetSigner() []types.Address {
	return []types.Address{tx.From}
}
func (tx CoinsTx) CalcGas() types.BigInt {
	return types.ZeroInt()
}

func (tx CoinsTx) GetGasPayer() types.Address {
	return tx.From
}
func (tx CoinsTx) GetSignData() (ret []byte) {
	ret = append(ret, tx.From.Bytes()...)
	ret = append(ret, []byte(tx.CoinType)...)
	for _, v := range tx.To {
		ret = append(ret, []byte(v.Address)...)
		ret = append(ret, types.Int2Byte(v.Amount.Int64())...)
	}
	return
}



type AOETx struct {
	From     types.Address
	To       []Recipient
}

func (tx AOETx) ValidateData(ctx context.Context) error {
	if tx.From == nil {
		return errors.New("address must not be empty")
	}
	for _, v := range tx.To {
		if v.Amount.IsNil() || v.Amount.IsZero() {
			return errors.New("amount must GT 0")
		}
	}
	return nil
}

//执行业务逻辑,
// crossTxQcp: 需要进行跨链处理的TxQcp。
// 业务端实现中crossTxQcp只需包含`to` 和 `txStd`
func (tx AOETx) Exec(ctx context.Context) (result types.Result, crossTxQcp *txs.TxQcp) {
	accMapper := ctx.Mapper(jianqian.AoeAccountMapperName).(*jianqian.AoeAccountMapper)
	for _, v := range tx.To {
		accMapper.AddBalance(v.Address, "AOE", v.Amount)
	}
	result = types.Result{
		Code: types.CodeOK,
	}
	return
}
func (tx AOETx) GetSigner() []types.Address {
	return []types.Address{tx.From}
}
func (tx AOETx) CalcGas() types.BigInt {
	return types.ZeroInt()
}

func (tx AOETx) GetGasPayer() types.Address {
	return tx.From
}
func (tx AOETx) GetSignData() (ret []byte) {
	ret = append(ret, tx.From.Bytes()...)
	for _, v := range tx.To {
		ret = append(ret, []byte(v.Address)...)
		ret = append(ret, types.Int2Byte(v.Amount.Int64())...)
	}
	return
}