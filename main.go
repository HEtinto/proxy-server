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
	httpListenAddress   string // http监听地址
	socks5ListenAddress string // socks5监听地址
	dnsListenAddress    string // dns监听地址
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// 定义命令行参数
	flag.StringVar(&httpListenAddress, "http", ":10802", "http listen address")
	flag.StringVar(&socks5ListenAddress, "socks5", ":10801", "socks5 listen address")
	flag.StringVar(&dnsListenAddress, "dns", ":53", "dns listen address")
	// 解析命令行参数
	flag.Parse()
	log.Printf("Init...(Use --help to get more information)\n")
	log.Printf("http listen address:%s\n", httpListenAddress)
	log.Printf("socks5 listen address:%s\n", socks5ListenAddress)
	log.Printf("dns listen address:%s\n", dnsListenAddress)
	log.Printf("Init done.\n")
}

func main() {
	go StartDNSServer()
	go StartSocks5Listen("10801")
	go StartHttpListen("10802")
	select {}
}
