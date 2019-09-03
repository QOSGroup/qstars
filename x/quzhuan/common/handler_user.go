package common

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/x/quzhuan"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

const (
	//注册用户
	User_RegisterUser  = "registerUser"
	//上传用户地址（充值提现使用）
	User_UploadAddress = "uploadAddress"
	//上传用户kycid
	User_AddKYCID      = "addKYCID"
)

type UserTx struct {
	Address     types.Address //签名地址
	ID          string        //会员唯一标识
	Source      string        //初次注册场景ID
	Inviter     string        //邀请人
	UserAddress string        //会员充值地址
	KYCID       string        //KYC id
	FuncName    string        //方法名
}

func (tx *UserTx) String() string {
	return tx.ID + "|" + tx.Source + "|" + tx.UserAddress + "|" + tx.Inviter + "|" + tx.KYCID
}

var _ RouterTx = (*UserTx)(nil)

func (tx *UserTx) ValidateData(ctx context.Context) error {
	if strings.TrimSpace(tx.ID) == "" {
		return errors.New("用户ID 不能为空")
	}
	userMapper := ctx.Mapper(quzhuan.UsersMapperName).(*quzhuan.UsersMapper)
	user := userMapper.GetUser(tx.ID)
	if tx.FuncName == User_RegisterUser {
		if user != nil {
			return errors.New("用户已注册 " + tx.ID)
		}
		scenesMapper := ctx.Mapper(quzhuan.ScenesMapperName).(*quzhuan.ScenesMapper)
		scenesId := tx.Source
		scenes := scenesMapper.GetScenesByID(scenesId)
		if scenes != nil {
			if scenes.Address != tx.Address.String() {
				return errors.New("场景来源与备案不符")
			}
		} else {
			return errors.New("场景来源未备案")
		}
	} else {
		if user == nil {
			return errors.New("用户未注册 " + tx.ID)
		}
	}
	return nil
}

//执行业务逻辑,
// crossTxQcp: 需要进行跨链处理的TxQcp。
// 业务端实现中crossTxQcp只需包含`to` 和 `txStd`
func (tx *UserTx) Exec(ctx context.Context) (result types.Result, crossTxQcp *txs.TxQcp) {
	userMapper := ctx.Mapper(quzhuan.UsersMapperName).(*quzhuan.UsersMapper)

	switch tx.FuncName {
	case User_RegisterUser:
		user := quzhuan.User{ID: tx.ID, Source: tx.Source, Inviter: tx.Inviter}
		userMapper.AddUser(user)
	case User_AddKYCID:
		//检查是否已经绑定过其他用户
		olduserid := userMapper.KYCExist(tx.KYCID)
		if olduserid != "" {
			userMapper.RemoveKYC(olduserid)
		}
		userMapper.AddKYCID(tx.ID, tx.KYCID)
	case User_UploadAddress:
		userMapper.AddAddress(tx.ID, tx.UserAddress)
	}
	result = types.Result{
		Code: types.CodeOK,
	}
	return

}

func (tx *UserTx) NewTx(funcName string, args []string, address types.Address) error {
	tx.FuncName = funcName
	tx.Address = address
	args_len := len(args)

	switch funcName {
	case User_RegisterUser:
		if args_len != para_len_3 {
			return errors.New(funcName + " UserTx args len error want " + strconv.Itoa(para_len_3) + " got " + strconv.Itoa(args_len))
		}
		tx.ID = args[0]
		tx.Source = args[1]
		tx.Inviter = args[2]
	case User_AddKYCID:
		if args_len != para_len_2 {
			return errors.New(funcName + " UserTx args len error want " + strconv.Itoa(para_len_2) + " got " + strconv.Itoa(args_len))
		}
		tx.ID = args[0]
		tx.KYCID = args[1]

	case User_UploadAddress:
		if args_len != para_len_2 {
			return errors.New(funcName + " UserTx args len error want " + strconv.Itoa(para_len_2) + " got " + strconv.Itoa(args_len))
		}
		tx.ID = args[0]
		tx.UserAddress = args[1]

	default:
		return errors.New(funcName + " funcName not support")
	}

	return nil

}
