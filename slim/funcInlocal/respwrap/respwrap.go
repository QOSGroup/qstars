package respwrap

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/tendermint/go-amino"
)

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

type RPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      string          `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
}

func ResponseWrapper(cdc *amino.Codec, result interface{}, err error) ([]byte, error) {
	if err != nil {
		return writeResponse(NewRPCErrorResponse("", 0, err.Error(), ""))
	} else {
		return writeResponse(NewRPCSuccessResponse(cdc, "", result))
	}
}

func writeResponse(res RPCResponse) ([]byte, error) {
	jsonBytes, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

func NewRPCErrorResponse(id string, code int, msg string, data string) RPCResponse {
	return RPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error:   &RPCError{Code: code, Message: msg, Data: data},
	}
}

func NewRPCSuccessResponse(cdc *amino.Codec, id string, res interface{}) RPCResponse {
	var rawMsg json.RawMessage

	if res != nil {
		var js []byte
		js, err := cdc.MarshalJSON(res)
		if err != nil {
			return RPCInternalError(id, errors.Wrap(err, "Error marshalling response"))
		}
		rawMsg = json.RawMessage(js)
	}

	return RPCResponse{JSONRPC: "2.0", ID: id, Result: rawMsg}
}

func RPCInternalError(id string, err error) RPCResponse {
	return NewRPCErrorResponse(id, -32603, "Internal error", err.Error())
}
