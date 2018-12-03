## Create your Own single node QStars
prepare your chainid and coin name (equal to chain name)

### build qstarsd binary
To create your own qstars, first each validator will need to install qstarsd in installation.md

### apply CA certification public key and private key
you will get qsc key pair and relay key pair

### Ask QOS public chain administrator to create your alien chain on QOS
see qos/docs/txcreateqsc_txissue_test.md
see QOS documment:

[QOS installation](https://github.com/QOSGroup/qos/blob/master/docs/install/networks.md)

[create QSC on QOS](https://github.com/QOSGroup/qos/blob/master/docs/client/qsc.md)
### setup a your relay server
see xxx
it will requires the relay key pair


### Build qstarsd binary
To create your own qstars, first each validator will need to install qstarsd in installation.md

### Apply CA certification public key and private key
you will get qsc key pair and relay key pair
- alien needs: qsc private key and relay public key
- relay needs: relay private key
- qos public chain needs: qsc public key

### Setup a your relay server
see xxx
it will requires the relay key pair

### Initial your qstars
```bash
qstarsd init 
```
This will generate a `genesis.json` in `$HOME/.qstarsd/config/genesis.json` 

### Configure qsc private key to following file
create  ~/.qstarsd/config/qstarsconf.toml and add following content.
```
// this is qstars_privatekey
QStarsPrivateKey = "rpt3O80wAFI1+ZqNYt8DqJ5PaQ+foDq7G/InFfycoFYT8tgGFJLp+BSVELW2fTQNGZ/yTzTIXbu9fg33gOmmzA=="
QOSChainName = "qosrace"
```

### Configure genesis.json
modify ~/.qstarsd/config/genesis.json
1. change the "chain_id" to "your-chain-id"
2. add following to genesis.json

- the name and chain_id is the qos public chain's name 
- the pub_key's value is cassini's public key


### Configure qstars client configuration file
```
[vagrant@vagrant-192-168-1-223 qstarsd]$ cat ~/.qstarscli/config/config.toml 
```
- qos_chain_id is public chain id
- qsc_chain_id is your alien chain id
- qos_node_uri is public chain abci url
- qsc_node_uri is your alien chain abci url

###Start QStars
./qstarsd start --with-tendermint
 or 
you can run startup.sh 

### Configure command line tool and run commands.
See [commands](./commands.md)

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

## Create your Own Multiple node QStars
download or build a tendermint of version 0.23.1
```
[tendermint]$ ./tendermint version
0.23.1
```
run tendermint testnet to create each node's configuration file
```cgo
[vagrant@vagrant-192-168-168-191 tendermint]$ ./tendermint testnet --n 4 --hostname-prefix qosracenode --starting-ip-address 172.19.222.64 --v 4
I[11-14|06:12:08.877] Found private validator                      module=main path=mytestnet/node0/config/priv_validator.json
I[11-14|06:12:08.877] Found node key                               module=main path=mytestnet/node0/config/node_key.json
I[11-14|06:12:08.877] Found genesis file                           module=main path=mytestnet/node0/config/genesis.json
I[11-14|06:12:08.877] Found private validator                      module=main path=mytestnet/node1/config/priv_validator.json
I[11-14|06:12:08.877] Found node key                               module=main path=mytestnet/node1/config/node_key.json
I[11-14|06:12:08.877] Found genesis file                           module=main path=mytestnet/node1/config/genesis.json
I[11-14|06:12:08.878] Found private validator                      module=main path=mytestnet/node2/config/priv_validator.json
I[11-14|06:12:08.878] Found node key                               module=main path=mytestnet/node2/config/node_key.json
I[11-14|06:12:08.878] Found genesis file                           module=main path=mytestnet/node2/config/genesis.json
I[11-14|06:12:08.878] Found private validator                      module=main path=mytestnet/node3/config/priv_validator.json
I[11-14|06:12:08.878] Found node key                               module=main path=mytestnet/node3/config/node_key.json
I[11-14|06:12:08.878] Found genesis file                           module=main path=mytestnet/node3/config/genesis.json
I[11-14|06:12:08.880] Generated private validator                  module=main path=mytestnet/node4/config/priv_validator.json
I[11-14|06:12:08.881] Generated node key                           module=main path=mytestnet/node4/config/node_key.json
I[11-14|06:12:08.881] Generated genesis file                       module=main path=mytestnet/node4/config/genesis.json
I[11-14|06:12:08.883] Generated private validator                  module=main path=mytestnet/node5/config/priv_validator.json
I[11-14|06:12:08.883] Generated node key                           module=main path=mytestnet/node5/config/node_key.json
I[11-14|06:12:08.883] Generated genesis file                       module=main path=mytestnet/node5/config/genesis.json
I[11-14|06:12:08.885] Generated private validator                  module=main path=mytestnet/node6/config/priv_validator.json
I[11-14|06:12:08.886] Generated node key                           module=main path=mytestnet/node6/config/node_key.json
I[11-14|06:12:08.886] Generated genesis file                       module=main path=mytestnet/node6/config/genesis.json
I[11-14|06:12:08.888] Generated private validator                  module=main path=mytestnet/node7/config/priv_validator.json
I[11-14|06:12:08.889] Generated node key                           module=main path=mytestnet/node7/config/node_key.json
I[11-14|06:12:08.889] Generated genesis file                       module=main path=mytestnet/node7/config/genesis.json
Successfully initialized 8 node directories

```
you will get several folders which has configuration files

### Manually create qstarsserverconf.toml file in every folder
This is as above step

### change genesis files
This is as above step

### startup every node
Copy each folder to each node