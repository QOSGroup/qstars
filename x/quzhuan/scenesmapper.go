package quzhuan

import (
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/types"
)

const (
	ScenesMapperName = "Scenes"
)

var _ mapper.IMapper = (*ScenesMapper)(nil)

const Prefix_ID  = "ID_"
const Prefix_Address  ="Address_"

type ScenesMapper struct {
	*mapper.BaseMapper
}

type Scenes struct {
	ID           string            //场景唯一标识
	Address      string            //场景地址
	Status       int
	RewardAmount map[string]Reward //场景奖励累计
}

type Reward struct {
	UserID string
	ScenesId string
	Amount types.BigInt
}

func NewScenesMapper(kvMapperName string) *ScenesMapper {
	var txMapper = ScenesMapper{}
	txMapper.BaseMapper = mapper.NewBaseMapper(nil, kvMapperName)
	return &txMapper
}
func (cm *ScenesMapper) Copy() mapper.IMapper {
	cpyMapper := &ScenesMapper{}
	cpyMapper.BaseMapper = cm.BaseMapper.Copy()
	return cpyMapper
}


func (cm *ScenesMapper) AddScenes(scenes Scenes)  {
	var result Scenes
	ok := cm.Get([]byte(scenes.ID), &result)
	if ok {
		return
	}
	cm.Set([]byte(Prefix_ID+scenes.ID), scenes)
	cm.Set([]byte(Prefix_Address+scenes.Address), scenes)
}




func (cm *ScenesMapper) GetScenesByID(id string) *Scenes {
	id=Prefix_ID+id
	var scenes Scenes
	ok := cm.Get([]byte(id), &scenes)
	if ok {
		return &scenes
	}
	return nil
}


func (cm *ScenesMapper) GetScenesByAddress(id string) *Scenes {
	id=Prefix_Address+id
	var scenes Scenes
	ok := cm.Get([]byte(id), &scenes)
	if ok {
		return &scenes
	}
	return nil
}


func (cm *ScenesMapper) AddReward(scenesId string,reward []*Reward) {
	scenes:=cm.GetScenesByID(scenesId)
	if scenes!=nil{
		curReward:=scenes.RewardAmount
		for _,v:=range reward{
			if  oldReward,ok:=curReward[v.UserID];ok{
				oldReward.Amount=oldReward.Amount.Add(v.Amount)
			}else{
				newReward:=Reward{UserID:v.UserID,ScenesId:v.ScenesId,Amount:v.Amount}
				curReward[v.UserID]=newReward
			}
		}
		cm.Set([]byte(scenesId),scenes)
	}
	return
}
