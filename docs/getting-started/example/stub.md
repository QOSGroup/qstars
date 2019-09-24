---
概述: QOS联盟链合约示例
---

stub中是启动本次示例合约的入口 注册在qstars后,由qstars启动时自动调用

##定义OrderStub对象

```go
type OrderStub struct {
}
```

要实现的BaseXTransaction接口 定义如下

```go
type BaseXTransaction interface {
	RegisterCdc(cdc *go_amino.Codec)  
	StartX(base *QstarsBaseApp) error
	ResultNotify(ctx context.Context, txQcpResult interface{}) *types.Result
	EndBlockNotify(ctx context.Context)
	CustomerQuery(ctx ctx.Context, route []string, req abci.RequestQuery) (res []byte, err types.Error)
	Name() string
}
```

##注册mapper

StartX中注册mapper 否则自定义mapper将不可调用
```go
func (astub OrderStub) StartX(base *baseapp.QstarsBaseApp) error {
	var orderMapper = NewOrderMapper(OrderMapperName)
	base.Baseapp.RegisterMapper(orderMapper)
	return nil
}
```
##注册交易结构体
将自定义的交易结构体序列化信息注册 这是交易由发起经传输后被识别的前提
```go
func (astub OrderStub) RegisterCdc(cdc *go_amino.Codec) {
	cdc.RegisterConcrete(&OrderTx{}, "example/supply/OrderTx", nil)
}
```


##处理跨链结果
跨链执行的交易会回调此方法 将执行结果返回进行下一步处理 在本示例中,跨链执行前将order用in.QcpOriginalExtends 存到了mapper中 在这里可以用此字段直接将order对象取出进行下一步处理
当跨链交易执行成功时  则删除以块高度+交易Hash组成的key保存的order 重新将order以订单id存储.

```go
func (astub OrderStub) ResultNotify(ctx context.Context, txQcpResult interface{}) *types.Result {
	in := txQcpResult.(*txs.QcpTxResult)
	log.Debugf("ResultNotify QcpOriginalSequence:%s, result:%+v", string(in.QcpOriginalSequence), txQcpResult)
	var resultCode types.CodeType
	qcpTxResult, ok := baseabci.ConvertTxQcpResult(txQcpResult)
	if ok == false {
		log.Errorf("ResultNotify ConvertTxQcpResult error.")
		resultCode = types.CodeTxDecode
	} else {
		log.Errorf("ResultNotify update status")

		orginalTxHash := in.QcpOriginalExtends //orginalTx.abc
		kvMapper := ctx.Mapper(common.QSCResultMapperName).(*common.KvMapper)
		initValue := ""
		kvMapper.Get([]byte(orginalTxHash), &initValue)
		if initValue != astub.Name() {
			log.Info("This is not my response.")
			return nil
		}
		//put result to map for client query
		c := strconv.FormatInt((int64)(qcpTxResult.Result.Code), 10)
		c = c + " " + qcpTxResult.Result.Log
		kvMapper.Set([]byte(orginalTxHash), c)

		orderMapper := ctx.Mapper(OrderMapperName).(*OrderMapper)
		order := orderMapper.GetOrder(orginalTxHash)
		if order != nil {
			orderMapper.SaveOrder(order.Id, order)
			orderMapper.DeleteOrder(orginalTxHash)
		}
		resultCode = types.CodeOK
	}
	rr := types.Result{
		Code: resultCode,
	}
	return &rr
}
```

##注册到qstars
打开qstars/star/app.go 在init方法中添加

```go
	registerType((*supply.Stub)(nil))
```

[cmd.go    命令行方式与链码交互入口](cmd.md)