package common

import (
	"encoding/json"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/x/quzhuan"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

const (

	//注册场景
	Scenes_RegisterScenes = "registerScenes"
	//活动奖励
	Scenes_ScenesReward   = "scenesReward"
    //场景充值
	Scenes_Recharge = "recharge"

	//只有合约owner才可以添加场景
	OWNER_Address = ""
)

type ScenesTx struct {
	Address       types.Address     //签名地址
	ScenesID      string            //场景唯一标识
	ScenesAddress string            //场景地址
	Rewards       []*quzhuan.Reward //奖励明细json
	Status        int               //状态
	Amount        types.BigInt      //充值金额
	FuncName      string            //方法名
}

type RewardDetail struct {
	UserId string //场景唯一标识
	Amount types.BigInt
}
type ScenesReward struct {
	ScenesId string
	Rewards  []RewardDetail
}

func (tx *ScenesTx) String() string {
	return tx.ScenesID + "|" + tx.ScenesAddress
}

var _ RouterTx = (*ScenesTx)(nil)

func (tx *ScenesTx) ValidateData(ctx context.Context) error {
	if strings.TrimSpace(tx.ScenesID) == "" {
		return errors.New("场景ID 不能为空")
	}
	scenesMapper := ctx.Mapper(quzhuan.ScenesMapperName).(*quzhuan.ScenesMapper)
	scenes := scenesMapper.GetScenesByID(tx.ScenesID)
	//添加场景 或场景充值需要权限验证
	if tx.FuncName == Scenes_RegisterScenes || tx.FuncName == Scenes_Recharge {
		if tx.Address.String() != OWNER_Address {
			return errors.New("权限不足")
		}
	}

	if tx.FuncName == Scenes_RegisterScenes {
		if scenes != nil {
			return errors.New("场景已注册 " + tx.ScenesID)
		}

	} else {
		if scenes == nil {
			return errors.New("场景未注册 " + tx.ScenesID)
		}
	}
	return nil
}

//执行业务逻辑,
// crossTxQcp: 需要进行跨链处理的TxQcp。
// 业务端实现中crossTxQcp只需包含`to` 和 `txStd`
func (tx *ScenesTx) Exec(ctx context.Context) (result types.Result, crossTxQcp *txs.TxQcp) {
	scenesMapper := ctx.Mapper(quzhuan.ScenesMapperName).(*quzhuan.ScenesMapper)
	switch tx.FuncName {
	case Scenes_RegisterScenes: //场景注册
		scenes := quzhuan.Scenes{ID: tx.ScenesID, Address: tx.ScenesAddress, Status: tx.Status}
		scenesMapper.AddScenes(scenes)
	case Scenes_ScenesReward: //活动奖励
		userMapper := ctx.Mapper(quzhuan.UsersMapperName).(*quzhuan.UsersMapper)
		scenesMapper.AddReward(tx.ScenesID, tx.Rewards)
		userMapper.AddReward(tx.Rewards)
		total := types.NewInt(0)
		//用户加余额
		coinsMapper := ctx.Mapper(quzhuan.CoinsMapperName).(*quzhuan.CoinsMapper)
		for _, v := range tx.Rewards {
			coinsMapper.UserAddBalance(v.UserID, v.Amount)
			total = total.Add(v.Amount)
		}
		//场景减余额
		coinsMapper.ScenesSubtractBalance(tx.ScenesID, total)

	case Scenes_Recharge: //场景充值
		coinsMapper := ctx.Mapper(quzhuan.CoinsMapperName).(*quzhuan.CoinsMapper)
		coinsMapper.ScenesAddBalance(tx.ScenesID, tx.Amount)
	}
	result = types.Result{
		Code: types.CodeOK,
	}
	return
}

func (tx *ScenesTx) NewTx(funcName string, args []string, address types.Address) error {
	tx.FuncName = funcName
	tx.Address = address
	args_len := len(args)

	switch funcName {
	case Scenes_RegisterScenes:
		if args_len != 3 {
			return errors.New(funcName + " ScenesTx args len error want " + strconv.Itoa(para_len_3) + " got " + strconv.Itoa(args_len))
		}
		tx.ScenesID = args[0]
		tx.ScenesAddress = args[1]
		status, err := strconv.Atoi(args[3])
		if err != nil {
			return err
		}
		tx.Status = status
	case Scenes_ScenesReward:
		if args_len != para_len_1 {
			return errors.New(funcName + " ScenesTx args len error want " + strconv.Itoa(para_len_1) + " got " + strconv.Itoa(args_len))
		}
		jsonstr := args[0]
		var sr ScenesReward
		err := json.Unmarshal([]byte(jsonstr), &sr)
		if err != nil {
			return err
		}
		tx.ScenesID = sr.ScenesId

		rewards := make([]*quzhuan.Reward, len(sr.Rewards))
		for i, v := range sr.Rewards {
			rewards[i] = &quzhuan.Reward{ScenesId: sr.ScenesId, Amount: v.Amount, UserID: v.UserId}
		}

		tx.Rewards = rewards

	case Scenes_Recharge:
		if args_len != para_len_2 {
			return errors.New(funcName + " ScenesTx args len error want " + strconv.Itoa(para_len_2) + " got " + strconv.Itoa(args_len))
		}
		tx.ScenesID = args[0]

		amount, ok := types.NewIntFromString(args[1])
		if !ok {
			return errors.New(funcName + " ScenesTx recharge amount format error")
		}
		tx.Amount = amount

	default:
		return errors.New(funcName + " funcName not support")
	}

	return nil

}
