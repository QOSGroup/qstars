// Copyright 2018 The QOS Authors

package common

import (
	"encoding/json"
	"github.com/QOSGroup/qstars/wire"
)

const (
	ResultCodeSuccess       = "0"
	ResultCodeQstarsTimeout = "-2"
	ResultCodeQOSTimeout    = "-1"
	ResultCodeInternalError = "500"
)

type Result struct {
	Code   string          `json:"code"`
	Height int64           `json:"height"`
	Hash   string          `json:"hash,omitempty"`
	Reason string          `json:"reason,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
}

func InternalError(reason string) Result {
	return NewErrorResult(ResultCodeInternalError, 0, "", reason)
}

func NewSuccessResult(cdc *wire.Codec, height int64, hash string, res interface{}) Result {
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
	result.Height = height
	result.Hash = hash
	result.Result = rawMsg
	result.Code = "0"

	return result
}

func NewErrorResult(code string, height int64, hash string, reason string) Result {
	var result Result
	result.Height = height
	result.Hash = hash
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
func ResultQstarsTimeoutError(height int64, hash string) Result {
	return NewErrorResult(ResultCodeQstarsTimeout, height, hash, "qstars timeout")
}

// ResultQOSTimeoutError 主链超时
func ResultQOSTimeoutError(height int64, hash string) Result {
	return NewErrorResult(ResultCodeQOSTimeout, height, hash, "qos timeout")
}
