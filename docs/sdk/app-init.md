# 开发自己的联盟链业务模块

在github.com/QOSGroup/qstars/x目录下新建自己的业务模块test目录并在test目录下创建: `./stub.go`. 用于注册mapper、自定义结构体编解码器、跨链交易结果接收等

在 `stub.go`, 需要实现 github.com/QOSGroup/qstars/baseapp/basex.go中BaseXTransaction接口中定义的方法

```go
package test

import (
	context "github.com/QOSGroup/qbase/context"
	types "github.com/QOSGroup/qbase/types"
	baseapp "github.com/QOSGroup/qstars/baseapp"
	go-amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
)
```

引入模块说明:

- [`context`]用于上下文信息传递
- [`types`]定义了账户地址、币等常用的类型.
- [`baseapp`]定义了要实现的接口.
- [`go-amino`]用于自定义交易结构体编解码器注册.
- [`abci`]跨链结果结构体定义.


定义结构体
```go
type TestStub struct {}
```
接下来实现以下baseapp.BaseXTransaction中定义的方法
```go
- func (ts TestStub) RegisterCdc(cdc *go_amino.Codec){          //注册交易结构体编解码器
	cdc.RegisterConcrete(&KvstoreTx{}, "test/KvstoreTx", nil)
  }
- func (ts TestStub) StartX(base *QstarsBaseApp) error{         //注册mapper
	    var testMapper = NewTestMapper(TestMapperName)
	    base.Baseapp.RegisterMapper(testMapper)
	    return nil
  }
- func (ts TestStub) ResultNotify(ctx context.Context, txQcpResult interface{}) *types.Result{return nil}    //接收跨链结果 无跨链返回nil
- func (ts TestStub) EndBlockNotify(ctx context.Context){}    //接收打块完成通知
- func (ts TestStub) CustomerQuery(ctx ctx.Context, route []string, req abci.RequestQuery) (res []byte, err types.Error){return nil, nil}    //自定义查询
- func (ts TestStub) Name() string{	return "TestStub"} 返回本模块的名字
```
###############################################################################

同目录下创建: `./mapper.go`.用于数据存储
在 `mapper.go`, 需要实现github.com/QOSGroup/qbase/mapper中 IMapper接口中定义的方法

```go
package test

import (
	mapper "github.com/QOSGroup/qbase/mapper"
	types "github.com/QOSGroup/qbase/types"
)
```
引入模块说明:

- [`mapper`]继承qbase中已封装好的通用方法
- [`types`]定义了账户地址、币等常用的类型.

接下来给自己的mapper命名 要求全局唯一 并定义TestMapper 继承mapper.BaseMapper
```go
const TestMapperName = "test"

type TestMapper struct {
	*mapper.BaseMapper
}
```

接下来定义初始化方法 在teststub.go中的StartX中引用NewTestMapper来注册
重新实现Copy方法 （Copy方法是IMapper接口中定义的BaseMapper已实现 这里需要重新实现）
```go
func NewTestMapper(MapperName string) *TestMapper {
	var testMapper = TestMapper{}
	testMapper.BaseMapper = mapper.NewBaseMapper(nil, MapperName)
	return &testMapper
}

func (mapper *TestMapper) Copy() mapper.IMapper {
	cpyMapper := &TestMapper{}
	cpyMapper.BaseMapper = mapper.BaseMapper.Copy()
	return cpyMapper
}
```
增加自己的业务逻辑（Getter Setter） 这里以kv为例
```go
func (mapper TestMapper SaveKV(key string, value string) {
	mapper.BaseMapper.Set([]byte(key), value)
}

func (mapper *TestMapper) GetKey(key string) (v string) {
	mapper.BaseMapper.Get([]byte(key), &v)
	return
}
```
###############################################################################

同目录下创建: `./handler.go`.定义交易结构及交易逻辑
在 `handler.go`, 需要实现github.com/QOSGroup/qbase/txs中 ITx接口中定义的方法

```go
package test

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
)
```
引入模块说明:
- [`context`]用于上下文信息传递
- [`types`]定义了账户地址、币等常用的类型.
- [`txs`]定义了交易及QCP交易及交易结果等类型

