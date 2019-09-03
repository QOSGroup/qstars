package jianqian

import (
	"encoding/json"
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/types"
	"time"
)
const AuctionMapperName = "auction"

type AuctionMapper struct {
	*mapper.BaseMapper
}

const MAXPRICEKEY  = "0"
type AuctionMap map[string]Auction


type Auction struct{
	Article      string
	Address      types.Address
	CoinsType    string                      // 币种类型
	//OtherAddr    string                      // 转出地址
	Amount       types.BigInt                // 竞拍金额
	AuctionTime  time.Time                   // 最后竞拍时间
}

func NewAuctionMapper(kvMapperName string) *AuctionMapper {
	var txMapper = AuctionMapper{}
	txMapper.BaseMapper = mapper.NewBaseMapper(nil, kvMapperName)
	return &txMapper
}



func (cm *AuctionMapper) Copy() mapper.IMapper {
	cpyMapper := &AuctionMapper{}
	cpyMapper.BaseMapper = cm.BaseMapper.Copy()
	return cpyMapper
}

var _ mapper.IMapper = (*AuctionMapper)(nil)


// Get 查询文章所有竞拍结果
func (cm *AuctionMapper) GetAuction(key string) (result AuctionMap,ok bool) {
	temp,ok:=cm.GetAuctionMap(key)
	if ok{
		delete(temp,MAXPRICEKEY)
	}
	return
}

// Get 查询指定人竞拍情况
func (cm *AuctionMapper) GetAuctionByAddress(article string,address string) (result Auction,exist bool) {
	temp,ok:=cm.GetAuctionMap(article)
	if ok{
		result,exist= temp[address]
	}
	return
}


// Get 获取临时状态  跨链确认使用
func (cm *AuctionMapper) GetTempAuction(key string) (result Auction,exist bool) {
	exist = cm.Get([]byte(key), &result)
	return
}


// Get 查询最高最价信息
func (cm *AuctionMapper) GetMaxAuction(article string) (result Auction,exist bool) {
	temp,ok:=cm.GetAuctionMap(article)
	if ok{
		result,exist= temp[MAXPRICEKEY]
	}
	return
}
// Get 查询最高最价信息
//func (cm *AuctionMapper) GetAuctionMap(article []byte) (result AuctionMap,ok bool) {
//	var temp[] byte
//	ok = cm.Get(article, &temp)
//	var auctionMap AuctionMap
//
//	json.Unmarshal(temp,&auctionMap)
//	return
//}
func (cm *AuctionMapper) GetAuctionMap(article string) (result AuctionMap,ok bool) {
	var temp[] byte
	ok = cm.Get([]byte(article), &temp)
	json.Unmarshal(temp,&result)
	return
}

// Set 保存活动奖励记录
func (cm *AuctionMapper) SetAuction(auction Auction) {
	am,ok:=cm.GetAuctionMap(auction.Article)
	if !ok{
		am=make(AuctionMap)
	}
	am[auction.Address.String()]=auction

	//判断最高出价人
	maxAuction,ex:=cm.GetMaxAuction(auction.Article)
	if !ex{
		am[MAXPRICEKEY]=auction
	}else{
		if maxAuction.Amount.Int64()<auction.Amount.Int64(){
			am[MAXPRICEKEY]=auction
		}
	}
    value,_:=json.Marshal(am)
    cm.Set([]byte(auction.Article), value)
}
// Set 保存活动奖励记录
func (cm *AuctionMapper) SetAuctionByKey(key []byte,i *Auction) {
	cm.Set(key, i)
	return
}

// Set 更新跨链记录
//func (cm *AuctionMapper) UpdateAuction(key []byte,code string) bool{
//	coins,ok:=cm.GetAuction(key)
//	if ok{
//		coins.Status=code
//		cm.SetAuction(coins)
//		return true
//	}
//	return false
//}



// Set 删除用户投资
func (bm *AuctionMapper) DeleteAuction(key []byte) {
	bm.Del(key)
	return
}


