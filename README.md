# 代理

一个简单的代理程序, 支持 http 和 socks5 代理.

## Linux 设置代理

```shell
export http_proxy=http://127.0.0.1:10809
export https_proxy=socks5://172.0.0.1:10808
```

## git 设置代理

```shell
git config --global http.proxy "xxx"
git config --global https.proxy "xxx"
```

## git 设置免密

只需输入一次密码

```shell
git config --global credential.helper store
```
