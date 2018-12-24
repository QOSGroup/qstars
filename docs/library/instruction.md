# Instruction on the interfaces in mobile app
How to invoke the functions/interfaces for mobile app

## Current version via commit hash:
```
aa3a2b5
```

## Interfaces exposed  
There are 13 interfaces exposed under the current version, let?s start with the interfaces one by one.

### 1. AccountCreate
Used for account creation with usage:
```
func AccountCreate(password string) string
```
Input `password ` string format for account generation with example output as below:
```
{
          "jsonrpc": "2.0",
          "id": "",
          "result": {
            "pubKey": "YRoVbOkSxltrudPcvZhZF4tbF/293Tozp8s+Pm9EHXk=",
            "privKey": "HDzw0o0BlBMP6+tJXp7U65vpH5/UKaWtHaqLKBlLDD5hGhVs6RLGW2u509y9mFkXi1sX/b3dOjOnyz4+b0QdeQ==",
            "addr": "address16rd2qgs6whuhycgm90yvx8uz6dwh83vh2jk2gn",
            "mnemonic": "siren venue foil chaos hen margin hockey cost decide critic process off allow layer position morning used plunge onion volume job tape into before",
            "type": "local"
          }
 }
```
Note: Write down the mnemonic, it is critical for the account recovery procedure.

### 2. AccountRecover
Used for account recovery with usage:
```
func AccountRecover(mncode, password string) string
```
Input the "mnemonic" as `mncode` and the `password` upon account creation, it could be recovered with example output as below:
```
{
          "jsonrpc": "2.0",
          "id": "",
          "result": {
            "pubKey": "YRoVbOkSxltrudPcvZhZF4tbF/293Tozp8s+Pm9EHXk=",
            "privKey": "HDzw0o0BlBMP6+tJXp7U65vpH5/UKaWtHaqLKBlLDD5hGhVs6RLGW2u509y9mFkXi1sX/b3dOjOnyz4+b0QdeQ==",
            "addr": "address16rd2qgs6whuhycgm90yvx8uz6dwh83vh2jk2gn",
            "mnemonic": "siren venue foil chaos hen margin hockey cost decide critic process off allow layer position morning used plunge onion volume job tape into before",
            "type": "local"
          }
 }
```
Obviously, it is the same format as that while account creation!

### 3. PubAddrRetrieval
Used for public key and address retrieval from private key input with usage:
```
func PubAddrRetrieval(priv string) string
```
Input the "privKey" of the account as `priv` to fetch the public key and address information, the output could be as below:
```
{
          "jsonrpc": "2.0",
          "id": "",
          "result": {
            "pubKey": "YRoVbOkSxltrudPcvZhZF4tbF/293Tozp8s+Pm9EHXk=",
            "addr": "address16rd2qgs6whuhycgm90yvx8uz6dwh83vh2jk2gn"
          }
 }
```

### 4. AesEncrypt
Used for the private key AES encryption alongside with usage:
```
func AesEncrypt(key, plainText string) string 
```
Input the "key" as the string output of internal HMAC, the byte size should be 16, and the "plainText" is the "privKey" of the account. The output could be as below:
```
90wz2NtXedwF5kMnjF88kty4yjxacijK9GdhpZhOvYwaeB-UV1QA6E0uuAv2rDUSlW3v81OS-u4I5sj2nBwF-vIxsyDc3UxwyIezFJPTHYo1J2yaSFX77WzxaxhOCXAQWTowgSO-aPs=
```
It would be stored in the app equipment and retrieved for the "plain" private key after decryption.

Note: It is worth mentioning however that there is a "salt" of sorts for AES as well in the form of the IV (initialization vector). This is simply random data (which should be different every time), that can be stored alongside the ciphertext, and ensures that even when encrypting the same data multiple times, the output is different.


### 5. AesDecrypt
Used for the private key retrieval on AES decryption with usage:
```
func AesDecrypt(key, cipherText string) string 
```
Input the same "key" as encryption, and the "cipherText" output of the private key after encryption. The "plain" private key could be shown as below:
```
HDzw0o0BlBMP6+tJXp7U65vpH5/UKaWtHaqLKBlLDD5hGhVs6RLGW2u509y9mFkXi1sX/b3dOjOnyz4+b0QdeQ==
```

### 6. SetBlockchainEntrance
Used for set the Restful entrance for the block chain with usage:
```
func SetBlockchainEntrance(sh, mh string)
```
Input "sh" means the restful host of block chain(e.g. qstars), the format could be "IP:1317"; the "mh" is reserved for Qmoon explorer, right now could be input as "forQmoonAddr" or any string.

