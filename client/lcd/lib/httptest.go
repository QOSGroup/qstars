// Copyright 2018 The QOS Authors

package lib

import (
	"encoding/json"
	amino "github.com/tendermint/go-amino"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	tmltypes "github.com/tendermint/tendermint/rpc/lib/types"
)

const (
	TestInternalUser         = "test1@internal.local"
	TestInternalUserPassword = "123456"
)

type HttpTest struct {
	t *testing.T
	h http.Handler

	req *http.Request
}

func NewHttpTest(t *testing.T, req *http.Request) *HttpTest {
	var ht HttpTest
	ht.t = t
	ht.req = req

	return &ht
}

func (ht *HttpTest) Do(r http.Handler, cdc *amino.Codec, v interface{}) (*httptest.ResponseRecorder, error) {
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, ht.req)

	if w, ok := v.(io.Writer); ok {
		io.Copy(w, rw.Body)
	} else {
		var tmresp tmltypes.RPCResponse
		err := json.NewDecoder(rw.Body).Decode(&tmresp)
		if err != nil {
			if err == io.EOF {
				err = nil // ignore EOF errors caused by empty response body
			} else {
				return rw, err
			}
		}

		if tmresp.Error != nil {
			return rw, tmresp.Error
		}

		err = cdc.UnmarshalJSON(tmresp.Result, v)
		if err != nil {
			return rw, err
		}
	}

	return rw, nil
}
