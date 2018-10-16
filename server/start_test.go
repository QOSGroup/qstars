package server

import (
	"testing"
)

func TestStartStandAlone(t *testing.T) {
	//home, err := ioutil.TempDir("", "mock-sdk-cmd")
	//require.Nil(t, err)
	//defer func() {
	//	os.RemoveAll(home)
	//}()
	//
	//logger := log.NewNopLogger()
	//cfg, err := tcmd.ParseConfig()
	//require.Nil(t, err)
	//ctx := NewContext(cfg, logger)
	//cdc := wire.NewCodec()
	//appInit := AppInit{
	//	AppGenState: mock.AppGenState,
	//	AppGenTx:    mock.AppGenTx,
	//}
	//initCmd := InitCmd(ctx, cdc, appInit)
	//err = initCmd.RunE(nil, nil)
	//require.NoError(t, err)
	//
	//app, err := mock.NewApp(home, logger)
	//require.Nil(t, err)
	//svrAddr, _, err := FreeTCPAddr()
	//require.Nil(t, err)
	//svr, err := server.NewServer(svrAddr, "socket", app)
	//require.Nil(t, err, "error creating listener")
	//svr.SetLogger(logger.With("module", "abci-server"))
	//svr.Start()
	//
	//timer := time.NewTimer(time.Duration(2) * time.Second)
	//select {
	//case <-timer.C:
	//	svr.Stop()
	//}
}
