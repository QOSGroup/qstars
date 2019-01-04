package context

import (
	"fmt"
	"github.com/QOSGroup/qbase/store"
	"github.com/pkg/errors"
	"io"
	"strings"

	"github.com/QOSGroup/qstars/client"
	"github.com/QOSGroup/qstars/wire"

	"github.com/spf13/viper"

	rpcclient "github.com/tendermint/tendermint/rpc/client"
	tmlite "github.com/tendermint/tendermint/lite"



	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/merkle"
	tmliteErr "github.com/tendermint/tendermint/lite/errors"
	tmliteProxy "github.com/tendermint/tendermint/lite/proxy"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	tmtypes "github.com/tendermint/tendermint/types"
)

const ctxAccStoreName = "qosaccount"
const ctxKVStoreName = "kv"

// CLIContext implements a typical CLI context created in SDK modules for
// transaction handling and queries.
type CLIContext struct {
	Codec           *wire.Codec
	Client          rpcclient.Client
	Logger          io.Writer
	Height          int64
	NodeURI         string
	FromAddressName string
	AccountStore    string
	KVStore         string
	TrustNode       bool
	UseLedger       bool
	Async           bool
	JSON            bool
	PrintResponse   bool
	Verifier      tmlite.Verifier
}

// NewCLIContext returns a new initialized CLIContext with parameters from the
// command line using Viper.
func NewCLIContext1(nodeURI string) CLIContext {
	var rpc rpcclient.Client

	if nodeURI != "" {
		rpc = rpcclient.NewHTTP(nodeURI, "/websocket")
	}

	return CLIContext{
		Client:          rpc,
		NodeURI:         nodeURI,
		AccountStore:    ctxAccStoreName,
		KVStore:         ctxKVStoreName,
		FromAddressName: viper.GetString(client.FlagFrom),
		Height:          viper.GetInt64(client.FlagHeight),
		TrustNode:       viper.GetBool(client.FlagTrustNode),
		UseLedger:       viper.GetBool(client.FlagUseLedger),
		Async:           viper.GetBool(client.FlagAsync),
		JSON:            viper.GetBool(client.FlagJson),
		PrintResponse:   viper.GetBool(client.FlagPrintResponse),
	}
}

// WithCodec returns a copy of the context with an updated codec.
func (ctx CLIContext) WithCodec(cdc *wire.Codec) CLIContext {
	ctx.Codec = cdc
	return ctx
}

// WithLogger returns a copy of the context with an updated logger.
func (ctx CLIContext) WithLogger(w io.Writer) CLIContext {
	ctx.Logger = w
	return ctx
}

// WithAccountStore returns a copy of the context with an updated AccountStore.
func (ctx CLIContext) WithAccountStore(accountStore string) CLIContext {
	ctx.AccountStore = accountStore
	return ctx
}

// WithKVStore returns a copy of the context with an updated KVStore.
func (ctx CLIContext) WithKVStore(kvStore string) CLIContext {
	ctx.KVStore = kvStore
	return ctx
}

// WithFromAddressName returns a copy of the context with an updated from
// address.
func (ctx CLIContext) WithFromAddressName(addrName string) CLIContext {
	ctx.FromAddressName = addrName
	return ctx
}

// WithTrustNode returns a copy of the context with an updated TrustNode flag.
func (ctx CLIContext) WithTrustNode(trustNode bool) CLIContext {
	ctx.TrustNode = trustNode
	return ctx
}

// WithNodeURI returns a copy of the context with an updated node URI.
func (ctx CLIContext) WithNodeURI(nodeURI string) CLIContext {
	ctx.NodeURI = nodeURI
	ctx.Client = rpcclient.NewHTTP(nodeURI, "/websocket")
	return ctx
}

// WithClient returns a copy of the context with an updated RPC client
// instance.
func (ctx CLIContext) WithClient(client rpcclient.Client) CLIContext {
	ctx.Client = client
	return ctx
}

// WithUseLedger returns a copy of the context with an updated UseLedger flag.
func (ctx CLIContext) WithUseLedger(useLedger bool) CLIContext {
	ctx.UseLedger = useLedger
	return ctx
}

// verifyProof perform response proof verification.
func (ctx CLIContext) verifyProof(queryPath string, resp abci.ResponseQuery) error {
	if ctx.Verifier == nil {
		return fmt.Errorf("missing valid certifier to verify data from distrusted node")
	}

	// the AppHash for height H is in header H+1
	commit, err := ctx.Verify(resp.Height + 1)
	if err != nil {
		return err
	}

	// TODO: Instead of reconstructing, stash on CLIContext field?
	prt := store.DefaultProofRuntime()

	// TODO: Better convention for path?
	storeName, err := parseQueryStorePath(queryPath)
	if err != nil {
		return err
	}

	kp := merkle.KeyPath{}
	kp = kp.AppendKey([]byte(storeName), merkle.KeyEncodingURL)
	kp = kp.AppendKey(resp.Key, merkle.KeyEncodingURL)

	err = prt.VerifyValue(resp.Proof, commit.Header.AppHash, kp.String(), resp.Value)
	if err != nil {
		return errors.Wrap(err, "failed to prove merkle proof")
	}

	return nil
}


// Verify verifies the consensus proof at given height.
func (ctx CLIContext) Verify(height int64) (tmtypes.SignedHeader, error) {
	check, err := tmliteProxy.GetCertifiedCommit(height, ctx.Client, ctx.Verifier)
	switch {
	case tmliteErr.IsErrCommitNotFound(err):
		return tmtypes.SignedHeader{}, ErrVerifyCommit(height)
	case err != nil:
		return tmtypes.SignedHeader{}, err
	}

	return check, nil
}

