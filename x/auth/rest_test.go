// Copyright 2018 The QOS Authors

package auth

import (
	qosaccount "github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qstars/client/lcd/lib"
	"github.com/QOSGroup/qstars/star"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"net/http"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()
}

func exampleTestRest(t *testing.T) {
	cdc := star.MakeCodec()
	req, err := http.NewRequest(http.MethodGet, "/accounts/"+"address1k0m8ucnqug974maa6g36zw7g2wvfd4sug6uxay", nil)
	assert.Nil(t, err)

	r := mux.NewRouter()
	RegisterRoutes(cdc, r)

	var res *qosaccount.QOSAccount
	_, err = lib.NewHttpTest(t, req).Do(r, cdc, res)
	assert.Nil(t, err)

	t.Logf("--%+v", res)
}
