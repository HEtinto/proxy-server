package main

import (
	"fmt"
	"log"
	"net"

	"github.com/miekg/dns"
)

func StartDNSServer() {
	dns.HandleFunc(".", handleDNSRequest)

	log.Println("Starting dns server on :53")
	server := &dns.Server{Addr: ":53", Net: "udp"}
	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Failed to start server: %s\n", err.Error())
	}
}

func QueryAddrs(domain string) ([]string, error) {
	addrs, err := net.LookupHost(domain)
	if err != nil {
		return []string{}, err
	}
	if len(addrs) == 0 {
		return []string{}, fmt.Errorf("address can not found for %s", domain)
	}
	return addrs, nil
}

func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false
	m.Authoritative = true // 设置权威响应标志

	switch r.Opcode {
	case dns.OpcodeQuery:
		for _, q := range r.Question {
			switch q.Qtype {
			case dns.TypeA:
				log.Printf("Query for A record: %s\n", q.Name)
				addrs, err := QueryAddrs(q.Name)
				if err != nil {
					log.Printf("Domain query address error: %v\n", err)
					dns.HandleFailed(w, r)
					return
				}
				log.Printf("Domain query address result: %s\n", addrs)
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, addrs[0]))
				if err != nil {
					log.Printf("Failed to create A record: %v\n", err)
					dns.HandleFailed(w, r)
					return
				}
				m.Answer = append(m.Answer, rr)

			case dns.TypeMX:
				log.Printf("Query for MX record: %s\n", q.Name)
				rr, err := dns.NewRR(fmt.Sprintf("%s MX 10 mail.%s", q.Name, q.Name))
				if err != nil {
					log.Printf("Failed to create MX record: %v\n", err)
					dns.HandleFailed(w, r)
					return
				}
				m.Answer = append(m.Answer, rr)

			case dns.TypeCNAME:
				log.Printf("Query for CNAME record: %s\n", q.Name)
				rr, err := dns.NewRR(fmt.Sprintf("%s CNAME target.%s", q.Name, q.Name))
				if err != nil {
					log.Printf("Failed to create CNAME record: %v\n", err)
					dns.HandleFailed(w, r)
					return
				}
				m.Answer = append(m.Answer, rr)

			default:
				log.Printf("Unsupported query type: %d for %s\n", q.Qtype, q.Name)
				// 对于不支持的类型，返回NOERROR状态但不包含任何记录
			}
		}
	}

	// 如果没有答案，设置响应码为NOERROR
	if len(m.Answer) == 0 {
		m.Rcode = dns.RcodeSuccess
	}

	if err := w.WriteMsg(m); err != nil {
		log.Printf("Failed to write message: %v\n", err)
	}
}
