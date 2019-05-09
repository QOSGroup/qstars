package comm

import (
	"errors"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	qstarstypes "github.com/QOSGroup/qstars/types"
	"strings"

	"github.com/QOSGroup/qstars/x/jianqian"
	"strconv"
)

type Recipient struct {
	Address string
	Amount  types.BigInt
}

type AOETx struct {
	From types.Address
	To   []Recipient
}

func (tx *AOETx) ValidateData(ctx context.Context) error {
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
func (tx *AOETx) Exec(ctx context.Context) (result types.Result, crossTxQcp *txs.TxQcp) {
	accMapper := ctx.Mapper(jianqian.AoeAccountMapperName).(*jianqian.AoeAccountMapper)
	for _, v := range tx.To {
		accMapper.AddBalance(v.Address, "AOE", v.Amount)
	}
	result = types.Result{
		Code: types.CodeOK,
	}
	return
}

func (tx *AOETx) NewTx(args []string) error {
	args_len := len(args)
	if args_len != para_len_3 {
		return errors.New("AdvertisersTx args len error want " + strconv.Itoa(para_len_3) + " got " + strconv.Itoa(args_len))
	}
	address, err := qstarstypes.AccAddressFromBech32(args[0])
	if err != nil {
		return err
	}
	toAddress := args[1]
	tocoins := args[2]
	addrs := strings.Split(toAddress, "|")
	addlen := len(addrs)
	cois := strings.Split(tocoins, "|")
	amounts := make([]types.BigInt, len(cois))
	for i, coinsv := range cois {
		if amou, ok := types.NewIntFromString(coinsv); ok {
			amounts[i] = amou
		} else {
			return errors.New("amount format error")
		}
	}
	toaddrss := make([]types.Address, addlen)
	for i, addrsv := range addrs {
		to, err := qstarstypes.AccAddressFromBech32(addrsv)
		if err != nil {
			return err
		}
		toaddrss[i] = to
	}

	addmap := make(map[string]Recipient)
	for i, coin := range amounts {
		if v, ok := addmap[toaddrss[i].String()]; ok {
			v.Amount.Add(coin)
		} else {
			recipient := Recipient{
				Address: toaddrss[i].String(),
				Amount:  coin,
			}
			addmap[toaddrss[i].String()] = recipient
		}
	}
	var newccs []Recipient
	for _, v := range addmap {
		newccs = append(newccs, v)
	}
	tx.From = address
	tx.To = newccs
	return nil
}
