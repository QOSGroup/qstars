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