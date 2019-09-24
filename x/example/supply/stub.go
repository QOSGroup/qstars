package supply

import (
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/baseapp"
	"github.com/QOSGroup/qstars/x/common"
	"github.com/prometheus/common/log"
	go_amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	"strconv"
)

type OrderStub struct {
}

func NewOrderStub() OrderStub {
	return OrderStub{}
}

func (astub OrderStub) StartX(base *baseapp.QstarsBaseApp) error {
	var orderMapper = NewOrderMapper(OrderMapperName)
	base.Baseapp.RegisterMapper(orderMapper)
	return nil
}
func (astub OrderStub) EndBlockNotify(ctx context.Context) {

}

func (astub OrderStub) RegisterCdc(cdc *go_amino.Codec) {
	cdc.RegisterConcrete(&OrderTx{}, "example/supply/OrderTx", nil)
}

func (astub OrderStub) ResultNotify(ctx context.Context, txQcpResult interface{}) *types.Result {
	in := txQcpResult.(*txs.QcpTxResult)
	log.Debugf("ResultNotify QcpOriginalSequence:%s, result:%+v", string(in.QcpOriginalSequence), txQcpResult)
	var resultCode types.CodeType
	qcpTxResult, ok := baseabci.ConvertTxQcpResult(txQcpResult)
	if ok == false {
		log.Errorf("ResultNotify ConvertTxQcpResult error.")
		resultCode = types.CodeTxDecode
	} else {
		log.Errorf("ResultNotify update status")

		orginalTxHash := in.QcpOriginalExtends //orginalTx.abc
		kvMapper := ctx.Mapper(common.QSCResultMapperName).(*common.KvMapper)
		initValue := ""
		kvMapper.Get([]byte(orginalTxHash), &initValue)
		if initValue != astub.Name() {
			log.Info("This is not my response.")
			return nil
		}
		//put result to map for client query
		c := strconv.FormatInt((int64)(qcpTxResult.Result.Code), 10)
		c = c + " " + qcpTxResult.Result.Log
		kvMapper.Set([]byte(orginalTxHash), c)

		orderMapper := ctx.Mapper(OrderMapperName).(*OrderMapper)
		order := orderMapper.GetOrder(orginalTxHash)
		if order != nil {
			orderMapper.SaveOrder(order.Id, order)
			orderMapper.DeleteOrder(orginalTxHash)
		}
		resultCode = types.CodeOK
	}
	rr := types.Result{
		Code: resultCode,
	}
	return &rr
}

func (cstub OrderStub) CustomerQuery(ctx context.Context, route []string, req abci.RequestQuery) (res []byte, err types.Error) {
	return nil, nil
}
func (cstub OrderStub) Name() string {
	return "OrderStub"
}
