package baseapp

import (

	go_amino "github.com/tendermint/go-amino"
	"github.com/QOSGroup/qbase/mapper"
)
type BaseContract interface{
	mapper.IMapper
	RegisterKVCdc(cdc *go_amino.Codec)
	StartX(base *QstarsBaseApp)
}