# How to use command line

###send a key and value:

```
./qstarscli kvset --key=bob --value=bbbb --private=1111 --sequence=1 --chain-id=test-chain-AE4XQo
Committed at block 50158 (tx hash: 0B562F9EB18390FA2A5C185160A204006C81620D)
```

###get a key
```
./qstarscli kvget --key=bob
bbbb
```
###create an account
using '#' to seperate base64 privatekey publickey and address
```$xslt
./qstarscli createaccount
XoAKH85xzK+Yx8vCg9zJ6Ntezr0u/yMA3gvU/hRUi/iHB1x9ihr106O0QozpdOdPgRk4iOKV4cD/xNrK0HubfQ==#cosmosaccpub1zcjduepqsur4clv2rt6a8ga5g2xwja88f7q3jwygu227rs8lcndv45rmnd7smdstg2#cosmosaccaddr1t77xlsd4pgq2z3g8d8t8564zfnh3jvps28z5x6

```
###query an account
after account is the address of the account
```
./qstarscli account cosmosaccaddr120ws5500u0q8q75k70uetqp2xnysus5t4x9ug9
{
  "QosAccount": {
    "base_account": {
      "account_address": "address1wmrup5xemdxzx29jalp5c98t7mywulg8wg8cc9",
      "public_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "SKf/Idl3KCPZE9Dwkp0H4K6+mjSSB9sIZ8++mmxK+uE="
      },
      "nonce": "0"
    },
    "qos": "5380394853",
    "qsc": [
      {
        "coin_name": "QSC1",
        "amount": "1234"
      },
      {
        "coin_name": "QSC2",
        "amount": "5678"
      }
    ]
  },
  "QCoins": [
    {
      "denom": "QSC1",
      "amount": "1234"
    },
    {
      "denom": "QSC2",
      "amount": "5678"
    }
  ]
}
```

###do a transaction
from is sender privatekey
mount is coins
to is receiver publickey
```
./qstarscli send --from=XoAKH85xzK+Yx8vCg9zJ6Ntezr0u/yMA3gvU/hRUi/iHB1x9ihr106O0QozpdOdPgRk4iOKV4cD/xNrK0HubfQ==# --amount=3QSC1 --to=cosmosaccaddr120ws5500u0q8q75k70uetqp2xnysus5t4x9ug9 --chain-id=test-chain-AE4XQo
Committed at block 50158 (tx hash: 0B562F9EB18390FA2A5C185160A204006C81620D)
```