// Copyright 2018 The QOS Authors

package investad

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	qbasetypes "github.com/QOSGroup/qbase/types"
	qostxs "github.com/QOSGroup/qos/txs"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/x/jianqian"
	"time"
)

type InvestTx struct {
	//transferTx  qostxs.TransferTx
	Std         *txs.TxStd
	ArticleHash []byte `json:"articleHash"` // 文章hash
}

var _ txs.ITx = (*InvestTx)(nil)

func checkArticle(article []byte) error {
	return nil
}

func (it InvestTx) ValidateData(ctx context.Context) error {
	if err := checkArticle(it.ArticleHash); err != nil {
		return err
	}

	return it.Std.ITx.ValidateData(ctx)
}

func (it InvestTx) Exec(ctx context.Context) (result qbasetypes.Result, crossTxQcps *txs.TxQcp) {
	result = qbasetypes.Result{
		Code: qbasetypes.ABCICodeOK,
	}
	//set for qos result
	investMapper := ctx.Mapper(jianqian.InvestMapperName).(*jianqian.InvestMapper)
	now := time.Now()
	transferTx, _ := it.Std.ITx.(*qostxs.TransferTx)
	for _, v := range transferTx.Senders {
		investor, ok := investMapper.GetInvestor(it.ArticleHash, v.Address)
		if ok {
			investMapper.SetInvestor(it.ArticleHash, investor)
		} else {
			investor = jianqian.Investor{
				Address:    v.Address,
				InvestTime: now,
			}
			investMapper.SetInvestor(it.ArticleHash, investor)
		}
	}

	crossTxQcps = &txs.TxQcp{}
	crossTxQcps.TxStd = it.Std
	crossTxQcps.To = config.GetServerConf().QOSChainName

	return
}

func (it InvestTx) GetSigner() []qbasetypes.Address {
	return it.Std.ITx.GetSigner()
}

func (it InvestTx) CalcGas() qbasetypes.BigInt {
	return it.Std.ITx.CalcGas()
}

func (it InvestTx) GetGasPayer() qbasetypes.Address {
	return it.Std.ITx.GetGasPayer()
}

func (it InvestTx) GetSignData() []byte {
	sd := it.Std.ITx.GetSignData()

	return append(sd, it.ArticleHash...)
}
