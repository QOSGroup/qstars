// Copyright 2018 The QOS Authors

package comm

import (
	"errors"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	qbasetypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/x/jianqian"
	"log"
	"strconv"
	"strings"
)

type InvestTx struct {
	Address     qbasetypes.Address `json:"address"`     // 投资者地址
	Invest      qbasetypes.BigInt  `json:"investad"`    // 投资金额
	ArticleHash []byte             `json:"articleHash"` // 文章hash
	cointype    string
}

var _ RouterTx = (*InvestTx)(nil)

func checkArticle(ctx context.Context, a *jianqian.Articles) error {
	if a == nil {
		return errors.New("invalid article")
	}
	log.Printf("--- checkArticle: EndInvestDate:%+v, Time:%+v", a.EndInvestDate, ctx.BlockHeader().Time)
	if a.EndInvestDate.Before(ctx.BlockHeader().Time) {
		return errors.New("超过投资期限")
	}
	return nil
}

func (it *InvestTx) ValidateData(ctx context.Context) error {

	articleMapper := ctx.Mapper(jianqian.ArticlesMapperName).(*jianqian.ArticlesMapper)
	a := articleMapper.GetArticle(string(it.ArticleHash))

	if err := checkArticle(ctx, a); err != nil {
		return err
	}

	it.cointype = a.CoinType
	aoeaccount := ctx.Mapper(jianqian.AoeAccountMapperName).(*jianqian.AoeAccountMapper)
	blance := aoeaccount.GetBalance(it.Address.String(), it.cointype)
	if blance.IsZero() || blance.IsNil() || blance.Int64() < 0 || it.Invest.GT(blance) {
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

func (it *InvestTx) Exec(ctx context.Context) (result qbasetypes.Result, crossTxQcps *txs.TxQcp) {
	investMapper := ctx.Mapper(jianqian.InvestMapperName).(*jianqian.InvestMapper)
	key := jianqian.GetInvestKey(it.ArticleHash, it.Address.String(), jianqian.InvestorTypeCommonInvestor)
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
			//OtherAddr:    it.OtherAddr,
			InvestTime: ctx.BlockHeader().Time,
			Invest:     it.Invest,
		}
		log.Printf("investad.InvestadStub investor create %+v\n", investor)
		investMapper.SetInvestor(key, investor)
	}
	//账户余额中减掉投资金额
	aoeaccount := ctx.Mapper(jianqian.AoeAccountMapperName).(*jianqian.AoeAccountMapper)
	aoeaccount.SubtractBalance(it.Address.String(), it.cointype, it.Invest)

	return
}

func (tx *InvestTx) NewTx(args []string) error {
	args_len := len(args)
	if args_len != para_len_4 {
		return errors.New("AdvertisersTx args len error want " + strconv.Itoa(para_len_4) + " got " + strconv.Itoa(args_len))
	}
	address, err := types.AccAddressFromBech32(args[0])
	if err != nil {
		return err
	}

	amount, ok := qbasetypes.NewIntFromString(args[1])
	if !ok {
		return errors.New("amount format error")
	}

	tx.Address = address
	tx.Invest = amount
	tx.ArticleHash = []byte(args[2])
	tx.cointype = args[3]
	return nil
}
