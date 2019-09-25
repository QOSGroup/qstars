---
概述: 执行交易
---


##发起交易

通过qstarcli执行

````go
./qstarscli NewOrder --orderfrom=<订单发起人的私钥地址> --orderto=<订单接收人的账户地址> --orderid=<订单ID> --orderamount=<订单金额>
	
````

- `--orderfrom`    订单发起人的私钥地址
- `--orderto`      订单接收人的账户地址
- `--orderid`      订单ID
 - `--orderamount` 订单金额
