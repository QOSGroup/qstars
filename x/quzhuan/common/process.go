package common

import (
	"encoding/hex"
	"encoding/json"
	"github.com/QOSGroup/qbase/account"
	qosaccount "github.com/QOSGroup/qos/types"

	"errors"
	"github.com/QOSGroup/qbase/txs"
	qstarstypes "github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/utility"
	"github.com/QOSGroup/qstars/wire"
	"github.com/QOSGroup/qstars/x/common"
	"log"

	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"sync"

	"github.com/QOSGroup/qbase/types"
)

const (
	ResultCodeSuccess       = "0"
	ResultCodeQstarsTimeout = "-2"
	ResultCodeQOSTimeout    = "-1"
	ResultCodeInternalError = "500"
)

var RPC rpcclient.Client

var lock sync.Mutex

func getRPC(url string) rpcclient.Client {

	lock.Lock()
	defer lock.Unlock()
	if RPC == nil {
		RPC = rpcclient.NewHTTP(url, "/websocket")
	}
	return RPC
}

type ResultResp struct {
	Code   string          `json:"code"`
	Height int64           `json:"height"`
	Hash   string          `json:"hash,omitempty"`
	Reason string          `json:"reason,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
}

func (ri ResultResp) Marshal() string {
	if ri.Code == ResultCodeSuccess {
		return string(hex.EncodeToString(ri.Result))
	}
	return string(ri.Result)
}

func CommHandler(cdc *wire.Codec, url,qscchainid,funcName, privatekey, argstr string) string {

	var result ResultResp
	result.Code = ResultCodeSuccess

	var args []string
	err := json.Unmarshal([]byte(argstr), &args)
	if err != nil {
		return common.NewErrorResult(common.ResultCodeInternalError, 0, "", err.Error()).Marshal()
	}
	tx, berr := commHandler(cdc,url, qscchainid,funcName, privatekey, args)
	if berr != "" {
		return berr
	}

	js, err := cdc.MarshalBinaryBare(tx)
	if err != nil {
		log.Printf("CommHandler err:%s", err.Error())
		result.Code = ResultCodeInternalError
		result.Reason = err.Error()
		return result.Marshal()
	}
	result.Result = json.RawMessage(js)


	txb:=result.Marshal()

	res,err:=BroadcastTransferTxToQSC(cdc,txb,url)

	if err!=nil{
		return common.NewErrorResult(common.ResultCodeInternalError, 0, "", err.Error()).Marshal()
	}

	if res.CheckTx.Code!=0{
		result.Code=ResultCodeInternalError
		result.Hash=res.Hash.String()
		result.Reason=res.CheckTx.Log
		result.Height=res.Height
	}
	if res.DeliverTx.Code!=0{
		result.Code=ResultCodeInternalError
		result.Hash=res.Hash.String()
		result.Reason=res.DeliverTx.Log
		result.Height=res.Height
	}

	result.Hash=res.Hash.String()
	result.Reason="success"
	result.Height=res.Height

	return result.Marshal()
}

func commHandler(cdc *wire.Codec, url, qscchainid,funcName, privatekey string, args []string) (*txs.TxStd, string) {
	_, addrben32, priv := utility.PubAddrRetrievalFromAmino(privatekey, cdc)
	from, _ := qstarstypes.AccAddressFromBech32(addrben32)
	gas := types.NewInt(int64(200000))
	//key := account.AddressStoreKey(from)
	var qscnonce int64 = 0
	qscacc, err := RpcQueryAccount(cdc,url, from)
	if err != nil||qscacc==nil {
		qscnonce = 0
	} else {
		qscnonce = int64(qscacc.Nonce)
	}
	qscnonce += 1
	tx := &QuZhuanTx{}
	tx.Address = []types.Address{from}
	tx.FuncName = funcName
	tx.Args = args
	tx.Gas = gas
	//qscchainid := config.GetCLIContext().Config.QSCChainID
	tx2 := txs.NewTxStd(tx, qscchainid, gas)
	signature2, _ := tx2.SignTx(priv, qscnonce, qscchainid, qscchainid)
	tx2.Signature = []txs.Signature{txs.Signature{
		Pubkey:    priv.PubKey(),
		Signature: signature2,
		Nonce:     qscnonce,
	}}
	return tx2, ""
}

// 提交到联盟链上
func BroadcastTransferTxToQSC(cdc *wire.Codec, txb, url string) (*ctypes.ResultBroadcastTxCommit,error) {
	txBytes, err := hex.DecodeString(txb)
	if err != nil {
		return nil,err
	}

	var res *ctypes.ResultBroadcastTxCommit


	res,err=getRPC(url).BroadcastTxCommit(txBytes)

	//switch broadcastModes {
	//case "sync":
	//	res, err = getRPC(url).BroadcastTxSync(txBytes)
	//	//默认异步
	//default:
	//	res, err = getRPC(url).BroadcastTxAsync(txBytes)
	//}
	if err != nil {
		return nil,err
	}
	return res,nil
}

func RpcQueryAccount(cdc *wire.Codec,url string, addr types.Address) (*qosaccount.QOSAccount, error) {
	key := account.AddressStoreKey(addr)
	opts := rpcclient.ABCIQueryOptions{
		Height: 0,
		Prove:  true,
	}
	result, err := getRPC(url).ABCIQueryWithOptions("/store/acc/key", key, opts)
	if err != nil {
		return nil, err
	}
	resp := result.Response
	if !resp.IsOK() {
		return nil, errors.New("query failed")
	}
	var acc *qosaccount.QOSAccount
	err = cdc.UnmarshalBinaryBare(resp.Value, &acc)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

func ParameterFormat(paras []string) (string, error) {
	argstr, err := json.Marshal(paras)

	if err != nil {
		return "", err
	}
	return string(argstr), nil
}
