package main_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"testing"

	"github.com/miekg/dns"
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

func QueryDNSTypeA(domain string) (string, error) {
	c := &dns.Client{}
	m := &dns.Msg{}
	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	m.RecursionDesired = true

	r, _, err := c.Exchange(m, "8.8.8.8:53")
	if err != nil {
		return "", err
	}

	if r.Rcode != dns.RcodeSuccess || len(r.Answer) == 0 {
		return "", fmt.Errorf("no A record found for %s", domain)
	}

	fmt.Printf("r:%+v", r)

	// 类型断言确保记录是 A 类型
	if a, ok := r.Answer[0].(*dns.A); ok {
		return a.A.String(), nil // 返回 IP 地址（如 "93.184.216.34"）
	}

	return r.Answer[0].String(), nil // 兜底：返回原始记录字符串
}

func TestCase3(t *testing.T) {
	// 1. 定义域名
	domain := "www.baidu.com"
	// 2. 调用 QueryDNS 函数 查询
	addr, err := QueryDNSTypeA(domain)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%s resolved to %s\n", domain, addr)
}
