package baseapp

import (
	"github.com/QOSGroup/qbase/mapper"
	go_amino "github.com/tendermint/go-amino"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qbase/context"
)

type BaseContract interface {
	mapper.IMapper
	RegisterKVCdc(cdc *go_amino.Codec)
	StartX(base *QstarsBaseApp) error
	ResultNotify(ctx context.Context, txQcpResult interface{}) *types.Result
}
