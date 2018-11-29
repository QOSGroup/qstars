package coins

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/config"
	starcommon "github.com/QOSGroup/qstars/x/common"
	"github.com/QOSGroup/qstars/x/jianqian"
	"github.com/prometheus/common/log"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/common"
)

//创建AOE  发行AOE  执行活动奖励

//在公链上创建AOE 只能创建一次
type CoinAOETx struct {
	Wrapper *txs.TxStd //已封装好的 TxCreateQSC 结构体
}

func GetResultKey(heigth1 string, tx1 string) string {
	qstarskey := "heigth:" + heigth1 + ",hash:" + tx1
	return qstarskey
}

func (createaoe CoinAOETx) ValidateData(ctx context.Context) error {
	return nil
}

//执行业务逻辑
func (tx CoinAOETx) Exec(ctx context.Context) (result types.Result, crossTxQcp *txs.TxQcp) {
	cross := txs.TxQcp{}
	crossTxQcp = &cross
	kvMapper := ctx.Mapper(QSCResultMapperName).(*starcommon.KvMapper)
	heigth1 := strconv.FormatInt(ctx.BlockHeight(), 10)
	tx1 := (common.HexBytes)(tmhash.Sum(ctx.TxBytes()))
	qstarskey := GetResultKey(heigth1, tx1.String())
	log.Info("CreateAOE new request key:" + qstarskey)
	qk := []byte(qstarskey)
	kvMapper.Set(qk, "-1")
	crossTxQcp.TxStd = tx.Wrapper
	crossTxQcp.To = config.GetServerConf().QOSChainName
	crossTxQcp.Extends = qstarskey
	r := types.Result{
		Code: types.ABCICodeOK,
	}
	return r, &cross
}
func (tx CoinAOETx) GetSigner() []types.Address {
	return []types.Address{}
}

func (tx CoinAOETx) CalcGas() types.BigInt {

	return types.ZeroInt()
}
func (tx CoinAOETx) GetGasPayer() types.Address {

	return types.Address{}
}
func (tx CoinAOETx) GetSignData() []byte {
	return nil
}

func NewCoinAOETx(wrapper *txs.TxStd) CoinAOETx {
	return CoinAOETx{wrapper}
}

////发行 可发行多次
//type IssueAOE struct {
//	Wrapper *txs.TxStd       //已封装好的发行 TxIssueQsc 结构体
//}
//func (tx IssueAOE ) ValidateData(ctx context.Context) error{
//	return nil
//}
//
////执行业务逻辑,
//// crossTxQcp: 需要进行跨链处理的TxQcp。
//// 业务端实现中crossTxQcp只需包含`to` 和 `txStd`
//func (tx IssueAOE ) Exec(ctx context.Context) (result types.Result, crossTxQcp *txs.TxQcp){
//	cross := txs.TxQcp{}
//	crossTxQcp = &cross
//
//	kvMapper := ctx.Mapper(QSCResultMapperName).(*starcommon.KvMapper)
//	heigth1 := strconv.FormatInt(ctx.BlockHeight(), 10)
//	tx1 := (common.HexBytes)(tmhash.Sum(ctx.TxBytes()))
//	qstarskey := GetResultKey(heigth1,tx1.String())
//	log.Info("IssueAOE new request key:"+qstarskey)
//	qk := []byte(qstarskey)
//	kvMapper.Set(qk, "-1")
//
//	crossTxQcp.TxStd = tx.Wrapper
//	crossTxQcp.To = config.GetServerConf().QOSChainName
//	crossTxQcp.Extends = qstarskey
//
//	r := types.Result{
//		Code: types.ABCICodeOK,
//	}
//	return r, &cross
//	}
//func (tx IssueAOE ) GetSigner() []types.Address {
//	return []types.Address{}
//}
//func (tx IssueAOE ) CalcGas() types.BigInt      {
//
//	return types.ZeroInt()
//}
//func (tx IssueAOE ) GetGasPayer() types.Address {
//
//	return  types.Address{}
//}
//func (tx IssueAOE ) GetSignData() []byte {
//	return nil
//}
//活动奖励
type DispatchAOETx struct {
	Wrapper *txs.TxStd //已封装好的 TxCreateQSC 结构体

	From       types.Address
	Address    []types.Address
	CoinAmount []types.BigInt
	CausesCode []string
	CausesStr  []string

	Gas types.BigInt
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
func (tx DispatchAOETx) Exec(ctx context.Context) (result types.Result, crossTxQcp *txs.TxQcp) {
	//本地存储
	coinsMapper := ctx.Mapper(jianqian.CoinsMapperName).(*jianqian.CoinsMapper)
	for i, v := range tx.Address {
		amountstr := fmt.Sprint(tx.CoinAmount[i])
		key := make([]byte, len(v)+len(ctx.TxBytes())+1)
		key = append(key, v...)
		key = append(key, "|"...)
		key = append(key, ctx.TxBytes()...)
		value := tx.CausesCode[i] + "|" + tx.CausesStr[i] + "|" + amountstr
		coinsMapper.Set(key, value)
	}
	//跨链执行
	cross := txs.TxQcp{}
	//crossTxQcp = &cross
	//kvMapper := ctx.Mapper(QSCResultMapperName).(*starcommon.KvMapper)
	//heigth1 := strconv.FormatInt(ctx.BlockHeight(), 10)
	//tx1 := (common.HexBytes)(tmhash.Sum(ctx.TxBytes()))
	//qstarskey := GetResultKey(heigth1, tx1.String())
	//log.Info("DispatchAOETx new request key:" + qstarskey)
	//qk := []byte(qstarskey)
	//kvMapper.Set(qk, "-1")
	//crossTxQcp.TxStd = tx.Wrapper
	//crossTxQcp.To = config.GetServerConf().QOSChainName
	//crossTxQcp.Extends = qstarskey
	r := types.Result{
		Code: types.ABCICodeOK,
	}
	return r, &cross
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

func NewDispatchAOE(Wrapper *txs.TxStd, From types.Address, to []types.Address, coinAmount []types.BigInt, causecode, causestr []string, gas types.BigInt) DispatchAOETx {

	return DispatchAOETx{Wrapper, From, to, coinAmount, causecode, causestr, gas}
}
