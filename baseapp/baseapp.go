package baseapp

import (
	"fmt"
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/tendermint/tendermint/libs/log"
	"os"
	"path/filepath"

	dbm "github.com/tendermint/tendermint/libs/db"
)

func NewAPP(rootDir string) QstarsBaseApp {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "main")
	qstarts := QstarsBaseApp{
		Logger:  logger,
		RootDir: rootDir,
	}
	return qstarts
}

type QstarsBaseApp struct {
	Contracts    BaseContract
	Baseapp      *baseabci.BaseApp
	ContractList []BaseContract
	Logger       log.Logger
	RootDir      string
}

func (base *QstarsBaseApp) Register(basecontract BaseContract) {
	base.ContractList = append(base.ContractList, basecontract)
}

func (base *QstarsBaseApp) loadX() {
	for index, c := range base.ContractList {
		fmt.Printf("arr[%d]=%d \n", index, c)
		c.RegisterKVCdc(base.Baseapp.GetCdc())
		c.StartX(base)
	}
}

func (base *QstarsBaseApp) Start() {

	db, err := dbm.NewGoLevelDB("kvstore", filepath.Join(base.RootDir, "data"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	base.Baseapp = baseabci.NewBaseApp("kvstore", base.Logger, db, nil)

	base.Baseapp.RegisterAccountProto(func() account.Account {
		return &account.BaseAccount{}
	})

	base.loadX()
}
