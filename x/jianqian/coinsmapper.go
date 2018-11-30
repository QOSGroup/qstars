package jianqian

import (
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/types"
)
const CoinsMapperName = "coins"

type CoinsMapper struct {
	*mapper.BaseMapper
}

type Coins struct{
	Tx         string
	From       types.Address
	Data       []ActivityAward
	Code       string //记录保存跨链结果
}

type ActivityAward struct{
	Address types.Address
	Amount types.BigInt
	CausesCode string
	CausesStr string
}

func NewCoinsMapper(kvMapperName string) *CoinsMapper {
	var txMapper = CoinsMapper{}
	txMapper.BaseMapper = mapper.NewBaseMapper(nil, kvMapperName)
	return &txMapper
}
func (cm *CoinsMapper) Copy() mapper.IMapper {
	cpyMapper := &CoinsMapper{}
	cpyMapper.BaseMapper = cm.BaseMapper.Copy()
	return cpyMapper
}

var _ mapper.IMapper = (*CoinsMapper)(nil)


// Get 查询活动奖励记录
func (cm *CoinsMapper) GetCoins(key []byte) (*Coins, bool) {
	var result Coins
	ok := cm.Get(key, &result)
	return &result, ok
}

// Set 保存活动奖励记录
func (cm *CoinsMapper) SetCoins(i *Coins) {
	cm.Set([]byte(i.Tx), i)
	return
}

// Set 更新跨链记录
func (cm *CoinsMapper) UpdateCoins(key []byte,code string) bool{
	coins,ok:=cm.GetCoins(key)
	if ok{
		coins.Code=code
		cm.SetCoins(coins)
		return true
	}
	return false
}
