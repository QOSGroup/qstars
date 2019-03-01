package jianqian

import (
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/types"
	"time"
)
const ArticlesMapperName = "article"

type ArticlesMapper struct {
	*mapper.BaseMapper
}

type Articles struct {
	AuthorAddr          types.Address  //作者地址(必填) 0qos 1cosmos
	AuthorOtherAddr     string         //作者其他帐户地址
	ArticleType         int            //是否原创 0原创 1转载
	ArticleHash         string         //作品唯一标识hash
	ShareAuthor         int            //作者收入比例(必填)
	ShareOriginalAuthor int            //原创收入比例(转载作品必填)
	ShareCommunity      int            //社区收入比例(必填)
	ShareInvestor       int            //投资者收入比例(必填)
	InvestDays          int            //可供投资的天数(必填)
	EndInvestDate       time.Time      //投资结束时间(必填)
	BuyDays             int            //可供购买广告位的天数(必填)
	EndBuyDate          time.Time      //广告位购买结果时间(必填)
	CoinType            string         //币种
	Gas                 types.BigInt
}


func NewArticlesMapper(MapperName string) *ArticlesMapper {
	var txMapper = ArticlesMapper{}
	txMapper.BaseMapper = mapper.NewBaseMapper(nil, MapperName)
	return &txMapper
}

func (mapper *ArticlesMapper) Copy() mapper.IMapper {
	cpyMapper := &ArticlesMapper{}
	cpyMapper.BaseMapper = mapper.BaseMapper.Copy()
	return cpyMapper
}

var _ mapper.IMapper = (*ArticlesMapper)(nil)

func (mapper *ArticlesMapper) SaveKV(key string, value string) {
	mapper.BaseMapper.Set([]byte(key), value)
}

func (mapper *ArticlesMapper) GetKey(key string) (v string) {
	mapper.BaseMapper.Get([]byte(key), &v)
	return
}

func (mapper *ArticlesMapper) GetArticle(articleHash string) *Articles {
	var articles Articles
	exist := mapper.Get([]byte(articleHash), &articles)
	if !exist {
		return nil
	}
	return &articles
}

func (mapper *ArticlesMapper) SetArticle(articleHash string, qscinfo *Articles) bool {
	mapper.Set([]byte(articleHash), qscinfo)
	return true
}