# 代币(QSC)发行

QOS公链上支持各联盟链发行自己的代币,代币在联盟链内和业务关联,经cassin到QOS公链结算
##证书
代币发行的前提准备[CA证书](https://github.com/QOSGroup/qos/blob/master/docs/spec/ca.md)


发行代币我们需要申请QSC证书[点击了解更多](https://github.com/QOSGroup/qos/blob/master/docs/spec/qsc.md)

QSC证书申请需要提供以下信息：
- 证书适用链
- 联盟币名称，大写字母、长度不超过8个字符
- Banker公钥，可为空。[go-amino](https://github.com/tendermint/go-amino)JSON 序列化ed25519编码信息
- 证书公钥，可通过`kepler genkey`命令行工具生成公私钥对，请妥善保存私钥文件，防止泄露。
- 个人身份证正反面照片或企业营业执照照片



##发行代币

> 创建联盟币前需要用到 qoscli发起命令[点击下载](https://github.com/QOSGroup/qos/blob/master/DOWNLOAD.md) 

联盟币相关指令：
* `qoscli tx create-qsc`    [创建联盟币](#创建联盟币)
* `qoscli query qsc`        [查询联盟币](#查询联盟币)
* `qoscli tx issue-qsc`     [发放联盟币](#发放联盟币)

#### 创建联盟币

`qoscli tx create-qsc --node <ip:port> --creator <key_name_or_account_address> --qsc.crt <qsc.crt_file_path> --accounts <account_qsc_s>`

主要参数：

- `--node`          QOS主链中任意全节点地址和端口
- `--creator`       创建账号
- `--qsc.crt`       证书位置
- `--accounts`      初始发放地址币值集合，[addr1],[amount];[addr2],[amount2],...，该参数可为空，即只创建联盟币

`Arya`在QOS网络中创建`QOE`，不含初始发放地址币值信息：
```bash
$ qoscli tx create-qsc --creator Arya --qsc.crt aoe.crt
Password to sign with 'Arya':<输入Arya本地密钥库密码>
```
> 假设`Arya`已在CA中心申请`aoe.crt`证书，`aoe.crt`中包含`banker`公钥，对应地址`address1rpmtqcexr8m20zpl92llnquhpzdua9stszmhyq`，已经导入到本地私钥库中，名字为`ATM`，。

执行结果：
```bash
{"check_tx":{},"deliver_tx":{},"hash":"BA45F8416780C76468C925E34372B05F5A7FEAAC","height":"200"}
```

#### 查询联盟币

`qoscli query qsc <qsc_name>`

`qsc_name`为联盟币名称

查询`AOE`信息：
```bash
$ qoscli query qsc QOE --indent
```
执行结果：
```bash
{
  "name": "AOE",
  "chain_id": "capricorn-1000",
  "extrate": "1:280.0000",
  "description": "",
  "banker": "address1rpmtqcexr8m20zpl92llnquhpzdua9stszmhyq"
}
```

#### 发放联盟币

针对使用包含`Banker`公钥创建的联盟币，可向`Banker`地址发放（增发）对应联盟币：

`qoscli tx issue-qsc --qsc-name <qsc_name> --banker <key_name_or_account_address> --amount <qsc_amount>`

主要参数：
- `--qsc-name`  联盟币名字
- `--banker`    Banker地址或私钥库中私钥名
- `--amount`    联盟币发放（增发）量

向联盟币AOE `Banker`中发放（增发）10000AOE：

```bash
$ qoscli tx issue-qsc --qsc-name AOE --banker ATM --amount 10000
Password to sign with 'ATM':<输入ATM本地密钥库密码>
```

执行结果：
```bash
{"check_tx":{},"deliver_tx":{},"hash":"BA45F8416780C76468C925E34372B05F5A7FEAAC","height":"223"}
```

- [业务应用实现](example/start.md)
