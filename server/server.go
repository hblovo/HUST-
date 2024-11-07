package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
)

func main() {
	// 加载服务器证书和私钥
	certificate, err := tls.LoadX509KeyPair("certs/server.crt", "certs/server.key")
	if err != nil {
		fmt.Printf("Error loading server certificates: %v\n", err)
		os.Exit(1)
	}

	// 加载客户端根证书 (CA)
	caCert, err := os.ReadFile("certs/ca.crt")
	if err != nil {
		fmt.Printf("Error reading CA certificate: %v\n", err)
		os.Exit(1)
	}

	// 配置TLS
	config := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    x509.NewCertPool(),
		ClientAuth:   tls.RequireAndVerifyClientCert, // 双向认证
	}

	// 将CA证书加入ClientCAs池中
	if ok := config.ClientCAs.AppendCertsFromPEM(caCert); !ok {
		fmt.Println("Failed to append CA certificate")
		os.Exit(1)
	}

	// 启动服务器
	listener, err := tls.Listen("tcp", ":4433", config)
	if err != nil {
		fmt.Printf("Error starting TLS listener: %v\n", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server listening on port 4433")

	// 监听客户端连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Client connected")

	// 发送欢迎消息
	_, err := conn.Write([]byte("Welcome to the secure server!\n"))
	if err != nil {
		fmt.Printf("Error writing to client: %v\n", err)
		return
	}
}
