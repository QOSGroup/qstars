// Copyright 2018 The QOS Authors

package jianqian

import (
	"github.com/QOSGroup/qbase/mapper"
	qbasetypes "github.com/QOSGroup/qbase/types"
	"time"
)

const (
	BuyMapperName = "buyad"
)

func getBuyKey(article []byte) []byte {
	return article
}

type CheckStatus int

const (
	CheckStatusInit CheckStatus = iota
	CheckStatusSuccess
	CheckStatusFail
)

// Buyer 买家
type Buyer struct {
	Address     qbasetypes.Address `json:"address"`   // 买家地址
	Buy         qbasetypes.BigInt  `json:"buyad"`     // 购买金额
	BuyTime     time.Time          `json:"buyTime"`   // 购买时间
	CheckStatus CheckStatus        `json:"isChecked"` // 验证状态
}

type BuyMapper struct {
	*mapper.BaseMapper
}

var _ mapper.IMapper = (*BuyMapper)(nil)

func (bm *BuyMapper) Copy() mapper.IMapper {
	cpyMapper := &BuyMapper{}
	cpyMapper.BaseMapper = bm.BaseMapper.Copy()
	return cpyMapper
}

func NewBuyMapper(mapperName string) *BuyMapper {
	var buyMapper = BuyMapper{}
	buyMapper.BaseMapper = mapper.NewBaseMapper(nil, mapperName)
	return &buyMapper
}

func (bm *BuyMapper) SaveKV(key string, value string) {
	bm.BaseMapper.Set([]byte(key), value)
}

func (bm *BuyMapper) GetKey(key string) (v string) {
	bm.BaseMapper.Get([]byte(key), &v)
	return
}

// Get 查询用户投资情况
func (bm *BuyMapper) GetBuyer(article []byte) (*Buyer, bool) {
	key := getBuyKey(article)
	var result Buyer
	ok := bm.Get(key, &result)
	return &result, ok
}

// Set 添加用户投资
func (bm *BuyMapper) SetBuyer(article []byte, i Buyer) {
	key := getBuyKey(article)
	bm.Set(key, i)
	return
}

// Set 删除用户投资
func (bm *BuyMapper) DeleteBuyer(article []byte) {
	key := getBuyKey(article)
	bm.Del(key)
	return
}
