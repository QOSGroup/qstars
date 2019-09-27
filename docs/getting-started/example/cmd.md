---
概述: QOS联盟链合约示例
---

命令行方式访问联盟链

##发起一个新订单交易

收集命令行传入的参数并调用process.go中的Send方法对交易进行签名及广播
````go
orderfrom := viper.GetString(flag_orderfrom)
orderto := viper.GetString(flag_orderto)
orderid := viper.GetString(flag_orderid)
orderamount := viper.GetString(flag_orderamount)
to, _ := sdk.AccAddressFromBech32(orderto)
coins, _ := sdk.ParseCoins(orderamount)
result, err := Send(cdc, orderfrom, to, coins, orderid, NewSendOptions(
	gas(viper.GetInt64("gas")),
	fee(viper.GetString("fee"))))

````


##查询一个订单交易

根据订单id查询订单交易信息 

````go
orderid := viper.GetString(flag_orderid)
result, err := config.GetCLIContext().QSCCliContext.QueryStore([]byte(orderid), OrderMapperName)
````

[hanlder.go 合约主逻辑](handler.md) 