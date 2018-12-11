// Copyright 2018 The QOS Authors

package tx

import (
	"github.com/QOSGroup/qbase/txs"
	qbasetypes "github.com/QOSGroup/qbase/types"
	qostxs "github.com/QOSGroup/qos/txs/transfer"
)

func warpperTransItem(addr qbasetypes.Address, coins []qbasetypes.BaseCoin) qostxs.TransItem {
	var ti qostxs.TransItem
	ti.Address = addr
	ti.QOS = qbasetypes.NewInt(0)

	for _, coin := range coins {
		if coin.Name == "qos" {
			ti.QOS = ti.QOS.Add(coin.Amount)
		} else {
			newcoin:=coin
			ti.QSCs = append(ti.QSCs, &newcoin)
		}
	}
	return ti
}

// NewTransfer ...
func NewTransfer(sender []qbasetypes.Address, receiver []qbasetypes.Address, coin []qbasetypes.BaseCoin) txs.ITx {
	var sendTx qostxs.TxTransfer
	for _,sv:=range sender{
		//一转多时 sender要提供转出总额
		name:=coin[0].Name
		var total qbasetypes.BigInt=qbasetypes.ZeroInt()
		for _,v:=range coin{
			total=total.Add(v.Amount)
		}
		newcoin:=qbasetypes.BaseCoin{name,total}
		sendTx.Senders = append(sendTx.Senders, warpperTransItem(sv, []qbasetypes.BaseCoin{newcoin}))
	}
	for i,rv:=range receiver {
		sendTx.Receivers = append(sendTx.Receivers, warpperTransItem(rv, []qbasetypes.BaseCoin{coin[i]}))
	}
	return sendTx
}
