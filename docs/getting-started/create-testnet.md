## Create your Own single node QStars
prepare your chainid and coin name (equal to chain name)
### build qstarsd binary
To create your own qstars, first each validator will need to install qstarsd in installation.md

### apply CA certification public key and private key
you will get qsc key pair and relay key pair

### Ask QOS public chain administrator to create your alien chain on QOS
see qos/docs/txcreateqsc_txissue_test.md

### setup a your relay server
see xxx
it will requires the relay key pair

### initial your qstars
```bash
qstarsd init 
```

This will generate a `genesis.json` in `$HOME/.qstarsd/config/genesis.json` distribute this file to all validators on your qstars.

配置 联盟链私钥到下面这个文件
配置 QOS公链
cat ~/.qstarsd/config/qstarsconf.toml 
# this is qstars_privatekey
QStarsPrivateKey = "0xa328891040b7c4ca726ee42e46e0c6cc76f1d68c0e06f9c2894c48289f570dae64d0e05c533b45e7a573d8927e23597c013e01b5c29d5a0b1d2dbae83d6257345870679794"
QOSChainName = "qosrace"

配置genesis.json，增加如下内容
modify ~/.qstarsd/config/genesis.json
1 change the "chain_id" to "your-chain-id"
2 add following to genesis.json

the name and chain_id is the qos public chain's name 
the pub_key's value is cassini's public key
```
 "app_state": {
	"qcps": [
          {
            "name": "qos-test",
            "chain_id": "qos-test",
            "pub_key": {
              "type": "tendermint/PubKeyEd25519",
              "value": "X9NorHcSXnCcEd+7G1ETU66dTiqy7RKCxzlQr37X3WY="
            }
          }
        ]
	}
```

配置config.toml

[vagrant@vagrant-192-168-1-223 qstarsd]$ cat ~/.qstarscli/config/config.toml 
qos_chain_id is public chain id
qsc_chain_id is your alien chain id
qos_node_uri is public chain abci url
qsc_node_uri is your alien chain abci url

```
# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

# Path to the JSON file containing the initial validator set and other meta data
qos_chain_id = "qos-test"

qsc_chain_id = "qstars-test"

qos_node_uri = "192.168.1.224:26657"

qsc_node_uri = "localhost:26657"

direct_to_qos = false

waiting_for_qos_result = 70
```

Genesis example

```
{
  "genesis_time": "2018-10-30T12:20:40.927988421Z",
  "chain_id": "qstars-test",
  "consensus_params": {
    "block_size_params": {
      "max_bytes": "22020096",
      "max_txs": "10000",
      "max_gas": "-1"
    },
    "tx_size_params": {
      "max_bytes": "10240",
      "max_gas": "-1"
    },
    "block_gossip_params": {
      "block_part_size_bytes": "65536"
    },
    "evidence_params": {
      "max_age": "100000"
    }
  },
  "validators": [
    {
      "pub_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "XsVLCKdI6FMhQex5gPzZqunXli8fREkZrOgRj6Lipdg="
      },
      "power": "10",
      "name": ""
    }
  ],
  "app_hash": "",
  "app_state": {
	"qcps": [
          {
            "name": "qos-test",
            "chain_id": "qos-test",
            "pub_key": {
              "type": "tendermint/PubKeyEd25519",
              "value": "DzUh6N/RPVM7kBqq2u3zRbJtQAoP1f1lwk979my/74E="
            }
          }
        ]
	}
}
```
