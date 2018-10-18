package baseapp

import (
	"testing"
	"github.com/QOSGroup/qstars/x/kvstore"
)

// TODO update
func TestInitCmd(t *testing.T) {
	app := QstarsBaseApp{}
	mock := kvstore.KVStub{}
	app.Register(mock)
	app.Start()

}

type MockBaseContract struct {

}

func (mb MockBaseContract ) A(){

}