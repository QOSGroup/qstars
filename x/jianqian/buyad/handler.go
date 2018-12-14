// Copyright 2018 The QOS Authors

package buyad

import (
	"errors"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	qbasetypes "github.com/QOSGroup/qbase/types"
	qostxs "github.com/QOSGroup/qos/txs/transfer"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/x/common"
	"github.com/QOSGroup/qstars/x/jianqian"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmcommon "github.com/tendermint/tendermint/libs/common"
	"log"
	"time"

	"strconv"
)

type BuyTx struct {
	Std         *txs.TxStd
	ArticleHash []byte `json:"articleHash"` // 文章hash
}

var _ txs.ITx = (*BuyTx)(nil)

func checkArticleBase(article *jianqian.Articles, now time.Time) error {
	if article.EndBuyDate.After(now) {
		return errors.New("投资还没结束期限")

	}
	if article.EndBuyDate.Before(now) {
		return errors.New("超过购买期限")
	}

	return nil
}

func check(ctx context.Context, articleKey []byte) error {
	articleMapper := ctx.Mapper(jianqian.ArticlesMapperName).(*jianqian.ArticlesMapper)
	a := articleMapper.GetArticle(string(articleKey))
	if a == nil {
		return errors.New("invalid article")
	}

	log.Printf("--- checkArticle: EndBuyDate:%+v, Time:%+v", a.EndBuyDate, ctx.BlockHeader().Time)
	if err := checkArticleBase(a, ctx.BlockHeader().Time); err != nil {
		return err
	}

	return nil
}

func (it BuyTx) ValidateData(ctx context.Context) error {
	if err := check(ctx, it.ArticleHash); err != nil {
		return err
	}

	buyMapper := ctx.Mapper(jianqian.BuyMapperName).(*jianqian.BuyMapper)
	buyer, ok := buyMapper.GetBuyer(it.ArticleHash)
	if ok && buyer.CheckStatus != jianqian.CheckStatusFail {
		return errors.New("已被购买")
	}

	transferTx, ok := it.Std.ITx.(*qostxs.TxTransfer)
	if !ok {
		return errors.New("std类型不支持")
	}

	if len(transferTx.Senders) == 0 || len(transferTx.Receivers) == 0 {
		return errors.New("无效的tx")
	}

	totalAmount := qbasetypes.NewInt(0)
	for _, v := range transferTx.Senders {
		totalAmount = totalAmount.Add(v.QOS)
	}
	if totalAmount.IsZero() {
		return errors.New("购买金额不能为0")
	}

	return nil
}

func (it BuyTx) Exec(ctx context.Context) (result qbasetypes.Result, crossTxQcps *txs.TxQcp) {
	log.Printf("buyad.handler Exec")

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
	if !ok {
		buyer = &jianqian.Buyer{}
	}

	buyer.ArticleHash = it.ArticleHash
	buyer.Address = buyerAddr
	buyer.Buy = qos
	buyer.BuyTime = ctx.BlockHeader().Time
	buyer.CheckStatus = jianqian.CheckStatusInit
	buyMapper.SetBuyer(it.ArticleHash, *buyer)

	heigth1 := strconv.FormatInt(ctx.BlockHeight(), 10)
	tx1 := (tmcommon.HexBytes)(tmhash.Sum(ctx.TxBytes()))
	key := "heigth:" + heigth1 + ",hash:" + tx1.String()
	kvMapper := ctx.Mapper(common.QSCResultMapperName).(*common.KvMapper)
	kvMapper.Set([]byte(key), BuyadStub{}.Name())

	// 设置临时状态,便于qos返回信息处理
	buyMapper.SetBuyer([]byte(key), *buyer)

	crossTxQcps = &txs.TxQcp{}
	crossTxQcps.TxStd = it.Std
	crossTxQcps.To = config.GetServerConf().QOSChainName
	crossTxQcps.Extends = key

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

func (it BuyTx) Name() string {
	return "BuyTx"
}
