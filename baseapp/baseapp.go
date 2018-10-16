package baseapp

import (
	"fmt"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/tendermint/tendermint/libs/log"
	"os"
)

func NewAPP(rootDir string) QstarsBaseApp{
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "main")
	qstarts:=QstarsBaseApp{
		Logger:logger,
		RootDir:rootDir,
	}
	return qstarts
}


type QstarsBaseApp struct {
	Contracts BaseContract
	Baseapp *baseabci.BaseApp
	ContractList []BaseContract
	Logger log.Logger
	RootDir string
}

func (base *QstarsBaseApp) Register(basecontract BaseContract) {
	base.ContractList = append(base.ContractList,basecontract)
}

func (base *QstarsBaseApp) loadX(){
	for index, c := range base.ContractList {
		fmt.Printf("arr[%d]=%d \n", index, c)
		c.RegisterKVCdc(base.Baseapp.GetCdc())
		c.StartX(base)
	}
}


func (base *QstarsBaseApp) Start() {

	base.loadX()


	// Start the ABCI server
	//srv, err := server.NewServer("0.0.0.0:26658", "socket", base.Baseapp)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//err = srv.Start()
	//if err != nil {
	//	cmn.Exit(err.Error())
	//}

	// Wait forever
	//cmn.TrapSignal(func() {
	//	// Cleanup
	//	err = srv.Stop()
	//	if err != nil {
	//		cmn.Exit(err.Error())
	//	}
	//})
}