# 使用说明

## 1. 下载并安装golang环境

下载地址：https://go.dev/dl/

安装教程：https://go.dev/doc/install

## 2. 配置项目环境
已经把项目的一些配置写进配置文件Toward2017/configure.sh里了，直接运行即可

    cd Toward2017
    source ./configure.sh

## 3. 安装依赖包
    cd src
    go mod tidy

可能出现的问题：

    github.com/Airscope/Toward2017/cloud imports
    github.com/sachaservan/bgn: module github.com/sachaservan/bgn: Get "https://proxy.golang.org/github.com/sachaservan/bgn/@v/list": dial tcp 142.251.42.241:443: i/o timeout

    github.com/Airscope/Toward2017/cloud imports
    github.com/sachaservan/paillier: module github.com/sachaservan/paillier: Get "https://proxy.golang.org/github.com/sachaservan/paillier/@v/list": dial tcp 142.251.42.241:443: i/o timeout

解决方法：
golang官网被墙，更改代理即可

    go env -w GOPROXY=https://goproxy.cn

## 4. 运行代码
在src目录下输入以下代码，运行结果保存在log目录下

    nohup go run . > ../log/log_yyyymmdd &

可能出现的问题：

    # github.com/Nik-U/pbc
    cgo: C compiler "gcc" not found: exec: "gcc": executable file not found in $PATH

    # github.com/ncw/gmp
    cgo: C compiler "gcc" not found: exec: "gcc": executable file not found in $PATH

未安装gcc，安装即可

    sudo apt-get install gcc

可能出现的问题2：

    # github.com/ncw/gmp
    /root/go/pkg/mod/github.com/ncw/gmp@v1.0.4/int.go:14:10: fatal error: gmp.h: No such file or directory
    14  | #include <gmp.h>
        |          ^~~~~~~
    compilation terminated.

    # github.com/Nik-U/pbc
    /root/go/pkg/mod/github.com/!nik-!u/pbc@v0.0.0-20181205041846-3e516ca0c5d6/element.go:25:10: fatal error: pbc/pbc.h: No such file or directory
    25  | #include <pbc/pbc.h>
        |          ^~~~~~~~~~~
    compilation terminated.

未安装pbc和gmp库（pbc依赖gmp），安装即可
- 安装教程：https://blog.csdn.net/qq_41977843/article/details/126765593

## 依赖项目/版本

- [[bgn](https://github.com/sachaservan/bgn/tree/d0ab3643939c7d8fde1194e9ddefa5728b858dbb)] BGN encryption scheme implementation using Go

    commit SHA: 45ab32a771996f67f204e42b94dfd6e60d25a77f

- [[pallier](https://github.com/sachaservan/paillier)] Paillier encryption scheme with threshold decryption and zero-knowledge proof capabilities

    commit SHA:  30237183ba29b1b833f964976a1591af397a529b


## References
Qiu, Shuo, et al. "Toward practical privacy-preserving frequent itemset mining on encrypted cloud data." IEEE Transactions on Cloud Computing 8.1 (2017): 312-323.