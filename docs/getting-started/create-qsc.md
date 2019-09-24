# 代币(QSC)发行

QOS公链上支持各联盟链发行自己的代币,代币在联盟链内和业务关联,经cassin到QOS公链结算
##证书
代币发行的前提准备[CA证书](https://github.com/QOSGroup/qos/blob/master/docs/spec/ca.md)


我们共需要申请两个证书  [QSC证收](https://github.com/QOSGroup/qos/blob/master/docs/spec/qsc.md)和[QCP证书](https://github.com/QOSGroup/qos/blob/master/docs/spec/qcp.md)
其中QSC证书是发行代币时使用
而QCP证书是部署跨链中继时使用

QSC证书申请需要提供以下信息：
- 证书适用链
- 联盟币名称，大写字母、长度不超过8个字符
- Banker公钥，可为空。[go-amino](https://github.com/tendermint/go-amino)JSON 序列化ed25519编码信息
- 证书公钥，可通过`kepler genkey`命令行工具生成公私钥对，请妥善保存私钥文件，防止泄露。
- 个人身份证正反面照片或企业营业执照照片



##待续