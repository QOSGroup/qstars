package jianqian

import (
	"github.com/QOSGroup/qbase/mapper"
)
const CoinsMapperName = "CoinsMapper"

type CoinsMapper struct {
	*mapper.BaseMapper
}

func NewCoinsMapper(kvMapperName string) *CoinsMapper {
	var txMapper = CoinsMapper{}
	txMapper.BaseMapper = mapper.NewBaseMapper(nil, kvMapperName)
	return &txMapper
}

func (mapper *CoinsMapper) Copy() mapper.IMapper {
	cpyMapper := &CoinsMapper{}
	cpyMapper.BaseMapper = mapper.BaseMapper.Copy()
	return cpyMapper
}

var _ mapper.IMapper = (*CoinsMapper)(nil)

func (mapper *CoinsMapper) SaveKV(key string, value string) {
	mapper.BaseMapper.Set([]byte(key), value)
}

func (mapper *CoinsMapper) GetKey(key string) (v string) {
	mapper.BaseMapper.Get([]byte(key), &v)
	return
}
