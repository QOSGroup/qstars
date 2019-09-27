package supply

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/config"
	starcommon "github.com/QOSGroup/qstars/x/common"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/common"
	"strconv"
	"strings"
)

type OrderTx struct {
	Wrapper     *txs.TxStd
	Address     types.Address //发起者地址
	Id          string        //订单编号
	OrderTo     types.Address //订单接收者地址
	OrderAmount types.BigInt  //订单金额
	Gas         types.BigInt  //gas
}

func NewOrderTx(wrapper *txs.TxStd,address types.Address,id string,orderTo types.Address,orderAmount,gas types.BigInt)*OrderTx{

	return &OrderTx{Wrapper:wrapper,Address:address,Id:id,OrderTo:orderTo,OrderAmount:orderAmount,Gas:gas}

}

func (tx *OrderTx) ValidateData(ctx context.Context) error {
	if tx.Address == nil || tx.Address.Empty() {
		return errors.New("发起者地址不能为空")
	}
	if strings.TrimSpace(tx.Id) == "" {
		return errors.New("订单不能为空")
	}
	if tx.OrderTo == nil || tx.OrderTo.Empty() {
		return errors.New("接单企业地址不能为空")
	}
	return nil
}

//执行业务逻辑,
// crossTxQcp: 需要进行跨链处理的TxQcp。
// 业务端实现中crossTxQcp只需包含`to` 和 `txStd`
func (tx *OrderTx) Exec(ctx context.Context) (result types.Result, crossTxQcps *txs.TxQcp) {
	kvMapper := ctx.Mapper(starcommon.QSCResultMapperName).(*starcommon.KvMapper)
	heigth := strconv.FormatInt(ctx.BlockHeight(), 10)
	txhash := (common.HexBytes)(tmhash.Sum(ctx.TxBytes()))
	key := GetResultKey(heigth, txhash.String())
	kvMapper.Set([]byte(key), OrderStub{}.Name())
	crossTxQcps.TxStd = tx.Wrapper
	crossTxQcps.To = config.GetServerConf().QOSChainName
	crossTxQcps.Extends = key
	result = types.Result{
		Code: types.CodeOK,
	}
	orderMapper := ctx.Mapper(OrderMapperName).(*OrderMapper)
	order:=&Order{Id:tx.Id,OrderName:tx.Address.String(),OrderTo:tx.OrderTo.String(),OrderAmount:tx.OrderAmount,Status:1}
	orderMapper.SaveOrder(key,order)
	return
}

func (tx *OrderTx) GetSigner() []types.Address {
	return []types.Address{tx.Address}
}

func (tx *OrderTx) CalcGas() types.BigInt {
	return tx.Gas
}

func (tx *OrderTx) GetGasPayer() types.Address {
	return types.Address{}
}

func (tx *OrderTx) GetSignData() (ret []byte) {
	ret = append(ret, tx.Wrapper.ITxs[0].GetSignData()...)
	ret = append(ret, tx.Address.Bytes()...)
	ret = append(ret, tx.Id...)
	ret = append(ret, tx.OrderTo...)
	ret = append(ret, tx.OrderAmount.String()...)
	ret = append(ret, tx.Gas.String()...)
	return ret
}

func GetResultKey(heigth string, tx string) string {
	key := "heigth:" + heigth + ",hash:" + tx
	return key
}
