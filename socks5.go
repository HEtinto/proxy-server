package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
)

// 过程(参考: RFC1928 SOCKS Protocol Version 5)
// (1). 客户端连接服务端, 并发送 版本标识符/认证方法选择数量
// ----------------------
// |VER|NMETHODS|METHODS|
// |1  |1       |1~255  |
// ----------------------
// (2). 服务端设置版本标识和认证方法序号
// -------------
// |VER|METHODS|
// |1  |1      |
// -------------
// (3). 客户端通告服务端它的目标地址
// +---+---+-----+----+--------+--------+
// |VER|CMD|RSV  |ATYP|DST.ADDR|DST.PORT|
// |1  |1  |X`00`|1   |Variable|2       |
// +---+---+-----+----+--------+--------+
// (4). 得到目标地址后, 服务端即可向目标发起连接, 转发目标数据

func StartSocks5Listen(address string) {
	log.Printf("Socks5 proxy listen addr %s\n", address)
	server, err := net.Listen("tcp", address)
	if err != nil {
		log.Printf("Listen failed: %v\n", err)
		return
	}

	for {
		client, err := server.Accept()
		if err != nil {
			log.Printf("Accept failed: %v\n", err)
			continue
		}
		go processSocks5(client)
	}
}

// 实现思路
// 1. 认证
// 2. 建立连接
// 3. 转发数据
func processSocks5(client net.Conn) {
	// 1. 认证
	if err := Socks5Auth(client); err != nil {
		log.Println("auth error:", err)
		client.Close()
		return
	}

	// 2. 建立连接
	target, err := Socks5Connect(client)
	if err != nil {
		log.Println("connect error:", err)
		client.Close()
		return
	}

	// 3. 转发数据
	Socks5Forward(client, target)
}

// RFC 1928
// 客户端首先发送数据: |VER|NMETHODS|METHODS|
// VER: 本次请求的协议版本号, 固定值0x05, 表示socks5
// NMETHODS: 客户端支持的认证方式数量 1~255
// METHODS: 可用的认证方式列表
func Socks5Auth(client net.Conn) (err error) {
	buf := make([]byte, 256)

	// 读取 VER 和 NMETHODS
	n, err := io.ReadFull(client, buf[:2])
	if n != 2 {
		return errors.New("reading header: " + err.Error())
	}

	ver, nMethods := int(buf[0]), int(buf[1])
	if ver != 5 {
		return errors.New("invalid version: " + strconv.Itoa(ver))
	}
	log.Printf("socket ver:%v, NMETHODS:%v\n", ver, nMethods)

	// 读取 METHODS 列表
	n, err = io.ReadFull(client, buf[:nMethods])
	if n != nMethods {
		return errors.New("reading methods: " + err.Error())
	}
	log.Printf("socket methods:%v\n", buf[:nMethods])

	// 设置认证方式 这里使用 00 表示不需要认证
	// * X`00` NO AUTHENTICATION REQUIRED
	// * X`01` GSSAPI
	// * X`02` USERNAME/PASSWORD
	// * X`03` to X`7F` IANA ASSIGNED
	// * X`80` to X`FE` RESERVED FOR PRIVATE METHODS
	// * X`FF` NO ACCEPTABLE METHODS
	n, err = client.Write([]byte{0x05, 0x00})
	if n != 2 || err != nil {
		return errors.New("write rsp: " + err.Error())
	}

	return nil
}

// 完成认证后 客户端要告知服务端它的目标地址
// VER: 协议版本号 0x05
// CMD: 连接方式, 0x01=CONNECT,0X02=BIND, 0X03=UDP ASSOCIATE
// RSV: 保留字段
// ATYP: 地址类型, 0x01=IPv4, 0x03=域名, 0x04=IPv6
// DST.ADDR: 目标地址
// DST.PORT: 目标端口, 2字节, 网络字节序
func Socks5Connect(client net.Conn) (net.Conn, error) {
	buf := make([]byte, 256)

	n, err := io.ReadFull(client, buf[:4])
	if n != 4 {
		return nil, errors.New("read header: " + err.Error())
	}

	ver, cmd, _, atyp := buf[0], buf[1], buf[2], buf[3]
	if ver != 5 || cmd != 1 {
		return nil, errors.New("invalid ver/cmd")
	}
	log.Printf("ver:%v, cmd:%v, atyp:%v\n", ver, cmd, atyp)

	addr := ""
	switch atyp {
	case 1:
		// IPV4 捕获到该值后 直接再读取四个字节即为IPV4地址, 即DST.ADDR
		// 即DST.ADDR = [IPV4地址](占用4个字节)
		n, err = io.ReadFull(client, buf[:4])
		if n != 4 {
			return nil, errors.New("invalid IPv4: " + err.Error())
		}
		addr = fmt.Sprintf("%d.%d.%d.%d", buf[0], buf[1], buf[2], buf[3])
		log.Printf("ATYP IPV4 ADDR:%v\n", addr)

	case 3:
		// 域名 先读取一个字节 该字节通告域名占n字节大小
		// 然后再读取n即为域名
		// 即DST.ADDR = [域名占用字节大小n](1字节) + [域名](n个字节)
		n, err = io.ReadFull(client, buf[:1])
		if n != 1 {
			return nil, errors.New("invalid hostname: " + err.Error())
		}
		addrLen := int(buf[0])

		n, err = io.ReadFull(client, buf[:addrLen])
		if n != addrLen {
			return nil, errors.New("invalid hostname: " + err.Error())
		}
		addr = string(buf[:addrLen])
		log.Printf("Domain Name Addr:%v\n", addr)

	case 4: // IPV6
		return nil, errors.New("IPv6: no supported yet")

	default:
		return nil, errors.New("invalid atyp")
	}

	// 获取2字节大小的端口号
	n, err = io.ReadFull(client, buf[:2])
	if n != 2 {
		return nil, errors.New("read port: " + err.Error())
	}
	port := binary.BigEndian.Uint16(buf[:2])
	log.Printf("port:%v\n", port)

	var destAddrPort string
	if atyp != 4 {
		destAddrPort = fmt.Sprintf("%s:%d", addr, port)
		log.Printf("destAddrPort:%v\n", destAddrPort)
	} else {
		destAddrPort = fmt.Sprintf("[%s]:%d", addr, port)
		log.Printf("destAddrPort:%v\n", destAddrPort)
	}

	// 与目标地址建立连接
	dest, err := net.Dial("tcp", destAddrPort)
	if err != nil {
		return nil, errors.New("dial dst: " + err.Error())
	}

	// 通知客户端本端已准备完毕
	// |VER|REP|RSV|ATYP|BND.ADDR|BND.PORT|
	// 此处以IPV4地址为例 ATYP 0x01
	_, err = client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	if err != nil {
		dest.Close()
		return nil, errors.New("write rsp: " + err.Error())
	}

	return dest, nil
}

func Socks5Forward(client, target net.Conn) {
	forward := func(src, dest net.Conn) {
		defer src.Close()
		defer dest.Close()
		io.Copy(src, dest)
	}
	go forward(client, target)
	go forward(target, client)
}
