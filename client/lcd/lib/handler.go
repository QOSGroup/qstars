// Copyright 2018 The QOS Authors

// Package pkg comments for pkg lib
// lib ...
package lib

import (
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
