// Copyright 2018 The QOS Authors

package jianqian

import (
	"github.com/QOSGroup/qbase/mapper"
	qbasetypes "github.com/QOSGroup/qbase/types"
	"github.com/tendermint/go-amino"

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

func (im *InvestUncheckedMapper) Copy() mapper.IMapper {
	cpyMapper := &InvestUncheckedMapper{}
	cpyMapper.BaseMapper = im.BaseMapper.Copy()
	return cpyMapper
}

func NewInvestUncheckedMapper(cdc *amino.Codec) *InvestUncheckedMapper {
	var im InvestUncheckedMapper
	im.BaseMapper = mapper.NewBaseMapper(nil, InvestUncheckedMapperName)

	return &im
}

// Get 查询用户投资情况
func (im *InvestUncheckedMapper) GetInvestUncheckeds(key []byte) (InvestUncheckeds, bool) {
	var result InvestUncheckeds
	ok := im.Get(key, &result)
	return result, ok
}

// Set 添加用户投资
func (im *InvestUncheckedMapper) SetInvestUncheckeds(key []byte, i InvestUncheckeds) {
	im.Set(key, i)
	return
}
