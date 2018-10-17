package lcd

import (
	"net/http"

	"github.com/QOSGroup/qstars/client/context"
	"github.com/QOSGroup/qstars/client/lcd/lib"
	"github.com/QOSGroup/qstars/version"
	"github.com/gorilla/mux"
)

type ResultCLIVersion struct {
	Version string `json:"version"`
}

func CLIVersionRegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/version",
		func(w http.ResponseWriter, r *http.Request) {
			result, err := CLIVersion()
			lib.HttpResponseWrapper(w, cliCtx.Codec, result, err)
		}).Methods("GET")
}

func CLIVersion() (*ResultCLIVersion, error) {
	v := version.GetVersion()

	result := &ResultCLIVersion{Version: string(v)}

	return result, nil
}

// connected node version REST handler endpoint
func NodeVersionRegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/node_version",
		func(w http.ResponseWriter, r *http.Request) {
			result, err := NodeVersion(NewReqNodeVersion(r))
			lib.HttpResponseWrapper(w, cliCtx.Codec, result, err)
		}).Methods("GET")
}

type ReqNodeVersion struct {
	CliCtx context.CLIContext
}

type ResultNodeVersion struct {
	Version string `json:"version"`
}

func NewReqNodeVersion(r *http.Request) *ReqNodeVersion {
	return &ReqNodeVersion{}
}

func NodeVersion(rnv *ReqNodeVersion) (*ResultNodeVersion, error) {
	v, err := rnv.CliCtx.Query("/app/version")
	if err != nil {
		return nil, err
	}

	result := &ResultNodeVersion{Version: string(v)}

	return result, nil
}
