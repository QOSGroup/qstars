package baseapp

import (
	"fmt"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/server"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/utility"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
	"os"
	"path/filepath"
	"strings"

	"github.com/QOSGroup/qbase/context"
	bctypes "github.com/QOSGroup/qbase/example/basecoin/types"
	go_amino "github.com/tendermint/go-amino"
	dbm "github.com/tendermint/tendermint/libs/db"
	ctx "github.com/QOSGroup/qbase/context"
	abci "github.com/tendermint/tendermint/abci/types"
)

type QStarsContext struct {
	ServerContext    *server.Context
	QStarsSignerPriv crypto.PrivKey
	QStarsTransactions []string
}

var qCtx *QStarsContext

func GetServerContext() *QStarsContext {
	return qCtx
}

func InitApp() {
	qCtx = &QStarsContext{
		ServerContext: server.NewDefaultContext(),
	}
}

/**
	startup a qstar chain instance
 */
func NewAPP(sconf *config.ServerConf , cdc *go_amino.Codec) (QstarsBaseApp, error) {
	_, _, qCtx.QStarsSignerPriv = utility.PubAddrRetrievalFromAmino(sconf.QStarsPrivateKey, cdc)
	qCtx.QStarsTransactions = strings.Split(sconf.QStarsTransactions, ",")
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "main")
	qstarts := QstarsBaseApp{
		Logger:  logger,
		RootDir: sconf.RootDir,
	}
	return qstarts, nil
}

type QstarsBaseApp struct {
	Transactions    BaseXTransaction
	Baseapp      *baseabci.BaseApp
	TransactionList []BaseXTransaction
	Logger       log.Logger
	RootDir      string
}

//call every transaction to register
func (base *QstarsBaseApp) Register(basecontract BaseXTransaction) {
	base.TransactionList = append(base.TransactionList, basecontract)
}

//Load transaction
func (base *QstarsBaseApp) loadX() error {
	for index, c := range base.TransactionList {
		base.Logger.Info("arr[%d]=%d \n", index, c)
		err := c.StartX(base)
		if err != nil {
			return err
		}
	}

	//qbase need qstars to call this, is right?
	if err := base.Baseapp.LoadLatestVersion(); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

//Rigster every x transaction
func (base *QstarsBaseApp) RegisterCDC(cdc *go_amino.Codec) {
	for _, c := range base.TransactionList {
		c.RegisterCdc(cdc)
	}
}

func (base *QstarsBaseApp) TxQcpResultHandler(ctx context.Context, txQcpResult interface{}) types.Result {
	var rr types.Result
	for _, c := range base.TransactionList {
		tmprr := c.ResultNotify(ctx, txQcpResult)
		if tmprr != nil {
			rr = *tmprr
		}
	}
	return rr
}

/**
	start transaction
 */
func (base *QstarsBaseApp) Start() error {
	//this store used to store chain information
	db, err := dbm.NewGoLevelDB("qstarstore", filepath.Join(base.RootDir, "data"))
	if err != nil {
		fmt.Println(err)
		return err
	}

	base.Baseapp = baseabci.NewBaseApp("qstarstore", base.Logger, db, base.RegisterCDC)

	//qbase need register account
	base.Baseapp.RegisterAccountProto(bctypes.NewAppAccount)
	//qbase need register result handler
	base.Baseapp.RegisterTxQcpResultHandler(base.TxQcpResultHandler)
	//qbase need register qstar(QCP) signer
	base.Baseapp.RegisterTxQcpSigner(GetServerContext().QStarsSignerPriv)

	base.Baseapp.SetEndBlocker(func(ctx ctx.Context, req abci.RequestEndBlock) abci.ResponseEndBlock{
		for _, c := range base.TransactionList {
			c.EndBlockNotify(ctx)
		}
		return abci.ResponseEndBlock{}
	})
	return base.loadX()
}

