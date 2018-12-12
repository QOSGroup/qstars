// Copyright 2018 The QOS Authors

package common

import (
	"encoding/json"
	"github.com/QOSGroup/qstars/wire"
)

type Result struct {
	Code   string          `json:"code"`
	Height int64           `json:"height"`
	Hash   string          `json:"hash,omitempty"`
	Reason string          `json:"reason,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
}

func InternalError(reason string) Result {
	return NewErrorResult("500", reason)
}

func NewSuccessResult(cdc *wire.Codec, height int64, hash, res interface{}) Result {
	var rawMsg json.RawMessage

	if res != nil {
		var js []byte
		js, err := cdc.MarshalJSON(res)
		if err != nil {
			return InternalError(err.Error())
		}
		rawMsg = json.RawMessage(js)
	}

	var result Result
	result.Result = rawMsg
	result.Code = "0"

	return result
}

func NewErrorResult(code, reason string) Result {
	var result Result
	result.Code = code
	result.Reason = reason

	return result
}

func (r Result) Marshal() string {
	jsonBytes, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return InternalError(err.Error()).Marshal()
	}
	return string(jsonBytes)
}

// ResultQstarsTimeoutError 联盟链超时
func ResultQstarsTimeoutError() Result {
	return NewErrorResult("-2", "qstars timeout")
}

// ResultQOSTimeoutError 主链超时
func ResultQOSTimeoutError() Result {
	return NewErrorResult("-1", "qos timeout")
}
