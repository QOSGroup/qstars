## Create your Own QStars

To create your own qstars, first each validator will need to install qstarsd and run gen-tx

```bash
qstarsd init 
```

This will generate a `genesis.json` in `$HOME/.qstarsd/config/genesis.json` distribute this file to all validators on your qstars.

### Export state

To export state and reload (useful for testing purposes):

```
qstarsd export > genesis.json; cp genesis.json ~/.qstarsd/config/genesis.json; qstarsd start
```

How to setup a tendmint testnet, please see tendermint

modify ~/.qstarsd/config/genesis.json
1 change the "chain_id" to "qstars-test"
2 add following to genesis.json
```
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
```

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
