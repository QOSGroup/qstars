# How to use restful

## 错误码说明
http 状态码一直为200
业务是否正常根据返回结果json中error对象,具体错误原因详见错误吗

| 状态码   |      说明      |
|----------|:-------------:|
|  |  |
|  |  |
|  |  |


## 响应体格式
```
{
    "jsonrpc": "2.0",
    "id": "",
    "result": {},
    "error": {
        "code": 000,
        "message": "错误信息",
        "data": "错误详细描述"
    }
}
```


## Start restful server
```
./qstarscli rest-server
I[10-09|07:26:56.410] Starting RPC HTTP server on tcp://localhost:1317 module=rest-server 
I[10-09|07:26:56.411] REST server started                          module=rest-server 
I[10-09|07:28:03.382] Served RPC HTTP response                     module=rest-server method=GET url=/accounts/cosmosaccaddr120ws5500u0q8q75k70uetqp2xnysus5t4x9ug9 status=200 duration=3 remoteAddr=127.0.0.1:38350

```

### send a key and value:

```
curl -d "{\"key\":\"1\",\"value\":\"3\",\"privatekey\":\"NQgkKn3vmPgxGut+k832nH6l0A1GDtLy8Nz6fMp6y6q2dbL7FVPkNw9PPcCsobEHJf9eDmC3zrG/iWn9qxLuvg\",\"chainid\":\"chainabc\"}"  http://localhost:1317/kv
```

response
```
{
    "jsonrpc": "2.0",
    "id": "",
    "result": {
        "hash":"D5C48EEA0FCBAB725062797F039636C9C743AE69",
    }
}
```

### get a key
```
curl http://localhost:1317/kv/1
```


response
```
{
    "jsonrpc": "2.0",
    "id": "",
    "result": {
        "value":"3",
    }
}

```



### create an account
using '#' to seperate base64 privatekey publickey and address
```$xslt
curl -d "" http://localhost:1317/accounts

```

response
```
{
    "jsonrpc": "2.0",
    "id": "",
    "result": {
        "pub_key": "Sq+bVJeW6vxDMo73XGcxM8xU7Re36yuLObQyH8+dQkE=",
        "priv_key": "Sq+bVJeW6vxDMo73XGcxM8xU7Re36yuLObQyH8+dQkE=",
        "addr":"cosmosaccaddr120ws5500u0q8q75k70uetqp2xnysus5t4x9ug9",
        "seed": "",
    }
}

```


### query an account
after account is the address of the account
```
curl http://localhost:1317/accounts/cosmosaccaddr120ws5500u0q8q75k70uetqp2xnysus5t4x9ug9
```

response
```
{
    "jsonrpc": "2.0",
    "id": "",
    "result": {
        "QosAccount": {
            "base_account": {
                "account_address": "address13jgd8cvase3gk9zecrqykfjrl0lkahk8kkg3fp",
                "public_key": {
                    "type": "tendermint/PubKeyEd25519",
                    "value": "5PkdAASMl7aDwngycV98nPI0G2XRmLuWaEYl3SWZBSE="
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
}

```

### do a transaction
privatekey in json is sender privatekey
mount in json is coins
the last part of url is receiver publickey
```
curl -d "{\"amount\":\"3QSC1\",\"privatekey\":\"GEPPkslt1Duwnb4B4W8OT1h311LYpo9GuJygHCE6mhH6iq1A17jIzMEzf6NiXUi6iGjDyoj9\",\"chain_id\":\"test-chain-AE4XQo\",\"account_number\":\"2\",\"sequence\":\"1\",\"gas\":\"1\"}" http://localhost:1317/accounts/cosmosaccaddr120ws5500u0q8q75k70uetqp2xnysus5t4x9ug9/send
```

response
```
{
    "jsonrpc": "2.0",
    "id": "",
    "result": {
        "hash":"CB13E9F042F3857DAA78E82BB594BA2366F2D256",
    }
}
```
