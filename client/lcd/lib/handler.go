// Copyright 2018 The QOS Authors

// Package pkg comments for pkg lib
// lib ...
package lib

import (
	"encoding/json"
	"net/http"

	"github.com/tendermint/go-amino"
	tmserver "github.com/tendermint/tendermint/rpc/lib/server"
	tmtype "github.com/tendermint/tendermint/rpc/lib/types"
)

func HttpResponseWrapper(w http.ResponseWriter, cdc *amino.Codec, result interface{}, err error) {
	if err != nil {
		tmserver.WriteRPCResponseHTTP(w, tmtype.NewRPCErrorResponse("", 0, err.Error(), ""))
	} else {
		tmserver.WriteRPCResponseHTTP(w, tmtype.NewRPCSuccessResponse(cdc, "", result))
	}
}

func ResponseWrapper(cdc *amino.Codec, result interface{}, err error) ([]byte, error) {
	if err != nil {
		return writeResponse(tmtype.NewRPCErrorResponse("", 0, err.Error(), ""))
	} else {
		return writeResponse(tmtype.NewRPCSuccessResponse(cdc, "", result))
	}
}

func writeResponse(res tmtype.RPCResponse) ([]byte, error) {
	jsonBytes, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}
