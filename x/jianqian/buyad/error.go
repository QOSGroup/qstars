// Copyright 2018 The QOS Authors

package buyad

type BuyadErr struct {
	code string
	err  string
}

func (be BuyadErr) Error() string {
	return be.err
}

func (be BuyadErr) Code() string {
	return be.code
}

var (
	InvalidArticleErr = &BuyadErr{code: InvalidArticleErrCode, err: "文章hash不正确"}
	CoinsErr          = &BuyadErr{code: CoinsErrCode, err: "coins不正确"}
	NoCommunityErr    = &BuyadErr{code: NoCommunityErrCode, err: "no community"}
	HasBeenBuyedErr   = &BuyadErr{code: HasBeenBuyedErrCode, err: "已经被购买"}
)

const (
	InvalidArticleErrCode = "401"
	CoinsErrCode          = "402"
	NoCommunityErrCode    = "403"
	HasBeenBuyedErrCode   = "404"
	SENDTXERRCode         = "410"
)

func NewBuyadErr(code, err string) *BuyadErr {
	return &BuyadErr{
		code: code,
		err:  err,
	}
}
