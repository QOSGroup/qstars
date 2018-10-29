package baseapp

import (
	"fmt"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
	"os"
	"path/filepath"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qbase/server"

	"github.com/QOSGroup/qbase/context"
	bctypes "github.com/QOSGroup/qbase/example/basecoin/types"
	go_amino "github.com/tendermint/go-amino"
	dbm "github.com/tendermint/tendermint/libs/db"
)
type QStarsContext struct {
	ServerContext *server.Context
	QStarsSignerPriv crypto.PrivKey
}

var qCtx *QStarsContext

func GetServerContext() *QStarsContext{
	return qCtx
}

func InitApp(){
	qCtx = &QStarsContext{
		ServerContext:server.NewDefaultContext(),

	}
}
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

func (base *QstarsBaseApp) loadX() error{
	for index, c := range base.ContractList {
		fmt.Printf("arr[%d]=%d \n", index, c)
		err := c.StartX(base)
		if err!=nil{
			return err
		}
	}
	return nil
}

func (base *QstarsBaseApp) RegisterCDC(cdc *go_amino.Codec){
	for _, c := range base.ContractList {
		c.RegisterKVCdc(cdc)
	}
}

func (base *QstarsBaseApp) TxQcpResultHandler (ctx context.Context, txQcpResult interface{}) types.Result {
	var rr types.Result
	for _, c := range base.ContractList {
		tmprr := c.ResultNotify(ctx,txQcpResult)
		if tmprr!=nil{
			rr = *tmprr
		}
	}
	return rr
}

func (base *QstarsBaseApp) Start() error{

	db, err := dbm.NewGoLevelDB("kvstore", filepath.Join(base.RootDir, "data"))
	if err != nil {
		fmt.Println(err)
		return err
	}

	base.Baseapp = baseabci.NewBaseApp("kvstore", base.Logger, db, base.RegisterCDC)

	base.Baseapp.RegisterAccountProto(bctypes.NewAppAccount)
	base.Baseapp.RegisterTxQcpResultHandler(base.TxQcpResultHandler)
	base.Baseapp.RegisterTxQcpSigner(GetServerContext().QStarsSignerPriv)
	return base.loadX()
}
