
jianqian错误号，分配：

* coins 101-199
```
101  参数长度不一致
102  私钥获取地址错误
103  交易出错
104  查询跨链结果错误
105  查询跨链结果错误
```

* article 201-299
```
201  上传者私钥获取地址错误
202  原作者获取地址错误
203  作者收入比例出错
204  原创入比例出错
205  社区收入比例出错
206  投资者收入比例出错
207  投资期限出错
208  购买期限比例出错
209  地址转换错误
210  非作者本人私钥
211  交易出错
212  查询跨链结果错误
```

* invested 301-399

```
301  文章hash不正确
302  coins不正确
310  上链错误
```

* buyad 401-499
```
401  文章hash不正确
402  coins不正确
403  没有社区地址
404  已经被购买
410  上链错误
```

100以内是跨链的错误号，如果跨链失败了，直接把跨链的错误号返回客户端。

### jianqian application
jianqian is first application of QOS and QStars Chain

qcpTxResult.Result.Log 是跨链的错误原因字符串


超时错误都是负数，

-2 代表未上联盟链超时了，

-1代表已经上主链，未上联盟链超时了。

0 是成功


如果像文章不存在这种没上链的错误，也要返回这个，就是没有hash和heigth而已

###上传新文章：
```
 ./qstarscli NewArticle --authoraddress="address1y9r4pjjnvkmpvw46de8tmwunw4nx4qnz2ax5ux" --originalAuthor="address1zsqzn6wdecyar6c6nzem3e8qss2ws95csr8d0r" --articleHash="814CBF7374D249564ED6220AC837DFC46F5EC82C" --shareAuthor="20" --shareOriginalAuthor="20" --shareCommunity="10" --shareInvestor="50" --endInvestDate="20" --endBuyDate="3"

### 查询文章
```
./qstarscli QueryArticle --articleHash="abcd"
```

### 投资文章广告
```
./qstarscli investad invest --articleHash=abcd --coins=1AOE --investor=maD8NeYMqx6fHWHCiJdkV4/B+tDXFIpY4LX4vhrdmAYIKC67z/lpRje4NAN6FpaMBWuIjhWcYeI5HxMh2nTOQg==
```

### 查询文章投资
```
./qstarscli investad query abcd
```

### 购买文章广告
```
./qstarscli buyad buyad --articleHash=abcd --coins=100QOS --buyer=maD8NeYMqx6fHWHCiJdkV4/B+tDXFIpY4LX4vhrdmAYIKC67z/lpRje4NAN6FpaMBWuIjhWcYeI5HxMh2nTOQg==
```

### 查询文章投资
```
./qstarscli buyad query abcd
```


```
###撒币：
```
 ./qstarscli DispatchAoe --address="address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355|address1zsqzn6wdecyar6c6nzem3e8qss2ws95csr8d0r" --coin="500|400" --causecode="2|3" --causestrings="qiandao|shiming"
```
### 测试账号：

地址：
address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355

私钥：
31PlT2p6UICjV63dG7Nh3Mh9W0b+7FAEU+KOAxyNbZ29rwqNzxQJlQPh59tZpbS1EdIT6TE5N6L72se9BUe9iw==

地址：
address1zsqzn6wdecyar6c6nzem3e8qss2ws95csr8d0r

私钥：
vAeIlHuWjvz/JmyGcB46ZHfCZdXCYuRogqxDgjYUM5wNwKIyIYQBs9VZxGyD9FS5J4XvZntnUaTtoGsEl7+3hg==

发钱地址：
address13mjc3n3xxj73dhkju9a0dfr4lrfvv3whxqg0dy

发钱私钥：

```go
# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

# Path to the JSON file containing the initial validator set and other meta data
qos_chain_id = "qos-test"

qsc_chain_id = "qstars-test"

qos_node_uri = "192.168.1.224:26657"

qstars_node_uri = "192.168.1.223:26657"

direct_to_qos = false

waiting_for_qos_result = 100
community = "9QkouVPl29N2v1lBO1+azUDqm38fAgs6d3Xo8DcnCus7xjMqsavhc190xCGzZuXcjapUahi7Y7v2DD4hzVCAsQ=="
authormock = "9QkouVPl29N2v1lBO1+azUDqm38fAgs6d3Xo8DcnCus7xjMqsavhc190xCGzZuXcjapUahi7Y7v2DD4hzVCAsQ=="
adbuyermock = "9QkouVPl29N2v1lBO1+azUDqm38fAgs6d3Xo8DcnCus7xjMqsavhc190xCGzZuXcjapUahi7Y7v2DD4hzVCAsQ=="
banker = "9QkouVPl29N2v1lBO1+azUDqm38fAgs6d3Xo8DcnCus7xjMqsavhc190xCGzZuXcjapUahi7Y7v2DD4hzVCAsQ=="
dappowner = "Ey+2bNFF2gTUV6skSBgRy3rZwo9nS4Dw0l2WpLrhVvV8MuMRbjN4tUK8orHiJgHTR+enkxyXcA8giVrsrIRM4Q=="

```

```go
"name": "acc3", 
 "addrbech32": "address1y9r4pjjnvkmpvw46de8tmwunw4nx4qnz2ax5ux", 
 "address": "214750CA5365B6163ABA6E4EBDBB9375666A8262", 
 "pub_key": { 
 "type": "tendermint/PubKeyEd25519", 
 "value": "O8YzKrGr4XNfdMQhs2bl3I2qVGoYu2O79gw+Ic1QgLE=" 
 }, 
 "priv_key": { 
 "type": "tendermint/PrivKeyEd25519", 
 "value": "9QkouVPl29N2v1lBO1+azUDqm38fAgs6d3Xo8DcnCus7xjMqsavhc190xCGzZuXcjapUahi7Y7v2DD4hzVCAsQ==" 
 } 
 } 
```