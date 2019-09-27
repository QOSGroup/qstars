### 联盟链（qcp）

QOS跨链协议QCP，支持跨链交易

创建联盟链的前提准备[CA证书](https://github.com/QOSGroup/qos/blob/master/docs/spec/ca.md)

创建联盟链我们需要申请QCP证书[点击了解更多](https://github.com/QOSGroup/qos/blob/master/docs/spec/qcp.md)

QCP证书申请需要提供以下信息：
- 证书适用链
- 联盟链ChainId
- 证书公钥，可通过`kepler genkey`命令行工具生成公私钥对，请妥善保存私钥文件，防止泄露。
- 个人身份证正反面照片或企业营业执照照片


> 创建联盟链需要用到 qoscli发起命令[点击下载](https://github.com/QOSGroup/qos/blob/master/DOWNLOAD.md) 


联盟链相关指令：
* `qoscli tx init-qcp`: [初始化联盟链](#初始化联盟链)
* `qoscli query qcp`:   [查询qcp信息](#查询联盟链)

#### 初始化联盟链

`qoscli tx init-qcp  --node<ip:port> --creator <key_name_or_account_address> --qcp.crt <qcp.crt_file_path>`

主要参数：

- `--node`          QOS主链中任意全节点地址和端口
- `--creator`       创建账号
- `--qcp.crt`       证书位置

> 假设`Arya`已在CA中心申请`qcp.crt`证书，`qcp.crt`中联盟链ID为`aoe-1000`

`Arya`在QOS网络中初始化联盟链信息：
```bash
$ qoscli tx init-qcp --creator Arya --qcp.crt qcp.crt
Password to sign with 'Arya':<输入Arya本地密钥库密码>
```

执行结果：
```bash
{"check_tx":{},"deliver_tx":{},"hash":"BA45F8416780C76468C925E34372B05F5A7FEAAC","height":"243"}
```

#### 查询联盟链

跨链协议是[qbase](https://www.github.com/QOSGroup/qbase)提供支持，主要有以下四个查询指令：
- `qoscli query qcp list`
- `qoscli query qcp out` 
- `qoscli query qcp in`
- `qoscli query qcp tx`

指令说明请参照[qbase-Qcp](https://github.com/QOSGroup/qbase/blob/master/docs/client/command.md#Qcp)。


