// Copyright 2018 The QOS Authors

package jianqian

import (
	"github.com/QOSGroup/qbase/mapper"
	qbasetypes "github.com/QOSGroup/qbase/types"
	"github.com/tendermint/go-amino"

	"time"
)

const (
	BuyMapperName = "buyad"
)

func getBuyKey(article []byte) []byte {
	return article
}

// Buyer 买家
type Buyer struct {
	Address qbasetypes.Address `json:"address"` // 买家地址
	Buy     qbasetypes.BigInt  `json:"buyad"`   // 购买金额
	BuyTime time.Time          `json:"buyTime"` // 购买时间
}

type BuyMapper struct {
	*mapper.BaseMapper
}

var _ mapper.IMapper = (*BuyMapper)(nil)

func (im *BuyMapper) Copy() mapper.IMapper {
	cpyMapper := &BuyMapper{}
	cpyMapper.BaseMapper = im.BaseMapper.Copy()
	return cpyMapper
}

func NewBuyMapper(cdc *amino.Codec) *BuyMapper {
	var im BuyMapper
	im.BaseMapper = mapper.NewBaseMapper(nil, BuyMapperName)

	return &im
}

// Get 查询用户投资情况
func (im *BuyMapper) GetBuyer(article []byte) (Buyer, bool) {
	key := getBuyKey(article)
	var result Buyer
	ok := im.Get(key, &result)
	return result, ok
}

// Set 添加用户投资
func (im *BuyMapper) SetBuyer(article []byte, i Buyer) {
	key := getBuyKey(article)
	im.Set(key, i)
	return
}
