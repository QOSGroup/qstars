#Build your own transaction
Qstar can accept message as following. So you could generate a kind of json string and send to chain to run your own transaction.
 
```
{
	"type": "qbase/txs/stdtx",
	"value": {
		"itx": {
			"type": "qos/txs/TxTransfer",
			"value": {
				"senders": [{
					"addr": "address13mjc3n3xxj73dhkju9a0dfr4lrfvv3whxqg0dy",
					"qos": "0",
					"qscs": [{
						"coin_name": "AOE",
						"amount": "2"
					}]
				}],
				"receivers": [{
					"addr": "address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355",
					"qos": "0",
					"qscs": [{
						"coin_name": "AOE",
						"amount": "2"
					}]
				}]
			}
		},
		"sigature": [{
			"pubkey": {
				"type": "tendermint/PubKeyEd25519",
				"value": "fDLjEW4zeLVCvKKx4iYB00fnp5Mcl3APIIla7KyETOE="
			},
			"signature": "6o1JeAA8JuFVaiEa7hg43riTqnRlWzFAMuM9mlFzexrTbKUaFCpWCpKpDPpdV6b78qjV20nzTe4nD/IucgZxDA==",
			"nonce": "46"
		}],
		"chainid": "qos-test",
		"maxgas": "0"
	}
}
```

For example:
```
send --from=Ey+2bNFF2gTUV6skSBgRy3rZwo9nS4Dw0l2WpLrhVvV8MuMRbjN4tUK8orHiJgHTR+enkxyXcA8giVrsrIRM4Q== --to=address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355 --toamount=2AOE --fromamount=2AOE
```

## input
- sender private key
base64: Ey+2bNFF2gTUV6skSBgRy3rZwo9nS4Dw0l2WpLrhVvV8MuMRbjN4tUK8orHiJgHTR+enkxyXcA8giVrsrIRM4Q==
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
2AOE

```

## output
json string

## step
- data is required to sign are following:

```go
Need to signdata hex:   8ee588ce2634bd16ded2e17af6a475f8d2c645d73032414f4557614e5da14a88514adc1995fd3447414e0b20643032414f45716f732d746573740000000000000000000000000000002e716f732d74657374
Need to signdata byte:  [142 229 136 206 38 52 189 22 222 210 225 122 246 164 117 248 210 198 69 215 48 50 65 79 69 87 97 78 93 161 74 136 81 74 220 25 149 253 52 71 65 78 11 32 100 48 50 65 79 69 113 111 115 45 116 101 115 116 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 46 113 111 115 45 116 101 115 116]
signature hex:      ea8d4978003c26e1556a211aee1838deb893aa74655b314032e33d9a51737b1ad36ca51a142a560a92a90cfa5d57a6fbf2a8d5db49f34dee270ff22e7206710c
signature byte:     [234 141 73 120 0 60 38 225 85 106 33 26 238 24 56 222 184 147 170 116 101 91 49 64 50 227 61 154 81 115 123 26 211 108 165 26 20 42 86 10 146 169 12 250 93 87 166 251 242 168 213 219 73 243 77 238 39 15 242 46 114 6 113 12]
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


