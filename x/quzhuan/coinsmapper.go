package quzhuan

import (
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/types"
)

const (
	CoinsMapperName = "Coins"
)

var _ mapper.IMapper = (*CoinsMapper)(nil)

const Prefix_User  = "User_"
const Prefix_Scenes  ="Scenes_"


type CoinsMapper struct {
	*mapper.BaseMapper
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


// Set 更新账户 增加余额
func (cm *CoinsMapper) UserAddBalance(id string, amount types.BigInt) {
	id=getUserid(id)
	var oldAmount types.BigInt
	ok := cm.Get([]byte(id), &oldAmount)
	if ok {
		amount=oldAmount.Add(amount)
	}
	cm.Set([]byte(id),amount)
	return
}

// Set 更新账户 增加余额
func (cm *CoinsMapper) ScenesAddBalance(id string, amount types.BigInt) {
	id=getScenesid(id)
	var oldAmount types.BigInt
	ok := cm.Get([]byte(id), &oldAmount)
	if ok {
		amount=oldAmount.Add(amount)
	}
	cm.Set([]byte(id),amount)
	return
}


// Set 更新账户 减少余额
func (cm *CoinsMapper) UserSubtractBalance(id string, amount types.BigInt) {
	id=getUserid(id)
	var oldAmount types.BigInt
	ok := cm.Get([]byte(id), &oldAmount)
	if ok {
		amount=oldAmount.Sub(amount)
	}
	cm.Set([]byte(id),amount)
	return
}


// Set 更新账户 减少余额
func (cm *CoinsMapper) ScenesSubtractBalance(id string, amount types.BigInt) {
	id=getScenesid(id)
	var oldAmount types.BigInt
	ok := cm.Get([]byte(id), &oldAmount)
	if ok {
		amount=oldAmount.Sub(amount)
	}
	cm.Set([]byte(id),amount)
	return
}



func getUserid(id string) string{
	return 	Prefix_User+id
}


func getScenesid(id string) string{
	return 	Prefix_Scenes+id
}