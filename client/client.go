package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
)

func main() {
	// 加载客户端证书和私钥
	certificate, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
	if err != nil {
		log.Fatalf("client: loadkeys: %s", err)
	}

	// 加载根证书 (CA)
	caCert, err := os.ReadFile("certs/ca.crt")
	if err != nil {
		log.Fatalf("client: read ca cert: %s", err)
	}

	// 配置TLS
	config := &tls.Config{
		Certificates:       []tls.Certificate{certificate},
		RootCAs:            x509.NewCertPool(),
		InsecureSkipVerify: true,
	}

	// 将CA证书加入RootCAs池中
	config.RootCAs.AppendCertsFromPEM(caCert)

	// 连接到服务器
	conn, err := tls.Dial("tcp", "localhost:4433", config)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}
	defer conn.Close()

	fmt.Println("Connected to the server!")

	// 读取服务器的消息
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatalf("client: read: %s", err)
	}

	// 输出服务器消息
	fmt.Printf("Server: %s", string(buf[:n]))
}
