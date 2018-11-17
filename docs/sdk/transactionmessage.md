#Build your own transaction
Qstar can accept message as following. So you could generate a kind of json string and send to chain to run your own transaction.
 
```
{
	"type": "qbase/txs/stdtx",
	"value": {
		"itx": {
			"type": "qos/txs/TransferTx",
			"value": {
				"senders": [{
					"addr": "address1k0m8ucnqug974maa6g36zw7g2wvfd4sug6uxay",
					"qos": "2",
					"qscs": null
				}],
				"receivers": [{
					"addr": "address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355",
					"qos": "2",
					"qscs": null
				}]
			}
		},
		"sigature": [{
			"pubkey": {
				"type": "tendermint/PubKeyEd25519",
				"value": "E/LYBhSS6fgUlRC1tn00DRmf8k80yF27vX4N94Dppsw="
			},
			"signature": "JUTk/5Itlqv7VfjFwvARaEeJiAxfPhT4mCbbMVcF+MzYKkxXuz8f+PYTZeDIQ0W89/uTzBvQpn6Y1J8cyaCeBg==",
			"nonce": "7"
		}],
		"chainid": "qos-test",
		"maxgas": "0"
	}
}
```

For example:
```
send --from=rpt3O80wAFI1+ZqNYt8DqJ5PaQ+foDq7G/InFfycoFYT8tgGFJLp+BSVELW2fTQNGZ/yTzTIXbu9fg33gOmmzA== --to=address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355 --amount=2qos
```

## input
- sender private key
base64: rpt3O80wAFI1+ZqNYt8DqJ5PaQ+foDq7G/InFfycoFYT8tgGFJLp+BSVELW2fTQNGZ/yTzTIXbu9fg33gOmmzA==
Bin:
```
[174 155 119 59 205 48 0 82 53 249 154 141 98 223 3 168 158 79 105 15 159 160 58 187 27 242 39 21 252 156 160 86 19 242 216 6 20 146 233 248 20 149 16 181 182 125 52 13 25 159 242 79 52 200 93 187 189 126 13 247 128 233 166 204]
```
- receiver address
```go
address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355
```

- coin type and mount
```
1qos

```

## output
json string

## step
- data is required to sign are following:
from address + 32 + to address + 32 + chain id + nonce
```
b3f67e6260e20beaefbdd223a13bc8539896d61c3257614e5da14a88514adc1995fd3447414e0b206432716f732d7465737400000000000000000000000000000007
```
go source code
```go
// 签名字节
func (tx TransferTx) GetSignData() (ret []byte) {
   for _, sender := range tx.Senders {
      ret = append(ret, sender.Address...)
      ret = append(ret, (sender.QOS.NilToZero()).String()...)
      ret = append(ret, sender.QSCs.String()...)
   }
   for _, receiver := range tx.Receivers {
      ret = append(ret, receiver.Address...)
      ret = append(ret, (receiver.QOS.NilToZero()).String()...)
      ret = append(ret, receiver.QSCs.String()...)
   }

   return ret
}
```


- from: (Hex format, above is base64)
```
b3f67e6260e20beaefbdd223a13bc8539896d61c
```

- to: (Hex format, above is base64)
```
57614e5da14a88514adc1995fd3447414e0b2064
```

- sign content result is
```
2544e4ff922d96abfb55f8c5c2f011684789880c5f3e14f89826db315705f8ccd82a4c57bb3f1ff8f61365e0c84345bcf7fb93cc1bd0a67e98d49f1cc9a09e06
```


