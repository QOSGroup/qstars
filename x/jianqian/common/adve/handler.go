package adver

import (
"github.com/QOSGroup/qbase/context"
"github.com/QOSGroup/qbase/txs"
"github.com/QOSGroup/qbase/types"
"github.com/QOSGroup/qstars/x/jianqian"
"github.com/pkg/errors"
)

type AdvertisersTx struct {
	Tx *jianqian.CoinsTx
}

func  ValidateData(ctx context.Context,tx AdvertisersTx) error {
	if tx.Tx.Address == nil {
		return errors.New("address must not be empty")
	}
	if tx.Tx.ChangeType == "" || (tx.Tx.ChangeType != jianqian.CHANGE_TYPE_PLUS && tx.Tx.ChangeType != jianqian.CHANGE_TYPE_MINUS) {
		return errors.New("changetype format error")
	}
	if tx.Tx.Amount.IsNil() || tx.Tx.Amount.IsZero() {
		return errors.New("amount format error")
	}
	adverMapper := ctx.Mapper(jianqian.AdvertisersMapperName).(*jianqian.AdvertisersMapper)
	accMapper := ctx.Mapper(jianqian.AoeAccountMapperName).(*jianqian.AoeAccountMapper)
	isaddver, err := adverMapper.IsAdvertisers(tx.Tx.Address.String())
	//成为广告商 相当于提现  减少余额  时判断余额
	if tx.Tx.ChangeType == jianqian.CHANGE_TYPE_MINUS {
		if isaddver {
			return errors.New(tx.Tx.Address.String() + " Already an advertiser")
		}
		blance := accMapper.GetBalance(tx.Tx.Address.String(), tx.Tx.Cointype)
		if tx.Tx.Amount.GT(blance) {
			//余额不足
			return errors.New("Insufficient balance")
		}
	} else {
		//赎回时判断当前是否是广告商
		if err != nil {
			return errors.New("advertisers account not exist")
		}
		if !isaddver {
			return errors.New(tx.Tx.Address.String() + " not advertisers")
		}
	}
	return nil
}

//执行业务逻辑,
// crossTxQcp: 需要进行跨链处理的TxQcp。
// 业务端实现中crossTxQcp只需包含`to` 和 `txStd`
func  Exec(ctx context.Context,tx AdvertisersTx) (result types.Result, crossTxQcp *txs.TxQcp) {
	acc:=tx.Tx.Address.String()
	adverMapper := ctx.Mapper(jianqian.AdvertisersMapperName).(*jianqian.AdvertisersMapper)
	accMapper := ctx.Mapper(jianqian.AoeAccountMapperName).(*jianqian.AoeAccountMapper)

	//赎回  增加余额
	if tx.Tx.ChangeType==jianqian.CHANGE_TYPE_PLUS{
		accMapper.AddBalance(acc,tx.Tx.Cointype,tx.Tx.Amount)
		//设置为非广告主身份
		adverMapper.SetNotAdvertisers(acc)
	}else{
		//抵押
		accMapper.SubtractBalance(acc,tx.Tx.Cointype,tx.Tx.Amount)
		//设置为广告主身份
		adverMapper.SetAdvertisers(acc)
	}
	result = types.Result{
		Code:  types.CodeOK,
	}
	return

}



func GetStruct(args []string) AdvertisersTx{




}
