package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

type httpProxy struct{}

func (p *httpProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	log.Printf("Received request %s %s %s\n", req.Method, req.Host, req.RemoteAddr)

	transport := http.DefaultTransport

	// 1.拷贝源请求 避免修改源请求
	outReq := new(http.Request)
	*outReq = *req

	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		// 将客户端 IP 追加到 X-Forwarded-For 后面
		if prior, ok := outReq.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		outReq.Header.Set("X-Forwarded-For", clientIP)
	}

	// 2. 转发请求到目标服务器
	res, err := transport.RoundTrip(outReq)
	if err != nil {
		rw.WriteHeader(http.StatusBadGateway)
		return
	}

	// 3. 遍历目标服务器返回的响应头的字段信息
	// 并复制添加到客户端响应的头信息中
	for key, value := range res.Header {
		for _, v := range value {
			log.Printf("Header key:%v, value:%v\n", key, v)
			rw.Header().Add(key, v)
		}
	}

	// 4. 将目标服务器返回的状态码写入客户端响应
	// 当写入状态信息后, 无法再通过Set方法修改头信息
	rw.WriteHeader(res.StatusCode)

	// 5. 写入响应体
	io.Copy(rw, res.Body)

	res.Body.Close()
}

func startHttpListen(port string) {
	listen_addr := fmt.Sprintf(":%s", port)
	log.Printf("Http proxy listen addr:%s\n", listen_addr)
	http.Handle("/", &httpProxy{})
	http.ListenAndServe(listen_addr, nil)
}
