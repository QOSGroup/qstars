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
	Address types.Address
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
func (cm *AoeAccountMapper) GetBalance(key []byte) types.BigInt {
	var balance types.BigInt = types.ZeroInt()
	var result AoeAccount
	ok := cm.Get(key, &result)
	if ok {
		balance = result.Amount
	}
	return balance
}


// Set 更新账户 增加余额
func (cm *AoeAccountMapper) AddBalance(key types.Address,amount types.BigInt) {
	balance := cm.GetBalance([]byte(key))
	balance.Add(amount)
	cm.Set(key.Bytes(), balance)
	return
}


// Set 更新账户 增加余额
func (cm *AoeAccountMapper) AddBalanceBatch(accs []AoeAccount) {
	for _,v:=range accs{
		key:=v.Address
		amount:=v.Amount
		cm.AddBalance(key,amount)
	}
	return
}

// Set 更新账户 减少余额
func (cm *AoeAccountMapper) SubtractBalance(key types.Address,amount types.BigInt) {
	balance := cm.GetBalance(key.Bytes())
	balance.Sub(amount)
	cm.Set(key.Bytes(), balance)
	return
}
