---
概述: QOS联盟链合约示例
---

# 定义交易结构体

首先我们要定义一个交易结构体 该结构体里的属性根据业务逻辑需要来定义,但必须要包含验证签名者的地址 也就是签署这个交易的私钥所对应的公钥.如果是多方签名 这里就要包含多方私钥所对应的公钥
本示例中我们定义的结构体如下

```go
type OrderTx struct {
	Wrapper     *txs.TxStd    //需要到主网执行的交易信息
	Address     types.Address //发起者地址
	Id          string        //订单编号
	OrderTo     types.Address //订单接收者地址
	OrderAmount types.BigInt  //订单金额
	Gas         types.BigInt  //gas
}
```
单方签名的一个订单的简单信息, 包含了要到QOS主链上转账的交易Wrapper属性, 其他属性信息保存在联盟链中
接下来就是围绕OrderTx这个对象展开操作.   在联盟链合约开发中交易对象必须实现ITx接口.该接口定义如下


```go
type ITx interface {
	ValidateData(ctx context.Context) error //基础检查 
	Exec(ctx context.Context) (result types.Result, crossTxQcp *TxQcp)	//执行业务逻辑,crossTxQcp: 需要进行跨链处理的TxQcp。
	GetSigner() []types.Address //签名者
	CalcGas() types.BigInt      //计算gas
	GetGasPayer() types.Address //gas付费人
	GetSignData() []byte        //获取签名字段
}
```
## 交易检查
ValidateData 中做初步检查 通常是数据的合规检查 如果检查未通过测此交易将不被执行  在本示例中我们仅检查了数据是否非空

```go
func (tx *OrderTx) ValidateData(ctx context.Context) error {
	if tx.Address == nil || tx.Address.Empty() {
		return errors.New("发起者地址不能为空")
	}
	if strings.TrimSpace(tx.Id) == "" {
		return errors.New("订单不能为空")
	}
	if tx.OrderTo == nil || tx.OrderTo.Empty() {
		return errors.New("接单企业地址不能为空")
	}
	return nil
}
```


## 获取签名者

验证签名数据(GetSignData())的地址 

```go
func (tx *OrderTx) GetSigner() []types.Address {
	return []types.Address{tx.Address}
}
```

## 获取签名数据

将需要保护的数据加到此方法中 如果数据在传输过程中被修改则验证是不会通过 
```go
func (tx *OrderTx) GetSignData() (ret []byte) {
	ret = append(ret, tx.Wrapper.ITxs[0].GetSignData()...)
	ret = append(ret, tx.Address.Bytes()...)
	ret = append(ret, tx.Id...)
	ret = append(ret, tx.OrderTo...)
	ret = append(ret, tx.OrderAmount.String()...)
	ret = append(ret, tx.Gas.String()...)
	return ret
}
```

## 执行主逻辑
本例中qsc的转账部分需要到QOS主链上执行 所以需要封装跨链交易TxQcp交由cassin将交易中继到QOS主链.主链返回的结果由[stub.go](stub.md)中监听
当跨链结果返回时 识别此交易的方法是通过 crossTxQcps.Extends = key 中定义的标识 它由块高度和交易Hash组成.来保证唯一性
 
```go
func (tx *OrderTx) Exec(ctx context.Context) (result types.Result, crossTxQcps *txs.TxQcp) {
	kvMapper := ctx.Mapper(starcommon.QSCResultMapperName).(*starcommon.KvMapper)
	heigth := strconv.FormatInt(ctx.BlockHeight(), 10)
	txhash := (common.HexBytes)(tmhash.Sum(ctx.TxBytes()))
	key := GetResultKey(heigth, txhash.String())
	kvMapper.Set([]byte(key), OrderStub{}.Name())
	crossTxQcps.TxStd = tx.Wrapper
	crossTxQcps.To = config.GetServerConf().QOSChainName
	crossTxQcps.Extends = key
	result = types.Result{
		Code: types.CodeOK,
	}
	orderMapper := ctx.Mapper(OrderMapperName).(*OrderMapper)
	order:=&Order{Id:tx.Id,OrderName:tx.Address.String(),OrderTo:tx.OrderTo.String(),OrderAmount:tx.OrderAmount,Status:1}
	orderMapper.SaveOrder(key,order)
	return
}
```

完整代码参考 handler.go


[mapper.go  数据存取(lvevl db)逻辑](mapper.md)