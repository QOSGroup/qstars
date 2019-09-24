package supply

import (
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/types"
)

const (
	OrderMapperName = "order"
)

type OrderMapper struct {
	*mapper.BaseMapper
}

type Order struct {
	Id          string       //订单编号
	OrderName   string       //订单名称
	OrderFrom   string       //订单发起者
	OrderTo     string       //订单接收者
	OrderAmount types.BigInt //订单金额
	Status      int          //转账状态 0成功  1未成功
}

var _ mapper.IMapper = (*OrderMapper)(nil)

func NewOrderMapper(kvMapperName string) *OrderMapper {
	var txMapper = OrderMapper{}
	txMapper.BaseMapper = mapper.NewBaseMapper(nil, kvMapperName)
	return &txMapper
}
func (mapper *OrderMapper) Copy() mapper.IMapper {
	cpyMapper := &OrderMapper{}
	cpyMapper.BaseMapper = mapper.BaseMapper.Copy()
	return cpyMapper
}

func (mapper *OrderMapper) SaveOrder(key string,order *Order) {
	mapper.Set([]byte(key), order)
}

func (mapper *OrderMapper) GetOrder(key string) *Order {
	var order Order
	exist := mapper.Get([]byte(key), &order)
	if !exist {
		return nil
	}
	return &order
}

// Set 删除用户投资
func (mapper *OrderMapper) DeleteOrder(key string) {
	mapper.Del([]byte(key))
}
