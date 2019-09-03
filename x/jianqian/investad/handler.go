// Copyright 2018 The QOS Authors

package investad

import (
	"errors"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	qbasetypes "github.com/QOSGroup/qbase/types"
	qostxs "github.com/QOSGroup/qos/module/bank/txs"
	"github.com/QOSGroup/qstars/x/common"
	"log"

	qostypes "github.com/QOSGroup/qos/types"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/x/jianqian"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmcommon "github.com/tendermint/tendermint/libs/common"
	"strconv"
)

type InvestTx struct {
	Std         *txs.TxStd
	ArticleHash []byte `json:"articleHash"` // 文章hash
}

var _ txs.ITx = (*InvestTx)(nil)

func checkArticle(ctx context.Context, articleKey []byte) error {
	articleMapper := ctx.Mapper(jianqian.ArticlesMapperName).(*jianqian.ArticlesMapper)
	a := articleMapper.GetArticle(string(articleKey))
	if a == nil {
		return errors.New("invalid article")
	}

	log.Printf("--- checkArticle: EndInvestDate:%+v, Time:%+v", a.EndInvestDate, ctx.BlockHeader().Time)
	if a.EndInvestDate.Before(ctx.BlockHeader().Time) {
		return errors.New("超过投资期限")
	}

	return nil
}

func (it InvestTx) ValidateData(ctx context.Context) error {
	if err := checkArticle(ctx, it.ArticleHash); err != nil {
		return err
	}

	transferTx, ok := it.Std.ITxs[0].(*qostxs.TxTransfer)
	if !ok {
		return errors.New("std类型不支持")
	}

	if len(transferTx.Senders) == 0 || len(transferTx.Receivers) == 0 {
		return errors.New("无效的tx")
	}

	totalAmount := qbasetypes.NewInt(0)

	for _, v := range transferTx.Senders {
		totalAmount = totalAmount.Add(getInvestAmount(v.QSCs))
	}

	if totalAmount.IsZero() {
		return errors.New("投资金额不能为0")
	}

	return nil
}

func getInvestAmount(qscs qostypes.QSCs) qbasetypes.BigInt {
	amount := qbasetypes.NewInt(0)
	for _, v := range qscs {
		if v.Name == coinsName {
			amount = amount.Add(v.Amount)
		}
	}

	return amount
}

func (it InvestTx) Exec(ctx context.Context) (result qbasetypes.Result, crossTxQcps *txs.TxQcp) {
	result = qbasetypes.Result{
		Code: qbasetypes.CodeOK,
	}
	//set for qos result
	investUncheckedMapper := ctx.Mapper(jianqian.InvestUncheckedMapperName).(*jianqian.InvestUncheckedMapper)
	heigth1 := strconv.FormatInt(ctx.BlockHeight(), 10)
	tx1 := (tmcommon.HexBytes)(tmhash.Sum(ctx.TxBytes()))
	key := "heigth:" + heigth1 + ",hash:" + tx1.String()

	transferTx, _ := it.Std.ITxs[0].(*qostxs.TxTransfer)
	var values jianqian.InvestUncheckeds
	for _, v := range transferTx.Senders {
		values = append(values, jianqian.InvestUnchecked{
			Article:    it.ArticleHash,
			Address:    v.Address,
			InvestTime: ctx.BlockHeader().Time,
			Invest:     getInvestAmount(v.QSCs),
			IsChecked:  false,
		})
	}
	investUncheckedMapper.Set([]byte(key), values)

	kvMapper := ctx.Mapper(common.QSCResultMapperName).(*common.KvMapper)
	value := InvestadStub{}.Name()
	kvMapper.Set([]byte(key), value)
	log.Printf("investad.handler kvMapper key:[%s], value:[%s]", key, value)

	crossTxQcps = &txs.TxQcp{}
	crossTxQcps.TxStd = it.Std
	crossTxQcps.To = config.GetServerConf().QOSChainName
	crossTxQcps.Extends = key

	return
}

func (it InvestTx) GetSigner() []qbasetypes.Address {
	return it.Std.ITxs[0].GetSigner()
}

func (it InvestTx) CalcGas() qbasetypes.BigInt {
	return it.Std.ITxs[0].CalcGas()
}

func (it InvestTx) GetGasPayer() qbasetypes.Address {
	return it.Std.ITxs[0].GetGasPayer()
}

func (it InvestTx) GetSignData() []byte {
	sd := it.Std.ITxs[0].GetSignData()
	return append(sd, it.ArticleHash...)
}

func (it InvestTx) Name() string {
	return "InvestTx"
}
