package jianqian

import (
	"github.com/QOSGroup/qbase/mapper"
	"github.com/pkg/errors"
)

const (
	AdvertisersMapperName = "advertisers"
)

type AdvertisersMapper struct {
	*mapper.BaseMapper
}



func NewAdvertisersMapper(kvMapperName string) *AdvertisersMapper {
	var txMapper = AdvertisersMapper{}
	txMapper.BaseMapper = mapper.NewBaseMapper(nil, kvMapperName)
	return &txMapper
}
func (cm *AdvertisersMapper) Copy() mapper.IMapper {
	cpyMapper := &AdvertisersMapper{}
	cpyMapper.BaseMapper = cm.BaseMapper.Copy()
	return cpyMapper
}

var _ mapper.IMapper = (*AdvertisersMapper)(nil)

func (cm *AdvertisersMapper) IsAdvertisers(addr string) (bool,error) {
	var result bool
	ok := cm.Get([]byte(addr), &result)
	if ok{
		return result,nil
	}else{
		return false,errors.New("account not exist")
	}
}


func (cm *AdvertisersMapper) SetAdvertisers(addr string)  {
	cm.Set([]byte(addr), true)
}


func (cm *AdvertisersMapper) SetNotAdvertisers(addr string) {
	cm.Set([]byte(addr), false)
}

