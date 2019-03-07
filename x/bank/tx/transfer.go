// Copyright 2018 The QOS Authors

package tx

import (
	"github.com/QOSGroup/qbase/txs"
	qbasetypes "github.com/QOSGroup/qbase/types"
	qostxtype "github.com/QOSGroup/qos/module/transfer/types"
	qostxs "github.com/QOSGroup/qos/module/transfer"
)

func warpperTransItem(addr qbasetypes.Address, coins []qbasetypes.BaseCoin) qostxtype.TransItem {

	var ti qostxtype.TransItem
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
func NewTransfer(sender qbasetypes.Address, receiver qbasetypes.Address, coin []qbasetypes.BaseCoin) txs.ITx {
	var sendTx qostxs.TxTransfer

	sendTx.Senders = append(sendTx.Senders, warpperTransItem(sender, coin))
	sendTx.Receivers = append(sendTx.Receivers, warpperTransItem(receiver, coin))

	return sendTx
}

// NewTransfer ...
func NewTransferMultiple(sender []qbasetypes.Address, receiver []qbasetypes.Address, sendercoins [][]qbasetypes.BaseCoin,receivercoins [][]qbasetypes.BaseCoin) txs.ITx {
	var sendTx qostxs.TxTransfer

	for i:= 0;i<len(sender);i++{
		sendTx.Senders = append(sendTx.Senders, warpperTransItem(sender[i], sendercoins[i]))
	}
	for i:= 0;i<len(receiver);i++{
		sendTx.Receivers = append(sendTx.Receivers, warpperTransItem(receiver[i], receivercoins[i]))
	}

	return sendTx
}