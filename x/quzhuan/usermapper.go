package quzhuan

import (
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/types"
	"github.com/pkg/errors"
)

const (
	UsersMapperName = "Users"
)

var _ mapper.IMapper = (*UsersMapper)(nil)

const Prefix_KYC  ="KYC_"

type UsersMapper struct {
	*mapper.BaseMapper
}

type User struct {
	ID           string            //会员唯一标识
	Source       string            //初次注册场景ID
	Inviter      string            //邀请人

	Address      string            //会员充值地址
	KYCID        string            //KYC id
	Level        string            //会员级别
	QosStake     types.BigInt      //币龄
	RewardAmount map[string]Reward //场景奖励
}

func NewUsersMapper(kvMapperName string) *UsersMapper {
	var txMapper = UsersMapper{}
	txMapper.BaseMapper = mapper.NewBaseMapper(nil, kvMapperName)
	return &txMapper
}
func (cm *UsersMapper) Copy() mapper.IMapper {
	cpyMapper := &UsersMapper{}
	cpyMapper.BaseMapper = cm.BaseMapper.Copy()
	return cpyMapper
}


func (cm *UsersMapper) GetUser(id string) *User {
	var user User
	ok := cm.Get([]byte(id), &user)
	if ok {
		return &user
	}
	return nil
}

func (cm *UsersMapper) AddUser(user User) (bool, error) {
	var result User
	ok := cm.Get([]byte(user.ID), &result)
	if ok {
		return false, errors.New("user already exist")
	}
	cm.Set([]byte(user.ID), user)
	return true, nil
}


func (cm *UsersMapper) AddKYCID(id, kycid string) {
	user := cm.GetUser(id)
	if user != nil {
		if user.KYCID != "" {
			return
		}
		user.KYCID = kycid
		cm.Set([]byte(id), user)
		cm.Set([]byte(Prefix_KYC+kycid), user.ID)
		return
	}
}


func (cm *UsersMapper) KYCExist(kycid string) string {
	var userid string
	ok := cm.Get([]byte(Prefix_KYC+kycid), &userid)
	if ok {
		return userid
	}
	return ""
}


func (cm *UsersMapper) AddAddress(id, address string)  {
	user := cm.GetUser(id)
	if user != nil {
		if user.Address != "" {
			return
		}
		user.Address = address
		cm.Set([]byte(id), user)
	}
	return
}

func (cm *UsersMapper) RemoveKYC(id string) (bool, error) {
	user := cm.GetUser(id)
	if user != nil {
		user.KYCID = ""
		cm.Set([]byte(id), user)
		return true, nil
	}
	return false, errors.New("user not exist")
}

func (cm *UsersMapper) UpdateLevelAndStake(id, level string, stake types.BigInt) (bool, error) {
	user := cm.GetUser(id)
	if user != nil {
		if stake.LT(user.QosStake) {
			return false, errors.New("QosStake must be greater than current")
		}
		user.QosStake = stake
		user.Level = level
		cm.Set([]byte(id), user)
		return true, nil
	}
	return false, errors.New("user not exist")
}

func (cm *UsersMapper) AddReward(reward []*Reward) {
	for _, v := range reward {
		user := cm.GetUser(v.UserID)
		if user != nil {
			if curReward, ok := user.RewardAmount[v.ScenesId]; ok {
				curReward.Amount = curReward.Amount.Add(v.Amount)
			} else {
				newReward := Reward{UserID: v.UserID, ScenesId: v.ScenesId, Amount: v.Amount}
				user.RewardAmount[v.ScenesId] = newReward
			}
			cm.Set([]byte(v.UserID), user)
		}
	}
}
