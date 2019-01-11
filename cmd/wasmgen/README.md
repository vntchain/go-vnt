# wasmgen工具安装

## 依赖clang 5.0 llvm 5.0

### ``ubuntu``

```
wget https://raw.githubusercontent.com/go-clang/gen/master/scripts/switch-clang-version.sh

sh ./switch-clang-version.sh 5.0

sudo ln -s /usr/lib/llvm-5.0/lib/libclang*so /usr/lib/

```

### ``mac``

```
brew install llvm@5

sudo ln -s /usr/local/opt/llvm@5/lib/libclang*dylib /usr/local/lib

```

### ``centos``

#### 1. 添加yum源

```
[alonid-llvm-5.0.0]
name=Copr repo for llvm-5.0.0 owned by alonid
baseurl=https://copr-be.cloud.fedoraproject.org/results/alonid/llvm-5.0.0/epel-7-$basearch/
type=rpm-md
skip_if_unavailable=True
gpgcheck=1
gpgkey=https://copr-be.cloud.fedoraproject.org/results/alonid/llvm-5.0.0/pubkey.gpg
repo_gpgcheck=0
enabled=1
enabled_metadata=1
```

#### 2. 更新源
``yum makecache``

#### 3. 安装llvm
``yum install llvm``


C语言合约完成后通过wasmgen工具生成abi文件和预编译代码precompile.c,文件输出在合约代码目录的output文件夹中

```
./wasmgen --code codepath
```

# clang合约编译成wasm

使用wasm在线编译工具**webassembly studio**进行在线编译，将wasmgen生成的预编译代码precompile.c copy到webassembly studio中，点击build按钮进行编译，编译完成后下载wasm

[webassembly studio 网页链接](https://webassembly.studio/)
