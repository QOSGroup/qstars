// Copyright 2018 The QOS Authors

package jianqian

import (
	"github.com/QOSGroup/qbase/mapper"
	qbasetypes "github.com/QOSGroup/qbase/types"
	"time"
)

const (
	InvestUncheckedMapperName = "investUnchecked"
)

func GetInvestUncheckedKey(article []byte, user qbasetypes.Address) []byte {
	return append(article, user...)
}

// InvestUnchecked 投资者待确认
type InvestUnchecked struct {
	Article    []byte             `json:"article"`    // 文章hash
	Address    qbasetypes.Address `json:"address"`    // 投资者地址
	Invest     qbasetypes.BigInt  `json:"invest"`     // 投资金额
	InvestTime time.Time          `json:"investTime"` // 投资时间
	IsChecked  bool               `json:"isChecked"`  // 已确认
}

type InvestUncheckeds []InvestUnchecked

type InvestUncheckedMapper struct {
	*mapper.BaseMapper
}

var _ mapper.IMapper = (*InvestUncheckedMapper)(nil)

func (ium *InvestUncheckedMapper) Copy() mapper.IMapper {
	cpyMapper := &InvestUncheckedMapper{}
	cpyMapper.BaseMapper = ium.BaseMapper.Copy()
	return cpyMapper
}

func NewInvestUncheckedMapper(mapperName string) *InvestUncheckedMapper {
	var investUncheckedMapper = InvestUncheckedMapper{}
	investUncheckedMapper.BaseMapper = mapper.NewBaseMapper(nil, mapperName)
	return &investUncheckedMapper
}

func (ium *InvestUncheckedMapper) SaveKV(key string, value string) {
	ium.BaseMapper.Set([]byte(key), value)
}

func (ium *InvestUncheckedMapper) GetKey(key string) (v string) {
	ium.BaseMapper.Get([]byte(key), &v)
	return
}

// Get 查询用户投资情况
func (ium *InvestUncheckedMapper) GetInvestUncheckeds(key []byte) (InvestUncheckeds, bool) {
	var result InvestUncheckeds
	ok := ium.Get(key, &result)
	return result, ok
}

// Set 添加用户投资
func (ium *InvestUncheckedMapper) SetInvestUncheckeds(key []byte, i InvestUncheckeds) {
	ium.Set(key, i)
	return
}
