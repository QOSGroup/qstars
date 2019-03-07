// Copyright 2018 The QOS Authors

package tx

import (
	"github.com/QOSGroup/qbase/txs"
	qbasetypes "github.com/QOSGroup/qbase/types"
	approve "github.com/QOSGroup/qos/module/approve"
	approvetype "github.com/QOSGroup/qos/module/approve/types"
	"github.com/QOSGroup/qos/module/approve/types"
)

type ApproveTx struct {
	from qbasetypes.Address
	to   qbasetypes.Address
}

func NewApproveTx(from qbasetypes.Address, to qbasetypes.Address) *ApproveTx {
	return &ApproveTx{
		from: from,
		to:   to,
	}
}

func (at *ApproveTx) Create(coins []qbasetypes.BaseCoin) txs.ITx {
	ti := warpperTransItem(nil, coins)

	return approve.TxCreateApprove{
		Approve: types.Approve{
			From: at.from,
			To:   at.to,
			QOS:  ti.QOS,
			QSCs: ti.QSCs,
		},
	}
}

func (at *ApproveTx) Increase(coins []qbasetypes.BaseCoin) txs.ITx {
	ti := warpperTransItem(nil, coins)

	return approve.TxIncreaseApprove{
		Approve: approvetype.Approve{
			From: at.from,
			To:   at.to,
			QOS:  ti.QOS,
			QSCs: ti.QSCs,
		},
	}
}

func (at *ApproveTx) Decrease(coins []qbasetypes.BaseCoin) txs.ITx {
	ti := warpperTransItem(nil, coins)

	return approve.TxDecreaseApprove{
		Approve: approvetype.Approve{
			From: at.from,
			To:   at.to,
			QOS:  ti.QOS,
			QSCs: ti.QSCs,
		},
	}
}

func (at *ApproveTx) Use(coins []qbasetypes.BaseCoin) txs.ITx {
	ti := warpperTransItem(nil, coins)

	return approve.TxUseApprove{
		Approve: approvetype.Approve{
			From: at.from,
			To:   at.to,
			QOS:  ti.QOS,
			QSCs: ti.QSCs,
		},
	}
}

func (at *ApproveTx) Cancel() txs.ITx {
	return approve.TxCancelApprove{
		From: at.from,
		To:   at.to,
	}
}
