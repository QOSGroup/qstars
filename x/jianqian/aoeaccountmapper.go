package jianqian

import (
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/types"
	"github.com/pkg/errors"
	"strings"
)

const (
	AoeAccountMapperName = "aoeaccount"
	CHANGE_TYPE_PLUS="1"
	CHANGE_TYPE_MINUS="2"
)

type AoeAccountMapper struct {
	*mapper.BaseMapper
}

type CoinsTx struct {
	Address types.Address
	Cointype string
	Amount types.BigInt
	ChangeType string     //0 plus  1 minus
}

type AOETokens = types.BaseCoins

func NewAccountMapper(kvMapperName string) *AoeAccountMapper {
	var txMapper = AoeAccountMapper{}
	txMapper.BaseMapper = mapper.NewBaseMapper(nil, kvMapperName)
	return &txMapper
}
func (cm *AoeAccountMapper) Copy() mapper.IMapper {
	cpyMapper := &AoeAccountMapper{}
	cpyMapper.BaseMapper = cm.BaseMapper.Copy()
	return cpyMapper
}

var _ mapper.IMapper = (*AoeAccountMapper)(nil)

// Get 查询指定币种账户余额
func (cm *AoeAccountMapper) GetBalance(addr,cointype string) types.BigInt {
	var balance = types.ZeroInt()
	var result AOETokens
	ok := cm.Get([]byte(addr), &result)
	if ok {
		for _,v:=range result{
			if v.Name==strings.ToUpper(cointype){
				return v.Amount
			}
		}
	}
	return balance
}
// Get 查询账户余额
func (cm *AoeAccountMapper) GetAllBalanceByAddr(addr string) (AOETokens,error) {
	var result AOETokens
	ok := cm.Get([]byte(addr), &result)
	if ok {
		return result,nil
	}
	return nil,errors.New(addr+" account not exist")
}

// Set 更新账户 增加余额
func (cm *AoeAccountMapper) AddBalance(addr ,cointype string, amount types.BigInt) {
	var newtokens AOETokens
	result:=types.BaseCoins{&types.BaseCoin{strings.ToUpper(cointype),amount}}
	oldTokens,err:=cm.GetAllBalanceByAddr(addr)
	if err!=nil{
		newtokens= result
	}else{
		newtokens=oldTokens.Plus(result)
	}
	cm.Set([]byte(addr), newtokens)
	return
}


// Set 更新账户 减少余额
func (cm *AoeAccountMapper) SubtractBalance(addr,cointype string,amount types.BigInt) {
	var newtokens AOETokens
	result:=types.BaseCoins{&types.BaseCoin{strings.ToUpper(cointype),amount}}
	oldTokens,err:=cm.GetAllBalanceByAddr(addr)
	if err!=nil{
		newtokens= result
	}else{
		newtokens=oldTokens.Minus(result)
	}
	cm.Set([]byte(addr), newtokens)
	return
}
