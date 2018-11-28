// Copyright 2018 The QOS Authors

package buyad

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	qbasetypes "github.com/QOSGroup/qbase/types"
	qostxs "github.com/QOSGroup/qos/txs"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/x/jianqian"
	"time"
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
	now := time.Now()
	transferTx, _ := it.Std.ITx.(*qostxs.TransferTx)
	if len(transferTx.Senders) != 1 {
		result.Code = qbasetypes.ToABCICode(qbasetypes.CodespaceRoot, qbasetypes.CodeInternal)
		return result, nil
	}

	buyer, ok := buyMapper.GetBuyer(it.ArticleHash)
	if ok {
		buyMapper.SetBuyer(it.ArticleHash, buyer)
	} else {
		buyer = jianqian.Buyer{
			Address: transferTx.Senders[0].Address,
			Buy:     transferTx.Senders[0].QOS,
			BuyTime: now,
		}
		buyMapper.SetBuyer(it.ArticleHash, buyer)
	}

	crossTxQcps = &txs.TxQcp{}
	crossTxQcps.TxStd = it.Std
	crossTxQcps.To = config.GetServerConf().QOSChainName

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
