# Introduction

# P.S.

* local update: `git pull vnt master:vnt`

* local push: `git push vnt vnt:master`


```go
log.Info("[yhx-info] setBootstrapNodes()", "urls", urls, "and tests", tests, "and url length", len(urls))
```

* `libp2p` database tips：

    * when `libp2p` store `key-value`, it will make `key` to `hash` of `key` and the `key` + `value` + `timestamp` will be `value`.

    * when `libp2p` store `key-value`, the `key` should start with `\`.

# develop tips.

* when get the url of bootnodes，we need change it to node，the code is at `cmd/utils/flags.go#L610`

```go
url = "vnode://123b13dfdb555c69290acf510f2e1c00b9d31a917d8c9e7cf8216812da10caffd6e242879bbae43aa09c3bd2ef49c45999bb0c502d38d8c8f9d60287f2220c0c@127.0.0.1:30301"

node, err := discover.ParseNode(url)
```

# VNT starts

## `bootnode` start

```bash
$ cd bootnode-dir
$ bootnode -genkey=node.key
$ bootnode -nodekey=node.key
```

> the database is ：`bootnode -datadir="./" -nodekey=node.key`。

## `membernode` starts 

> ATTENTION：
> * `datadir` can't be too long (https://github.com/vntchain/go-vnt/issues/16342)。
> * `vntbootnode` is `vnt's``bootnodeURL`
> * `port` is `vntnode's` port

```bash
$ gvnt --datadir=./datadir1 --vntbootnode=/ip4/127.0.0.1/tcp/30301/ipfs/QmW1zhpCHrfoyXjWRkJMaTgtfy7BiqhZfHajgK3Xnysoxx --port 30306
```

miners node starts

```bash
$ gvnt account new --datadir ./datadir2

$ gvnt --datadir=./datadir2 --vntbootnode=/ip4/127.0.0.1/tcp/30301/ipfs/QmW1zhpCHrfoyXjWRkJMaTgtfy7BiqhZfHajgK3Xnysoxx --mine --minerthreads=1 --etherbase=0xf6f5038a406a7fe78229a80850ca8ed42fe03bfd --port 30307
```

now，`vntdb` is support `--datadir` tag





# whisper of VNT

NO.1 `gvnt` start

```bash
$ gvnt --datadir dir1
```

NO.1 `wnode` start

```bash
$ wnode -topic=70a4beef -verbosity=4 -ip=:30304 
Please enter the peer's vnode: /ip4/127.0.0.1/tcp/30303/ipfs/1kHJWBz9NHQdMZt1ZmdFwknvNezNQYhcAHq6Fx3nDkNbDf6
```

NO.2 `wnode` start

```bash
$ wnode -topic=70a4beef -verbosity=4 -ip=:30305 
Please enter the peer's vnode: /ip4/127.0.0.1/tcp/30303/ipfs/1kHJWBz9NHQdMZt1ZmdFwknvNezNQYhcAHq6Fx3nDkNbDf6
```

# bzz of VNT

create a vnt account in `datadir2`

```bash
$ gvnt --datadir dir2 account new
```

start a `gvnt` node in `datadir1`

```bash
$ gvnt --datadir dir1
```

start `swarm` node：

```bash
$ swarm --bzzaccount bea9faa39f67da4580c09a65af9521048a52b8f9 --datadir=dir2 --swap-api=dir1/gvnt.ipc --nodiscover
```

get `127.0.0.1:8500` in browser

upload file in swarm

```bash
$ swarm up genesis.json
6a5694e49f29ecb8c53f3392a1ada8c6a2838e5d9372e24816afd20bc51725fd
```

input `6a5694e49f29ecb8c53f3392a1ada8c6a2838e5d9372e24816afd20bc51725fd` to find file
