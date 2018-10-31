package common

import (
	"github.com/QOSGroup/qbase/mapper"
)


type KvMapper struct {
	*mapper.BaseMapper
}

func NewKvMapper(kvMapperName string) *KvMapper {
	var txMapper = KvMapper{}
	txMapper.BaseMapper = mapper.NewBaseMapper(nil,kvMapperName)
	return &txMapper
}

func (mapper *KvMapper) Copy() mapper.IMapper {
	cpyMapper := &KvMapper{}
	cpyMapper.BaseMapper = mapper.BaseMapper.Copy()
	return cpyMapper
}


var _ mapper.IMapper = (*KvMapper)(nil)

func (mapper *KvMapper) SaveKV(key string, value string) {
	mapper.BaseMapper.Set([]byte(key), value)
}

func (mapper *KvMapper) GetKey(key string) (v string) {
	mapper.BaseMapper.Get([]byte(key), &v)
	return
}
