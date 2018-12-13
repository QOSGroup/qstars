package coins

import (
	"errors"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/x/common"
	"github.com/QOSGroup/qstars/x/jianqian"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmcommon "github.com/tendermint/tendermint/libs/common"
	"strconv"
)

//活动奖励
type DispatchAOETx struct {
	Wrapper *txs.TxStd //已封装好的 TxCreateQSC 结构体

	From       types.Address
	Address    []types.Address
	CoinAmount []types.BigInt
	CausesCode []string
	CausesStr  []string
	Gas        types.BigInt
}

func (tx DispatchAOETx) ValidateData(ctx context.Context) error {
	if len(tx.Address) == 0 {
		return errors.New("DispatchAOE address must not empty")
	}
	if len(tx.Address) > 100 {
		return errors.New("DispatchAOE address number more than 100")
	}
	if len(tx.Address) != len(tx.CoinAmount) || len(tx.Address) != len(tx.CausesCode) || len(tx.Address) != len(tx.CausesStr) {
		return errors.New("DispatchAOE address|amount|causes nnequal length")
	}
	return nil
}

//执行业务逻辑,
// crossTxQcp: 需要进行跨链处理的TxQcp。
// 业务端实现中crossTxQcp只需包含`to` 和 `txStd`
func (tx DispatchAOETx) Exec(ctx context.Context) (result types.Result, crossTxQcps *txs.TxQcp) {

	awards := make([]jianqian.ActivityAward, 0)
	coinsMapper := ctx.Mapper(jianqian.CoinsMapperName).(*jianqian.CoinsMapper)
	for i, v := range tx.Address {
		award := jianqian.ActivityAward{v, tx.CoinAmount[i], tx.CausesCode[i], tx.CausesStr[i]}
		awards = append(awards, award)
	}
	tx1 := (tmcommon.HexBytes)(tmhash.Sum(ctx.TxBytes()))
	heigth1 := strconv.FormatInt(ctx.BlockHeight(), 10)
	key := "heigth:" + heigth1 + ",hash:" + tx1.String()

	coins := jianqian.Coins{tx1.String(), tx.From, awards, "-1"}
	coinsMapper.SetCoins(&coins)

	kvMapper := ctx.Mapper(common.QSCResultMapperName).(*common.KvMapper)
	kvMapper.Set([]byte(key), CoinsStub{}.Name())

	//跨链
	crossTxQcps = &txs.TxQcp{}
	crossTxQcps.TxStd = tx.Wrapper
	crossTxQcps.To = config.GetServerConf().QOSChainName
	crossTxQcps.Extends=key
	result = types.Result{
		Code: types.ABCICodeOK,
	}
	return
}

func (tx DispatchAOETx) GetSigner() []types.Address {
	return []types.Address{tx.From}
}
func (tx DispatchAOETx) CalcGas() types.BigInt {
	return tx.Gas
}
func (tx DispatchAOETx) GetGasPayer() types.Address {
	return types.Address{}
}
func (tx DispatchAOETx) GetSignData() (ret []byte) {
	ret = append(ret, tx.Wrapper.ITx.GetSignData()...)
	ret = append(ret, []byte(tx.From)...)
	for i, _ := range tx.Address {
		ret = append(ret, tx.Address[i]...)
		ret = append(ret, types.Int2Byte(tx.CoinAmount[i].Int64())...)
		ret = append(ret, tx.CausesCode[i]...)
		ret = append(ret, tx.CausesStr[i]...)
	}
	return ret
}

func (tx DispatchAOETx) Name() string {
	return "DispatchAOETx"
}

func NewDispatchAOE(Wrapper *txs.TxStd, From types.Address, to []types.Address, coinAmount []types.BigInt, causecode, causestr []string, gas types.BigInt) DispatchAOETx {

	return DispatchAOETx{Wrapper, From, to, coinAmount, causecode, causestr, gas}
}
