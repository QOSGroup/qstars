package baseapp

import (
	"fmt"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/server"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	account "github.com/QOSGroup/qos/types"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/utility"
	"github.com/QOSGroup/qstars/x/common"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
	"os"
	"path/filepath"
	"strings"

	"github.com/QOSGroup/qbase/context"
	ctx "github.com/QOSGroup/qbase/context"
	go_amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
)

type QStarsContext struct {
	ServerContext      *server.Context
	QStarsSignerPriv   crypto.PrivKey
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
func NewAPP(sconf *config.ServerConf, cdc *go_amino.Codec) (QstarsBaseApp, error) {
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
	Baseapp         *baseabci.BaseApp
	TransactionList []BaseXTransaction
	Logger          log.Logger
	RootDir         string
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

func (base *QstarsBaseApp) TxQcpResultHandler(ctx context.Context, txQcpResult interface{}) {
	defer func() {
		if msg := recover(); msg != nil {
			fmt.Printf("baseapp TxQcpResultHandler panic %+v\n", msg)
		}
	}()

	fmt.Printf("baseapp TransactionList %+v\n", base.TransactionList)
	for _, c := range base.TransactionList {
		fmt.Printf("baseapp kvMapper c:%+v\n", c)
		initValue := ""
		fmt.Printf("baseapp txQcpResult %+v\n", txQcpResult)
		in := txQcpResult.(*txs.QcpTxResult)
		key := in.QcpOriginalExtends
		fmt.Printf("baseapp kvMapper-1 key:%s, value:%s, name:%s\n", key, initValue, c.Name())
		kvMapper := ctx.Mapper(common.QSCResultMapperName).(*common.KvMapper)
		kvMapper.Get([]byte(key), &initValue)
		fmt.Printf("baseapp kvMapper-2 key:%s, value:%s, name:%s\n", key, initValue, c.Name())
		if initValue == c.Name() {
			c.ResultNotify(ctx, txQcpResult)
			return
		}
	}
	return
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
	base.Baseapp.RegisterAccountProto(account.ProtoQOSAccount)
	//qbase need register result handler
	base.Baseapp.RegisterTxQcpResultHandler(base.TxQcpResultHandler)
	//qbase need register qstar(QCP) signer
	base.Baseapp.RegisterTxQcpSigner(GetServerContext().QStarsSignerPriv)

	var handler baseabci.CustomQueryHandler
	handler = func(ctx ctx.Context, route []string, req abci.RequestQuery) (res []byte, err types.Error) {
		for _, c := range base.TransactionList {
			response, err := c.CustomerQuery(ctx, route, req)
			if (response != nil) && (err == nil) {
				return response, nil
			} else {
				if (response == nil) && (err != nil) {
					return response, err
				}
			}
		}
		return nil, nil
	}

	base.Baseapp.RegisterCustomQueryHandler(handler)
	base.Baseapp.SetEndBlocker(func(ctx ctx.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		for _, c := range base.TransactionList {
			c.EndBlockNotify(ctx)
		}
		return abci.ResponseEndBlock{}
	})
	return base.loadX()
}