Note: This function should be invoked every time Restful operations are taking.

### 7. QSCQueryAccount
Used for querying the QSC account information with usage:
```
func QSCQueryAccount(addr string) string
```
Input "addr" of the account you want to query. 

Note: Remember this is only used for query the account activated in QSC other than QOS, error would be returned if input the wrong address.

### 8. QOSQueryAccount
Used for querying the QOS account information with usage:
```
func QOSQueryAccount(addr string) string
```
Input "addr" of the account you want to query. 

Note: Remember this is only used for query the account activated in QOS other than QSC, error would be returned if input the wrong address.

### 9.  QSCKVStoreSet
Used for setting the keystore information in QSC with usage:
```
func QSCKVStoreSet(k, v, privkey, chain string) string
```
Input "k", "v" is the key-value as name, "chain" is the chain-id of the QSC, which is also corresponding to the Restful host with `SetBlockchainEntrance`; for security reason, "privkey" would be deprecated, now you can input any string.

Note: It should not input as plain "privkey" by security reason!

### 10.  QSCKVStoreGet
Used for setting the keystore information in QSC with usage:
```
func QSCKVStoreGet(k string) string
```
Input "k" is the key set in the previous `QSCKVStoreSet`

### 11.  QSCtransferSend
Used for transaction via the QSC then direct to QOS with usage:
```
func QSCtransferSend(addrto, coinstr, privkey, chainid string) string
```

Input "addrto" is the address of the receiver; "coinstr" is the transaction amount, e.g. "1QOS"; "privkey" is the "plain" private key of the sender; "chainid" is the chain-id of the QSC which is also corresponding to the Restful host with `SetBlockchainEntrance`

### 12.  QOSCommitResultCheck
Used for QSC transaction check with usage:
```
func QOSCommitResultCheck(txhash, height string) string
```
Input "txhash" ,"height" are the transaction hash and block height respectively after the transaction result feedback. 
Note: This function is partially implemented while QOS unsupported.

### 13.  JQInvestAd
Used for generation the transaction message with usage:
```
func JQInvestAd(QOSchainId, QSCchainId, articleHash, coins, privatekey string) string 
```
Input "QOSchainId", "QSCchainId" are the chain-id for the QOS and QSC respectively; "articleHash" is the hash of the article which has already created in QSC; "coins" is the investing amount; "privatekey" is the "plain" private key of the Dapp owner.
The output would be a transaction message on json, to be sent to Java background to push to block chain. An output example could be as below:
```
{
    "type": "qbase/txs/stdtx",
    "value": {
      "itx": {
        "type": "qstars/InvestTx",
        "value": {
          "Std": {
            "itx": {
              "type": "qos/txs/TxTransfer",
              "value": {
                "senders": [
                  {
                    "addr": "address13mjc3n3xxj73dhkju9a0dfr4lrfvv3whxqg0dy",
                    "qos": "0",
                    "qscs": [
                      {
                        "coin_name": "AOE",
                        "amount": "5"
                      }
                    ]
                  }
                ],
                "receivers": [
                  {
                    "addr": "address18yunjwfe8yunjwfe8yunjwfe8yunjwfex078fv",
                    "qos": "0",
                    "qscs": [
                      {
                        "coin_name": "AOE",
                        "amount": "5"
                      }
                    ]
                  }
                ]
              }
            },
            "sigature": [
              {
                "pubkey": {
                  "type": "tendermint/PubKeyEd25519",
                  "value": "fDLjEW4zeLVCvKKx4iYB00fnp5Mcl3APIIla7KyETOE="
                },
                "signature": "8BSwrU4aRzeKuaedOmfX0y7vliKtBxSquvMm6EFJU3d+cv8fK+/Vyen7ZqDTUizlTuOr7VP7Ydtw2MiSGSt7Ag==",
                "nonce": "18"
              }
            ],
            "chainid": "qos-test",
            "maxgas": "0"
          },
          "articleHash": "MTU0NTEyODMyMTE0Mg=="
        }
      },
      "sigature": [
        {
          "pubkey": {
            "type": "tendermint/PubKeyEd25519",
            "value": "fDLjEW4zeLVCvKKx4iYB00fnp5Mcl3APIIla7KyETOE="
          },
          "signature": "7/5IhnUZEzEX+nBCqBDBSJzltkCDjgZVYkP8NUc7RZ7m7Gkpb2p5tInzBaqIop8YCiyt9qPOXFW2Xje3/+B/Dg==",
          "nonce": "28"
        }
      ],
      "chainid": "qstars-test",
      "maxgas": "0"
    }
  }
```
