// Copyright 2018 The QOS Authors

package jianqian

import (
	"github.com/QOSGroup/qstars/x/common"
)

// article 101-199
// buyad 201-299
// coins 301-399
// invested 401-499

// ResultArticleNotFoundError 文章不存在
func ResultArticleNotFoundError() common.Result {
	return common.NewErrorResult("101", "article not found")
}

// ResultBuyadSaledError 已售
func ResultBuyadSaledError() common.Result {
	return common.NewErrorResult("201", "has been saled")
}

// ResultBuyadExpiredError 购买过期
func ResultBuyadExpiredError() common.Result {
	return common.NewErrorResult("202", "buy expire")
}

// ResultBuyadNotStartError 还未到购买时间
func ResultBuyadNotStartError() common.Result {
	return common.NewErrorResult("203", "buy not start")
}

// ResultInvestadExpiredError 投资过期
func ResultInvestadExpiredError() common.Result {
	return common.NewErrorResult("401", "invest expire")
}
