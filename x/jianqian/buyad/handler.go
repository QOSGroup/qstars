// Copyright 2018 The QOS Authors

package buyad

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	qbasetypes "github.com/QOSGroup/qbase/types"
	qostxs "github.com/QOSGroup/qos/txs/transfer"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/x/common"
	"github.com/QOSGroup/qstars/x/jianqian"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmcommon "github.com/tendermint/tendermint/libs/common"

	"strconv"
)

type BuyTx struct {
	Std         *txs.TxStd
	ArticleHash []byte `json:"articleHash"` // 文章hash
}

var _ txs.ITx = (*BuyTx)(nil)

func checkArticle(article []byte) error {
	return nil
}

func (it BuyTx) ValidateData(ctx context.Context) error {
	if err := checkArticle(it.ArticleHash); err != nil {
		return err
	}

	return it.Std.ITx.ValidateData(ctx)
}

func (it BuyTx) Exec(ctx context.Context) (result qbasetypes.Result, crossTxQcps *txs.TxQcp) {
	result = qbasetypes.Result{
		Code: qbasetypes.ABCICodeOK,
	}
	//set for qos result
	buyMapper := ctx.Mapper(jianqian.BuyMapperName).(*jianqian.BuyMapper)
	transferTx, _ := it.Std.ITx.(*qostxs.TxTransfer)
	if len(transferTx.Senders) != 1 {
		result.Code = qbasetypes.ToABCICode(qbasetypes.CodespaceRoot, qbasetypes.CodeInternal)
		return result, nil
	}

	qos := transferTx.Senders[0].QOS
	buyerAddr := transferTx.Senders[0].Address
	buyer, ok := buyMapper.GetBuyer(it.ArticleHash)
	if ok {
		if buyer.CheckStatus != jianqian.CheckStatusFail {
			result.Code = qbasetypes.ToABCICode(qbasetypes.CodespaceRoot, qbasetypes.CodeInternal)
			return result, nil
		}
	} else {
		buyer = &jianqian.Buyer{}
	}

	buyer.Address = buyerAddr
	buyer.Buy = qos
	buyer.BuyTime = ctx.BlockHeader().Time
	buyer.CheckStatus = jianqian.CheckStatusInit
	buyMapper.SetBuyer(it.ArticleHash, *buyer)

	heigth1 := strconv.FormatInt(ctx.BlockHeight(), 10)
	tx1 := (tmcommon.HexBytes)(tmhash.Sum(ctx.TxBytes()))
	key := "heigth:" + heigth1 + ",hash:" + tx1.String()
	kvMapper := ctx.Mapper(common.QSCResultMapperName).(*common.KvMapper)
	kvMapper.Set([]byte(key), "-1")

	crossTxQcps = &txs.TxQcp{}
	crossTxQcps.TxStd = it.Std
	crossTxQcps.To = config.GetServerConf().QOSChainName
	crossTxQcps.Extends = string(it.ArticleHash)

	return
}

func (it BuyTx) GetSigner() []qbasetypes.Address {
	return it.Std.ITx.GetSigner()
}

func (it BuyTx) CalcGas() qbasetypes.BigInt {
	return it.Std.ITx.CalcGas()
}

func (it BuyTx) GetGasPayer() qbasetypes.Address {
	return it.Std.ITx.GetGasPayer()
}

func (it BuyTx) GetSignData() []byte {
	sd := it.Std.ITx.GetSignData()

	return append(sd, it.ArticleHash...)
}
