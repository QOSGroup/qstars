---
概述: QOS联盟链合约示例
---


## 签名交易并广播
用发起者私钥将交易签名 并广播到链上 通常做SDK使用 供cmd命令行调用 或者封装RESTful api接口  甚至可用于对接数字钱包


通过私钥获取其对应的地址 并半 地址转换成存储通用的key用来查询nonce
```go
_, addrben32, priv := utility.PubAddrRetrievalFromAmino(privatestr, cdc)
from, err := types.AccAddressFromBech32(addrben32)
key := account.AddressStoreKey(from)
````


## 获取nonce
从QOS链上查询发起交易账户的当前nonce 并加1

```go
var qosnonce int64 = 0
acc, err := config.GetCLIContext().QOSCliContext.GetAccount(key, cdc)
if err != nil {
qosnonce = 0
} else {
qosnonce = int64(acc.Nonce)
}
qosnonce++
````

## 封装QOS主链转账交易

QOS公链转账交易需要提供 

form 出账账户地址

to   收款账户地址

ccs  转账币种及数量


```go
var ccs []qbasetypes.BaseCoin
for _, coin := range coins {
  ccs = append(ccs, qbasetypes.BaseCoin{
  Name:   coin.Denom,
  Amount: qbasetypes.NewInt(coin.Amount.Int64()),
 })
}

t := tx.NewTransfer(from, to, ccs)
````


### 签名及广播


```go
cliCtx = *config.GetCLIContext().QSCCliContext
result, err1, qscnonce := queryQSCAccount(cdc, key)
if result != nil {
  return result, err1
}
qscnonce++
order := &OrderTx{Address: from, OrderTo: to, OrderAmount: ccs[0].Amount, Gas: qbasetypes.NewInt(0)}
msg = genStdWrapTx(cdc, t, priv, config.GetCLIContext().Config.QOSChainID, config.GetCLIContext().Config.QSCChainID, qosnonce, qscnonce, order)

````

先提交到联盟链  再由联盟链转到QOS主链上 双层链执行交易


[stub.go    负责合约启动 注册mapper及交易结构体等 监听跨链结果](stub.md) 