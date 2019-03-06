package jianqian

import (
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/types"
)

const AoeAccountMapperName = "aoeaccount"

type AoeAccountMapper struct {
	*mapper.BaseMapper
}

//账户
type AoeAccount struct {
	Address string
	Amount  types.BigInt
}

func NewAccountMapper(kvMapperName string) *AoeAccountMapper {
	var txMapper = AoeAccountMapper{}
	txMapper.BaseMapper = mapper.NewBaseMapper(nil, kvMapperName)
	return &txMapper
}
func (cm *AoeAccountMapper) Copy() mapper.IMapper {
	cpyMapper := &AoeAccountMapper{}
	cpyMapper.BaseMapper = cm.BaseMapper.Copy()
	return cpyMapper
}

var _ mapper.IMapper = (*AoeAccountMapper)(nil)

// Get 查询账户余额
func (cm *AoeAccountMapper) GetBalance(key string) types.BigInt {
	var balance = types.ZeroInt()
	var result AoeAccount
	ok := cm.Get([]byte(key), &result)
	if ok {
		balance = result.Amount
	}
	return balance
}


// Set 更新账户 增加余额
func (cm *AoeAccountMapper) AddBalance(aoe AoeAccount) {
	balance := cm.GetBalance(aoe.Address)
	balance=balance.Add(aoe.Amount)
	aoe.Amount=balance
	cm.Set([]byte(aoe.Address), aoe)

	return
}


// Set 更新账户 增加余额
func (cm *AoeAccountMapper) AddBalanceBatch(accs []AoeAccount) {
	for _,v:=range accs{
		cm.AddBalance(v)
	}
	return
}

// Set 更新账户 减少余额
func (cm *AoeAccountMapper) SubtractBalance(key types.Address,amount types.BigInt) {
	balance := cm.GetBalance(key.String())
	balance.Sub(amount)
	cm.Set(key.Bytes(), balance)
	return
}
