
##QStars编译
编译qstars/cmd/qstarsd/main.go为qstarsd
编译qstars/cmd/qstarscli/main.go为qstarscli



##QStar初始化

执行初始化命令 ~/.qstarsd/config/qstarsconf.toml

````cgo
./qstarsd init --moniker supply --chain-id supply-001
````
 - `--moniker`   联盟链别名
 - `--chain-id`  联盟链ID 

输出的结果 就是创世块文件 `~/.qstarsd/config/genesis.json` 

```cgo
{
  "pubKey": "lBuYIKH/IBYQjvNUL9IPihBO4EjzgLprkMEuIkAbsTk=",
  "privKey": "CiYeXHzSXCvbENWQucKJifrK+XU9AJne9vkPCJhKuz6UG5ggof8gFhCO81Qv0g+KEE7gSPOAumuQwS4iQBuxOQ==",
  "addr": "address1mj63dehcraekn9afpfdfsm6ny2juv9v0hywul4",
  "mnemonic": "model silent gloom age diagram hidden zoo walnut truth vote option marine hybrid ice receive vast attract fix hunt lounge castle evidence unit legend",
  "type": "local"
}
{
 "moniker": "supply",
 "chain_id": "supply-001",
 "node_id": "c18b97b29d8a47d121d8e6514117163ad763d621",
 "gentxs_dir": "",
 "app_message": {
  "qcps": [
   {
    "name": "qos",
    "chain_id": "qos",
    "pub_key": {
     "type": "tendermint/PubKeyEd25519",
     "value": "ish2+qpPsoHxf7m+uwi8FOAWw6iMaDZgLKl1la4yMAs="
    }
   }
  ],
  "accounts": [
   {
    "address": "address1mj63dehcraekn9afpfdfsm6ny2juv9v0hywul4",
    "coins": [
     {
      "coin_name": "qstar",
      "amount": "100000000"
     }
    ]
   }
  ]
 }
}

```

### 配置 qsc私钥

在~/.qstarsd/config/下创建qstarsconf.toml文件并输入以下内容 

```
// this is qstars_privatekey
QStarsPrivateKey = "rpt3O80wAFI1+ZqNYt8DqJ5PaQ+foDq7G/InFfycoFYT8tgGFJLp+BSVELW2fTQNGZ/yTzTIXbu9fg33gOmmzA=="
QOSChainName = "qosrace"
```


###配置genesis.json
修改~/.qstarsd/config/genesis.json中qcp标签下的内容
-name和chain_id是qos公共链的名称
-pub_key的值是跨链中继cassin的公钥


###启动QStars

运行 ./qstarsd start 启动
