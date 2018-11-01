// Copyright 2018 The QOS Authors

// Package pkg comments for pkg lib
// lib ...
package lib

import (
	"encoding/json"
	"io"
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

func ResponseWrapper(w io.Writer, cdc *amino.Codec, result interface{}, err error) {
	if err != nil {
		writeResponse(w, tmtype.NewRPCErrorResponse("", 0, err.Error(), ""))
	} else {
		writeResponse(w, tmtype.NewRPCSuccessResponse(cdc, "", result))
	}
}

func writeResponse(w io.Writer, res tmtype.RPCResponse) {
	jsonBytes, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		panic(err)
	}
	w.Write(jsonBytes)
}
