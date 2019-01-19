# wasmgen

c语言智能合约通过``wasmgen``生成所需的``wasm``文件和``abi``文件，已经对合约进行预编译的代码，通过``wasm``和``abi``文件可以将智能合约部署到区块链网络并进行访问

## 依赖clang 5.0 llvm 5.0

### ``ubuntu``

```
wget https://raw.githubusercontent.com/go-clang/gen/master/scripts/switch-clang-version.sh

/bin/bash ./switch-clang-version.sh 5.0

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


## 下载wasmception

### ```mac```

[wasmception]https://github.com/ooozws/clang-heroku-slug/blob/master/precomp/wasmception-darwin-bin.tar.gz

### ```linux```

[wasmception]https://github.com/ooozws/clang-heroku-slug/blob/master/precomp/wasmception-linux-bin.tar.gz

下载wasmception并解压wasmception,设置wasmception的环境变量

```
echo export VNT_WASMCEPTION="/[PATH]/wasmception-[XXX]-bin" >> ~/.bash_profile
source ~/.bash_profile
``` 

wasmgen参数说明

* ``-I``:添加合约引用的头文件所在文件夹,默认为合约代码的文件夹
* ``-code``:合约代码的路径
* ``-output``:wasm，abi和预编译代码输出文件夹，默认路径为在合约代码的文件夹下新建output


# clang合约在线编译

使用wasm在线编译工具**webassembly studio**进行在线编译，将wasmgen生成的预编译代码precompile.c复制到webassembly studio中，点击build按钮进行编译，编译完成后下载wasm

[webassembly studio 网页链接](https://webassembly.studio/)


