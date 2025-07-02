# 代理

一个简单的代理程序, 支持 http 和 socks5 代理, 以及 dns 代理.

## 设置端口

运行时附加 --help 参数查看帮助, 修改http代理端口和socks5代理端口, dns代理端口通常是53, 不建议进行修改.

## 设置代理

### Linux 设置代理

xxx.xxx.xxx.xxx 为代理服务器的 ip 地址.

```shell
export http_proxy=http://xxx.xxx.xxx.xxx:10802
export https_proxy=socks5://xxx.xxx.xxx.xxx:10801
```

### git 设置代理

```shell
git config --global http.proxy "xxx"
git config --global https.proxy "xxx"
```

### git 设置免密

只需输入一次密码

```shell
git config --global credential.helper store
```

### python 配置代理

python 使用 socks5 代理需要安装额外的支持, 在没有安装 pysocks 的环境中配置了 socks5 代理将会导致错误.

安装命令:

```shell
pip install pysocks
```

### 配置 DNS 代理

linux 系统下, 在`/etc/resolv.conf` 文件中添加以下内容(xxx.xxx.xxx.xxx 为代理服务器的 ip 地址):

```shell
nameserver xxx.xxx.xxx.xxx
```
