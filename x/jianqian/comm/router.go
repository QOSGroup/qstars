package comm

import (
	"errors"
	"github.com/QOSGroup/qbase/context"

	"github.com/QOSGroup/qbase/txs"
	qbasetypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/x/jianqian"
)

const (
	para_len_1  int = 1
	para_len_2  int = 2
	para_len_3  int = 3
	para_len_4  int = 4
	para_len_5  int = 5
	para_len_6  int = 6
	para_len_7  int = 7
	para_len_8  int = 8
	para_len_9  int = 9
	para_len_10 int = 10
	para_len_11 int = 11
)

type RouterTx interface {
	ValidateData(ctx context.Context) error //检测
	Exec(ctx context.Context) (result qbasetypes.Result, crossTxQcp *txs.TxQcp)
	NewTx(args []string,address qbasetypes.Address) error
}

func getStruct(funcName string, args []string,address qbasetypes.Address) (routertx RouterTx, err error) {
	switch funcName {
	case AdvertisersTxFlag:
		routertx = &AdvertisersTx{}
	case ArticleTxFlag:
		routertx = &ArticleTx{}
	case AuctionTxFlag:
		routertx = &AuctionTx{}
	case BuyTxFlag:
		routertx = &BuyTx{}
	case AOETxFlag:
		routertx = &AOETx{}
	case InvestTxFlag:
		routertx = &InvestTx{}
	case RechargeTxFlag:
		routertx = &RechargeTx{}
	case ExtractTxFlag:
		routertx = &ExtractTx{}
	default:
		err = errors.New(funcName + " funcName not support")
	}
	if err != nil {
		return
	}
	if routertx != nil {
		err = routertx.NewTx(args,address)
	}
	return
}

func GetCoins(address qbasetypes.Address, cointype, changetype, amount string) (*jianqian.CoinsTx, error) {
	coinsTx := &jianqian.CoinsTx{}
	//address, err := types.AccAddressFromBech32(addr)
	//if err != nil {
	//	return nil, err
	//}
	coins, ok := qbasetypes.NewIntFromString(amount)
	if !ok {
		return nil, errors.New("amount format error")
	}
	coinsTx.Address = address
	coinsTx.Cointype = cointype
	coinsTx.ChangeType = changetype
	coinsTx.Amount = coins
	return coinsTx, nil
}
