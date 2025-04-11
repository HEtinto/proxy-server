package main_test

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestCase1(t *testing.T) {
	proxyUrl, err := url.Parse("http://127.0.0.1:10802")
	if err != nil {
		log.Printf("Parse proxy address err:%v\n", err)
		return
	}
	// 创建一个自定义传输配置, 设置代理
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}
	// 创建一个使用自定义传输配置的客户端
	client := &http.Client{
		Transport: transport,
	}
	// 1. 定义url
	url := "http://www.baidu.com/"
	// 2. 创建请求
	req, err := http.NewRequest("Get", url, nil)
	if err != nil {
		log.Printf("Create request err:%v", err)
		return
	}
	// 3.发起请求
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Request err:%v", err)
	}
	// 函数结束时关闭响应体
	defer resp.Body.Close()
	// 读取响应头
	log.Printf("Rsp header:\n")
	for k, v := range resp.Header {
		log.Printf("%v : %v\n", k, v)
	}

	// 读取响应体内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取响应体出错:%v\n", err)
		return
	}

	// 打印响应状态码和内容
	log.Printf("状态码: %d\n", resp.StatusCode)
	log.Printf("响应内容: %s\n", string(body))
}

func TestCase2(t *testing.T) {
	s := []string{"client", "proxy1", "proxy2"}
	s1 := strings.Join(s, ", ") + ", " + "172"
	log.Println(s1)
}