根据业务定义自己的交易结构体及初始化方法 这里以kv为例
```go
type KvstoreTx struct {
	Key   []byte
	Value []byte
	Bytes []byte
}

func NewKvstoreTx(key []byte, value []byte) KvstoreTx {
	return KvstoreTx{
		Key:   key,
		Value: value,
	}
}
```
接下来实现ITx接口
```go
- func (ts KvstoreTx)	ValidateData(ctx context.Context) error {}//检测
- func (ts KvstoreTx)	Exec(ctx context.Context) (result types.Result, crossTxQcp *TxQcp){}	//执行业务逻辑, crossTxQcp: 需要进行跨链处理的TxQcp 业务端实现中crossTxQcp只需包含`to` 和 `txStd`
- func (ts KvstoreTx)	GetSigner() []types.Address{} 签名者
- func (ts KvstoreTx)	CalcGas() types.BigInt{}   计算gas
- func (ts KvstoreTx)	GetGasPayer() types.Address{} gas付费人
- func (ts KvstoreTx)	GetSignData() []byte{}     获取签名字段
```
###############################################################################

同目录下创建: `./process.go` 定义供客户端（cmd 或RESTFul）调用发起交易用
```go
package test

import (
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/client/utils"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/utility"
	"github.com/QOSGroup/qstars/wire"
	"github.com/tendermint/tendermint/crypto/ed25519"
	 qstartypes "github.com/QOSGroup/qstars/types"
	)
```
引入模块说明:
- [`txs`]定义了交易及QCP交易及交易结果等类型
- [`types`]定义了账户地址、币等常用的类型.
- [`utils`]用于发起交易.
- [`config`]用于获取chainid等配置信息.
- [`utility`]用于私钥转换.
- [`ed25519`]公钥对类型.
- [`wire`]编解码器注册.
- [`qstartypes`]用于账户地址转换.

根据提供的私钥对交易签名
```go
func WrapToStdTx(cdc *wire.Codec,privkey,key, value, chainID string) string {
	cliCtx := *config.GetCLIContext().QSCCliContext                              //获取context
	kv := NewKvstoreTx([]byte(key), []byte(value))                               //初始化交易
	_, addrben32, priv := utility.PubAddrRetrievalFromAmino(privkey, cdc)        //根据私钥字符串转换成私钥对象
	from, _ := qstartypes.AccAddressFromBech32(addrben32)                        //根据私钥获取发起人账户地址
	account:=account.AddressStoreKey(from)                                       //根据私钥获取账户信息
	var nonce int64=0
	acc, err := cliCtx.GetAccount(account, cdc)                                  //联盟链内获取账户nonce信息
	if err != nil {
		nonce=0
	}else{
		nonce=int64(acc.Nonce)
	}
	nonce++                                                                       //nonce加1 防止双花
	tx:=txs.NewTxStd(kv, chainID, types.NewInt(int64(10000)))                     //封装标准交易结构体
	signdata,_:=tx.SignTx(priv,nonce,chainID)                                     //用私钥签名返回签名信息
	tx.Signature = []txs.Signature{txs.Signature{                                 //将签名信息填充到交易结构体
		Pubkey:    priv.PubKey(),
		Signature: signdata,
		Nonce:     nonce,
	}}
	hash, _, err := utils.SendTx(cliCtx, cdc, tx)                                 //发起交易并返回apphash
	return hash
	}
```
###############################################################################

同目录下创建: `./cmd.go` 定义供客户端命令行方式发起交易

```go
package test
import (
	"github.com/QOSGroup/qstars/wire"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)
```
引入模块说明:

- [`wire`]编解码器注册.
- [`cobra`]定义命令行参数
- [`viper`]读取命令行输入

定义命令行参数
```go
const (
	flagKey        = "key"
	flagValue      = "value"
	flagPrivateKey = "private"
	chainIdFlag    = "chain-id"
)
```
定义命令行逻辑
```go
// SendTxCmd will create a send tx and sign it with the given key.
func SendKVCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kvset",
		Short: "Create and sign a send set kv tx",
		RunE: func(cmd *cobra.Command, args []string) error {

			privatekey := viper.GetString(flagPrivateKey)
			key := viper.GetString(flagKey)
			value := viper.GetString(flagValue)
			result := WrapToStdTx(cdc, privatekey, key, value, viper.GetString(chainIdFlag))
			fmt.Println(result)
			return nil
		},
	}
	cmd.Flags().String(flagKey, "", "Key")
	cmd.Flags().String(flagValue, "", "Value")
	cmd.Flags().String(flagPrivateKey, "", "Private key")

	return cmd
}
```

###############################################################################

最后让我们把此模块代码加到启动项中
客户端:
```go
github.com/QOSGroup/qstars/cmd/qstarscli/main.go中main方法中增加
	rootCmd.AddCommand(
		test.SendKVCmd(cdc),
	)
```
服务端:
```go
github.com/QOSGroup/qstars/app/star.go中init方法中增加
	registerType((*test.TestStub)(nil))
```

编译		github.com/QOSGroup/qstars/cmd/qstars/main.go 启动服务端
编译		github.com/QOSGroup/qstars/cmd/qstarscli/main.go 启动客户端



