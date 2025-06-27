package main

import (
	"flag"
	"log"
)

// 实现思路
// 1. 监听一个tcp端口, 用于处理socks5的连接
// 2. 启动一个goroutine来处理每个连接
// 测试方法
// 1. curl -v --socks5-hostname 172.17.8.219:10801 "https://cn.bing.com/"
// 其中--socks5-hostname参数指定使用socks5代理

var (
	HttpListenAddress   string // http监听地址
	Socks5ListenAddress string // socks5监听地址
	DNSListenAddress    string // dns监听地址
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// 定义命令行参数
	flag.StringVar(&HttpListenAddress, "http", ":10802", "http listen address")
	flag.StringVar(&Socks5ListenAddress, "socks5", ":10801", "socks5 listen address")
	flag.StringVar(&DNSListenAddress, "dns", ":53", "dns listen address")
	// 解析命令行参数
	flag.Parse()
	log.Printf("Init...(Use --help to get more information)\n")
	log.Printf("http listen address:%s\n", HttpListenAddress)
	log.Printf("socks5 listen address:%s\n", Socks5ListenAddress)
	log.Printf("dns listen address:%s\n", DNSListenAddress)
	log.Printf("Init done.\n")
}

func main() {
	go StartDNSServer(DNSListenAddress)
	go StartSocks5Listen(Socks5ListenAddress)
	go StartHttpListen(HttpListenAddress)
	select {}
}
