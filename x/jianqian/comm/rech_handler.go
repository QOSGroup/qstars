package comm

import (
	"errors"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	qstarstypes "github.com/QOSGroup/qstars/types"

	"github.com/QOSGroup/qstars/x/jianqian"
	"strconv"
)

//充值
type RechargeTx struct {
	Tx *jianqian.CoinsTx
}

var _ RouterTx = (*RechargeTx)(nil)

func (tx *RechargeTx) ValidateData(ctx context.Context) error {
	if tx.Tx.Address == nil {
		return errors.New("address must not be empty")
	}
	if tx.Tx.ChangeType == "" || (tx.Tx.ChangeType != jianqian.CHANGE_TYPE_PLUS && tx.Tx.ChangeType != jianqian.CHANGE_TYPE_MINUS) {
		return errors.New("changetype format error")
	}
	if tx.Tx.Amount.IsNil() || tx.Tx.Amount.IsZero() {
		return errors.New("amount format error")
	}
	return nil
}

//执行业务逻辑,
// crossTxQcp: 需要进行跨链处理的TxQcp。
// 业务端实现中crossTxQcp只需包含`to` 和 `txStd`
func (tx *RechargeTx) Exec(ctx context.Context) (result types.Result, crossTxQcp *txs.TxQcp) {
	acc := tx.Tx.Address.String()
	accMapper := ctx.Mapper(jianqian.AoeAccountMapperName).(*jianqian.AoeAccountMapper)
	//充值
	//if tx.Tx.ChangeType == jianqian.CHANGE_TYPE_PLUS {
	accMapper.AddBalance(acc, tx.Tx.Cointype, tx.Tx.Amount)
	//} else {
	//	//提现
	//	accMapper.SubtractBalance(acc, tx.Tx.Cointype, tx.Tx.Amount)
	//}
	result = types.Result{
		Code: types.CodeOK,
	}
	return
}

func (tx *RechargeTx) NewTx(args []string, address types.Address) error {
	args_len := len(args)
	if args_len != para_len_5 {
		return errors.New("AdvertisersTx args len error want " + strconv.Itoa(para_len_5) + " got " + strconv.Itoa(args_len))
	}

	//充值是代签  不能用私钥地址
	address, err := qstarstypes.AccAddressFromBech32(args[1])
	if err != nil {
		return err
	}
	//tx.Address = args[0]
	coinsTx, err := GetCoins(address, args[2], args[3], args[4])
	if err != nil {
		return err
	}
	tx.Tx = coinsTx
	return nil
}

//提现
type ExtractTx struct {
	Tx *jianqian.CoinsTx
}

var _ RouterTx = (*ExtractTx)(nil)

func (tx *ExtractTx) ValidateData(ctx context.Context) error {
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
	blance := accMapper.GetBalance(tx.Tx.Address.String(), tx.Tx.Cointype)
	if tx.Tx.Amount.GT(blance) {
		//余额不足
		return errors.New("Insufficient balance")
	}
	return nil
}

//执行业务逻辑,
// crossTxQcp: 需要进行跨链处理的TxQcp。
// 业务端实现中crossTxQcp只需包含`to` 和 `txStd`
func (tx *ExtractTx) Exec(ctx context.Context) (result types.Result, crossTxQcp *txs.TxQcp) {
	acc := tx.Tx.Address.String()
	accMapper := ctx.Mapper(jianqian.AoeAccountMapperName).(*jianqian.AoeAccountMapper)
	//充值
	//if tx.Tx.ChangeType == jianqian.CHANGE_TYPE_PLUS {
	//	accMapper.AddBalance(acc, tx.Tx.Cointype, tx.Tx.Amount)
	//} else {
	//	//提现
	accMapper.SubtractBalance(acc, tx.Tx.Cointype, tx.Tx.Amount)
	//}
	result = types.Result{
		Code: types.CodeOK,
	}
	return
}

func (tx *ExtractTx) NewTx(args []string, address types.Address) error {
	args_len := len(args)
	if args_len != para_len_4 {
		return errors.New("AdvertisersTx args len error want " + strconv.Itoa(para_len_4) + " got " + strconv.Itoa(args_len))
	}
	coinsTx, err := GetCoins(address, args[1], args[2], args[3])
	if err != nil {
		return err
	}
	tx.Tx = coinsTx
	return nil
}
