package main

// 实现思路
// 1. 监听一个tcp端口, 用于处理socks5的连接
// 2. 启动一个goroutine来处理每个连接
// 测试方法
// 1. curl -v --socks5-hostname 172.17.8.219:10801 "https://cn.bing.com/"
// 其中--socks5-hostname参数指定使用socks5代理
func main() {
	go startSocks5Listen("10801")
	go startHttpListen("10802")
	select {}
}
