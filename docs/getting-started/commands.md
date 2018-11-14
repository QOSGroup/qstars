
# 命令行
## 查账户余额
```
./qstarscli account address1k0m8ucnqug974maa6g36zw7g2wvfd4sug6uxay
```

## 转账
```
./qstarscli send --from=rpt3O80wAFI1+ZqNYt8DqJ5PaQ+foDq7G/InFfycoFYT8tgGFJLp+BSVELW2fTQNGZ/yTzTIXbu9fg33gOmmzA== --to=address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355 --amount=1qos

```

## 预授权
```
./qstarscli approve --command=create  --from=rpt3O80wAFI1+ZqNYt8DqJ5PaQ+foDq7G/InFfycoFYT8tgGFJLp+BSVELW2fTQNGZ/yTzTIXbu9fg33gOmmzA== --to=address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355 --amount=2qos

./qstarscli approve --command=increase  --from=rpt3O80wAFI1+ZqNYt8DqJ5PaQ+foDq7G/InFfycoFYT8tgGFJLp+BSVELW2fTQNGZ/yTzTIXbu9fg33gOmmzA== --to=address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355 --amount=3qos

./qstarscli approve --command=decrease  --from=rpt3O80wAFI1+ZqNYt8DqJ5PaQ+foDq7G/InFfycoFYT8tgGFJLp+BSVELW2fTQNGZ/yTzTIXbu9fg33gOmmzA== --to=address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355 --amount=1qos

./qstarscli approve --command=cancel  --from=rpt3O80wAFI1+ZqNYt8DqJ5PaQ+foDq7G/InFfycoFYT8tgGFJLp+BSVELW2fTQNGZ/yTzTIXbu9fg33gOmmzA== --to=address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355 --amount=1qos

 ./qstarscli approve --command=use  --to=0xa328891040df53e54f6a7a5080a357addd1bb361dcc87d5b46feec500453e28e031c8d6d9dbdaf0a8dcf14099503e1e7db59a5b4b511d213e9313937a2fbdac7bd0547bd8b --from=address1k0m8ucnqug974maa6g36zw7g2wvfd4sug6uxay --amount=1qos
```

查询跨链交易seq:
./transfer -m=qcpseq -fromchain=qstars-test
查询具体跨链交易:
./transfer -m=qcpquery -fromchain=qstars-test -qcpseq=3
