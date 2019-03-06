// Copyright 2018 The QOS Authors

package investad

import (
	"errors"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	qbasetypes "github.com/QOSGroup/qbase/types"
	qostypes "github.com/QOSGroup/qos/types"
	"github.com/QOSGroup/qstars/x/jianqian"
	"log"
	"strings"
)

type InvestTx struct {
	Address      qbasetypes.Address `json:"address"`      // 投资者地址
	OtherAddr    string             `json:"otherAddr"`      // 投资者其他地址
	Invest       qbasetypes.BigInt  `json:"investad"`     // 投资金额
	ArticleHash  []byte             `json:"articleHash"`  // 文章hash
	Gas          qbasetypes.BigInt

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
	aoeaccount := ctx.Mapper(jianqian.AoeAccountMapperName).(*jianqian.AoeAccountMapper)
	blance:=aoeaccount.GetBalance(it.Address.String())
	if blance.IsZero()||blance.IsNil()||blance.Int64()<0{
		return errors.New("投资者余额不足")
	}

	if strings.TrimSpace(it.Address.String()) == "" {
		return errors.New("投资者地址不能为空")
	}
	if it.Invest.IsZero() {
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
	investMapper := ctx.Mapper(jianqian.InvestMapperName).(*jianqian.InvestMapper)
	key := jianqian.GetInvestKey(it.ArticleHash, it.OtherAddr, jianqian.InvestorTypeCommonInvestor)
	investor, ok := investMapper.GetInvestor(key)
	if ok {
		investor.Invest = investor.Invest.Add(it.Invest)
		investor.InvestTime = ctx.BlockHeader().Time
		log.Printf("investad.InvestadStub investor update %+v\n", investor)
		investMapper.SetInvestor(key, investor)
	} else {
		investor = jianqian.Investor{
			InvestorType: jianqian.InvestorTypeCommonInvestor,
			//Address:      it.Address,
			OtherAddr:    it.OtherAddr,
			InvestTime:   ctx.BlockHeader().Time,
			Invest:       it.Invest,
		}
		log.Printf("investad.InvestadStub investor create %+v\n", investor)
		investMapper.SetInvestor(key, investor)
	}
    //账户余额中减掉投资金额
	aoeaccount := ctx.Mapper(jianqian.AoeAccountMapperName).(*jianqian.AoeAccountMapper)
	aoeaccount.SubtractBalance(it.Address,it.Invest)

	return
}

func (it InvestTx) GetSigner() []qbasetypes.Address {
	return []qbasetypes.Address{it.Address}
}

func (it InvestTx) CalcGas() qbasetypes.BigInt {
	return it.Gas
}

func (it InvestTx) GetGasPayer() qbasetypes.Address {
	return it.Address
}

func (it InvestTx) GetSignData() (ret []byte) {
	ret = append(ret, it.ArticleHash...)
	ret = append(ret, it.Address.Bytes()...)
	ret = append(ret, qbasetypes.Int2Byte(it.Invest.Int64())...)
	ret = append(ret, it.OtherAddr...)
	return
}

func (it InvestTx) Name() string {
	return "InvestTx"
}
