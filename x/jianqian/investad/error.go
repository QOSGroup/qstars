// Copyright 2018 The QOS Authors

package investad

type InvestadErr struct {
	code string
	err  string
}

func (ie InvestadErr) Error() string {
	return ie.err
}

func (ie InvestadErr) Code() string {
	return ie.code
}

var (
	InvalidArticleErr = &InvestadErr{code: InvalidArticleErrCode, err: "文章hash不正确"}
	CoinsErr          = &InvestadErr{code: CoinsErrCode, err: "coins不正确"}
)

const (
	InvalidArticleErrCode = "301"
	CoinsErrCode          = "302"
	SENDTXERRCode         = "310"
)

func NewInvestadErr(code, err string) *InvestadErr {
	return &InvestadErr{
		code: code,
		err:  err,
	}
}
