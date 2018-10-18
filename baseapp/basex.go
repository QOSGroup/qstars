package baseapp

import (
	"github.com/QOSGroup/qbase/mapper"
	go_amino "github.com/tendermint/go-amino"
)

type BaseContract interface {
	mapper.IMapper
	RegisterKVCdc(cdc *go_amino.Codec)
	StartX(base *QstarsBaseApp) error
}
