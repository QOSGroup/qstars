package context

import (
	"fmt"
	"github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qstars/wire"
	"io"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"
	cmn "github.com/tendermint/tendermint/libs/common"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

const ACCOUNT_NOT_EXIST = "Account is not exsit."

// GetNode returns an RPC client. If the context's client is not defined, an
// error is returned.
func (ctx CLIContext) GetNode() (rpcclient.Client, error) {
	if ctx.Client == nil {
		return nil, errors.New("no RPC client defined")
	}

	return ctx.Client, nil
}

// Query performs a query for information about the connected node.
func (ctx CLIContext) Query(path string) (res []byte, err error) {
	return ctx.query(path, nil)
}

// QueryStore performs a query from a Tendermint node with the provided key and
// store name.
func (ctx CLIContext) QueryStore(key cmn.HexBytes, storeName string) (res []byte, err error) {
	return ctx.queryStore(key, storeName, "key")
}

// QueryStore performs a query from a Tendermint node with the provided key and
// store name.
func (ctx CLIContext) QueryInvestadCustom(key cmn.HexBytes) (res []byte, err error) {
	return ctx.query("/custom/jianqian/investad", key)
}

// QueryStore performs a query from a Tendermint node with the provided key and
// store name.
func (ctx CLIContext) QueryQOSAccount(key cmn.HexBytes) (res []byte, err error) {

	return ctx.query("/store/acc/key", key)

}

// QueryStore performs a query from a Tendermint node with the provided key and
// store name.
func (ctx CLIContext) QueryKV(key cmn.HexBytes) (res []byte, err error) {
	path := "/store/kv/key"
	return ctx.query(path, key)

}

// QuerySequenceIn performs a query from a Tendermint node with the provided chainid
func (ctx CLIContext) QuerySequenceIn(chainid string) (in int64, err error) {
	key := fmt.Sprintf("sequence/in/%s", chainid)
	path := "/store/qcp/key"
	res, err := ctx.query(path, []byte(key))
	if err != nil {
		return 0, err
	}

	err = ctx.Codec.UnmarshalBinaryBare(res, &in)

	return
}

// QuerySequenceOut performs a query from a Tendermint node with the provided chainid
func (ctx CLIContext) QuerySequenceOut(chainid string) (out int64, err error) {
	key := fmt.Sprintf("sequence/out/%s", chainid)
	path := "/store/qcp/key"

	res, err := ctx.query(path, []byte(key))
	if err != nil {
		return 0, err
	}

	err = ctx.Codec.UnmarshalBinaryBare(res, &out)

	return
}

// GetAccount queries for an account given an address and a block height. An
// error is returned if the query or decoding fails.
func (ctx CLIContext) GetAccount(address []byte, cdc *wire.Codec) (*account.QOSAccount, error) {

	result, err := ctx.QueryQOSAccount(address)
	if err != nil {
		return nil, err
	} else if len(result) == 0 {
		return nil, errors.New(ACCOUNT_NOT_EXIST)
	}

	var acc *account.QOSAccount
	err = cdc.UnmarshalBinaryBare(result, &acc)
	if err != nil {
		return nil, err
	}
	json, err := cdc.MarshalJSON(acc)
	fmt.Println(fmt.Sprintf("query addr is  %s", json))

	return acc, nil
}

// BroadcastTx broadcasts transaction bytes to a Tendermint node.
func (ctx CLIContext) BroadcastTx(tx []byte) (*ctypes.ResultBroadcastTxCommit, error) {
	node, err := ctx.GetNode()
	if err != nil {
		return nil, err
	}

	res, err := node.BroadcastTxCommit(tx)
	if err != nil {
		return res, err
	}

	if !res.CheckTx.IsOK() {
		return res, errors.Errorf("checkTx failed: (%d) %s",
			res.CheckTx.Code,
			res.CheckTx.Log)
	}

	if !res.DeliverTx.IsOK() {
		return res, errors.Errorf("deliverTx failed: (%d) %s",
			res.DeliverTx.Code,
			res.DeliverTx.Log)
	}

	return res, err
}

// BroadcastTxAsync broadcasts transaction bytes to a Tendermint node
// asynchronously.
func (ctx CLIContext) BroadcastTxAsync(tx []byte) (*ctypes.ResultBroadcastTx, error) {
	node, err := ctx.GetNode()
	if err != nil {
		return nil, err
	}

	res, err := node.BroadcastTxAsync(tx)
	if err != nil {
		return res, err
	}

	return res, err
}

// EnsureBroadcastTx broadcasts a transactions either synchronously or
// asynchronously based on the context parameters. The result of the broadcast
// is parsed into an intermediate structure which is logged if the context has
// a logger defined.
func (ctx CLIContext) EnsureBroadcastTx(txBytes []byte) (*ctypes.ResultBroadcastTxCommit, error) {
	if ctx.Async {
		_, err := ctx.ensureBroadcastTxAsync(txBytes)
		return nil, err
	}

	return ctx.ensureBroadcastTx(txBytes)
}

func (ctx CLIContext) ensureBroadcastTxAsync(txBytes []byte) (*ctypes.ResultBroadcastTx, error) {
	res, err := ctx.BroadcastTxAsync(txBytes)
	if err != nil {
		return res, err
	}

	if ctx.JSON {
		type toJSON struct {
			TxHash string
		}

		if ctx.Logger != nil {
			resJSON := toJSON{res.Hash.String()}
			bz, err := ctx.Codec.MarshalJSON(resJSON)
			if err != nil {
				return res, err
			}

			ctx.Logger.Write(bz)
			io.WriteString(ctx.Logger, "\n")
		}
	} else {
		if ctx.Logger != nil {
			io.WriteString(ctx.Logger, fmt.Sprintf("Async tx sent (tx hash: %s)\n", res.Hash))
		}
	}

	return res, nil
}

func (ctx CLIContext) ensureBroadcastTx(txBytes []byte) (*ctypes.ResultBroadcastTxCommit, error) {
	res, err := ctx.BroadcastTx(txBytes)
	if err != nil {
		return res, err
	}

	if ctx.JSON {
		// since JSON is intended for automated scripts, always include
		// response in JSON mode.
		type toJSON struct {
			Height   int64
			TxHash   string
			Response string
		}

		if ctx.Logger != nil {
			resJSON := toJSON{res.Height, res.Hash.String(), fmt.Sprintf("%+v", res.DeliverTx)}
			bz, err := ctx.Codec.MarshalJSON(resJSON)
			if err != nil {
				return res, err
			}

			ctx.Logger.Write(bz)
			io.WriteString(ctx.Logger, "\n")
		}

		return res, nil
	}

	if ctx.Logger != nil {
		resStr := fmt.Sprintf("Committed at block %d (tx hash: %s)\n", res.Height, res.Hash.String())

		if ctx.PrintResponse {
			resStr = fmt.Sprintf("Committed at block %d (tx hash: %s, response: %+v)\n",
				res.Height, res.Hash.String(), res.DeliverTx,
			)
		}

		io.WriteString(ctx.Logger, resStr)
	}

	return res, nil
}

// query performs a query from a Tendermint node with the provided store name
// and path.
func (ctx CLIContext) query(path string, key common.HexBytes) (res []byte, err error) {
	node, err := ctx.GetNode()
	if err != nil {
		return res, err
	}

	opts := rpcclient.ABCIQueryOptions{
		Height:  ctx.Height,
		Trusted: ctx.TrustNode,
	}

	result, err := node.ABCIQueryWithOptions(path, key, opts)
	if err != nil {
		return res, err
	}

	resp := result.Response
	if !resp.IsOK() {
		return res, errors.Errorf("query failed: (%d) %s", resp.Code, resp.Log)
	}
	// data from trusted node or subspace query doesn't need verification
	if ctx.TrustNode || !isQueryStoreWithProof(path) {
		return resp.Value, nil
	}

	err = ctx.verifyProof(path, resp)
	if err != nil {
		return nil, err
	}

	return resp.Value, nil
}

// queryStore performs a query from a Tendermint node with the provided a store
// name and path.
func (ctx CLIContext) queryStore(key cmn.HexBytes, storeName, endPath string) ([]byte, error) {
	path := fmt.Sprintf("/store/%s/%s", storeName, endPath)
	return ctx.query(path, key)
}

// isQueryStoreWithProof expects a format like /<queryType>/<storeName>/<subpath>
// queryType must be "store" and subpath must be "key" to require a proof.
func isQueryStoreWithProof(path string) bool {
	if !strings.HasPrefix(path, "/") {
		return false
	}

	paths := strings.SplitN(path[1:], "/", 3)
	switch {
	case len(paths) != 3:
		return false
	case paths[0] != "store":
		return false
	case store.RequireProof("/" + paths[2]):
		return true
	}

	return false
}