---
概述: QOS联盟链合约示例
---

现在一切就绪我们就可以开始创建合约了，在qstars工程中的x包下新建一个example示例目录
然后在example眼浅上新建supply目录.这就是存放本次合约示例的目录了
在supply中创建 cmd.go handler.go mapper.go process.go stub.go 五个文件 
目录结构如下

	qstars
	|-- x
	|   |-- example
	|          |--supply
	|          `-- cmd.go
	|          `-- handler.go    
	|          `-- mapper.go    
	|          `-- process.go    
    |          `-- stub.go              


下面介绍每个go文件的作用 

[hanlder.go 合约主逻辑](handler.md) 

[mapper.go  数据存取(lvevl db)逻辑](mapper.md) 

[process.go 调用合约SDK 负责签名和广播交易](process.md) 

[stub.go    负责合约启动 注册mapper及交易结构体等 监听跨链结果](stub.md) 

[cmd.go    命令行方式与链码交互入口](cmd.md)