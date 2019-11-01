
## Hubble Network

[![API Reference](https://camo.githubusercontent.com/915b7be44ada53c290eb157634330494ebe3e30a/68747470733a2f2f676f646f632e6f72672f6769746875622e636f6d2f676f6c616e672f6764646f3f7374617475732e737667)](https://godoc.org/github.com/vntchain/go-vnt) [![Build Status](https://www.travis-ci.org/vntchain/go-vnt.svg?branch=master)](https://www.travis-ci.org/vntchain/go-vnt)

VNT Chain的公链——Hubble Network是下一代商业级智能价值传递的基础设施，致力于解决当前公链高能耗、低效率，难以支撑大规模商业应用的问题。它借鉴了以太坊架构，融合DPOS与BFT技术以获取更高的性能和安全性，服务于规模化商业应用，实现价值的无摩擦转移。基于VNT Chain的公有链，用户可建立和发布下一代分布式智能应用，如游戏、金融资产交易、知识产权交易、电子签约、社交等。

## 特点
1. 采用Vortex双层共识设计，可迅速验证交易及产生区块，高效节约。
2. Vortex具有快速终局性，避免分叉。
3. Vortex可容错1/3的见证人节点。
4. 采用基于WASM的智能合约执行引擎——WAVM，支持多开发语言，降低开发者门槛。
5. WAVM指令精简，可提高交易执行效率。
6. WAVM基于沙箱运行，内存方面具有良好的安全性。
7. WAVM是平台无关的，可运行在各种架构的CPU平台上，且性能等同于原生应用，具有良好的可移植性。
8. WAVM采用模块化设计，可以针对特定功能实现模块化库的开发，提高代码利用率。
9. 采用基于libp2p的P2P网络对已存在节点更加偏好，可以有效抵御恶意新节点的洪泛攻击。
10. P2P网络具有灵活简单特性，可从任意节点发起查找，且查找算法对数收敛，能快速定位内容。
11. P2P子协议框架具有稳定拓展特性，可以在其上任意扩展子协议处理，而底层框架无需做任何改动。


## 从源码安装Gvnt

安装`gvnt`需要Go编译器（版本大于1.9）和C编译器。

首先，克隆仓库`go-vnt`到路径`$GOPATH/src/github.com/vntchain`，并进入项目目录:

    mkdir -p $GOPATH/src/github.com/vntchain
    cd $GOPATH/src/github.com/vntchain
    git clone https://github.com/vntchain/go-vnt
    cd go-vnt
如果有以前的版本，请更新到最新：

    git pull

然后，使用下面命令编译`gvnt`:

    make gvnt

或者使用下面命令编译`gvnt`和所有工具:

    make all

经过以上可以在`$GOPATH/src/github.com/vntchain/go-vnt/build/bin/`目录调用`gvnt`，为方便使用可以将`gvnt`安装到`$GOPATH/bin`:

    go install -v ./cmd/gvnt
 
 查看当前使用的gvnt版本号
 
    gvnt version

## 运行Gvnt

当前VNT的公链主网和测试网络都已经启动，资料参考：[如何加入测试网](https://github.com/vntchain/vnt-documentation/blob/master/developer-guide/05-network/connect-to-hubble-testnet.md), [如何加入主网](https://github.com/vntchain/vnt-documentation/blob/master/developer-guide/05-network/connect-to-hubble-network.md)

你可以在本地搭建一个小型测试网，资料请参考：[如何搭建VNT网络](https://github.com/vntchain/vnt-documentation/blob/master/introduction/set-up-vnt-network/set-up-4-node-vnt-network.md)。

## 资源

1. [VNT Chain官网](http://vntchain.io/)
2. [VNT Chain开发者文档](https://github.com/vntchain/vnt-documentation)
3. [VNT Chain白皮书](https://github.com/vntchain/vnt-documentation/blob/master/VNT-white-paper-CH.pdf)


## 贡献源码

欢迎PR，感谢您为`go-vnt`做的任何一点改进。您可以fork项目到个人仓库后、修复问题进行提交，然后向`go-vnt`仓库发起PR。

贡献代码请遵循以下规则，方便`go-vnt`核心开发人员对代码进行Review。

1. 所有代码经gofmt进行格式化。
2. PR请遵循以下规则：
    1. 标题格式：`[fixed/style/test] #Issue PR标题`，`fixed/style/test`代表了修复/调整格式/修改测试，`#Issue`为本PR相关的Issue编号，PR标题为一句简洁话，描述本次PR的目的。
    1. PR内容：描述本次PR具体的内容，希望能尽可能详细，这样能让Review本PR的开发人员了解你的意图。


## 许可证

所有`go-vnt`仓库生成的二进制程序都采用GNU General Public License v3.0许可证, 具体请查看[COPYING](https://github.com/vntchain/go-vnt/blob/master/COPYING)。
