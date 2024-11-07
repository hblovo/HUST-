package main

import (
	"SSLChat/utils"
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
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

	fmt.Println("oooooo   oooooo     oooo           oooo                                                           .             \n `888.    `888.     .8'            `888                                                         .o8             \n  `888.   .8888.   .8'    .ooooo.   888   .ooooo.   .ooooo.  ooo. .oo.  .oo.    .ooooo.       .o888oo  .ooooo.  \n   `888  .8'`888. .8'    d88' `88b  888  d88' `\"Y8 d88' `88b `888P\"Y88bP\"Y88b  d88' `88b        888   d88' `88b \n    `888.8'  `888.8'     888ooo888  888  888       888   888  888   888   888  888ooo888        888   888   888 \n     `888'    `888'      888    .o  888  888   .o8 888   888  888   888   888  888    .o        888 . 888   888 \n      `8'      `8'       `Y8bod8P' o888o `Y8bod8P' `Y8bod8P' o888o o888o o888o `Y8bod8P'        \"888\" `Y8bod8P' ")
	fmt.Println("  .oooooo.   oooo                      .   ooooooooo.                                         \n d8P'  `Y8b  `888                    .o8   `888   `Y88.                                       \n888           888 .oo.    .oooo.   .o888oo  888   .d88'  .ooooo.   .ooooo.  ooo. .oo.  .oo.   \n888           888P\"Y88b  `P  )88b    888    888ooo88P'  d88' `88b d88' `88b `888P\"Y88bP\"Y88b  \n888           888   888   .oP\"888    888    888`88b.    888   888 888   888  888   888   888  \n`88b    ooo   888   888  d8(  888    888 .  888  `88b.  888   888 888   888  888   888   888  \n `Y8bood8P'  o888o o888o `Y888\"\"8o   \"888\" o888o  o888o `Y8bod8P' `Y8bod8P' o888o o888o o888o \n                                                                                              ")

	//设置口令
	password := utils.SetPassword()
	var record string
	fmt.Println(password)
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
	record += "Server:" + string(buf[:n])
	// 启动一个 Goroutine 来持续接收服务器的消息
	messageChan := make(chan string)
	go utils.ReadLoop(conn, 0, messageChan)
	go func() {
		for message := range messageChan {
			record += message // 将接收到的消息追加到 record 中
		}
	}()
	// 主线程持续获取用户输入并发送消息
	scanner := bufio.NewScanner(os.Stdin)
	for {
		//fmt.Print("Enter message (type 'quit' to exit): ")
		scanner.Scan()
		text := scanner.Text()

		if strings.ToLower(text) == "quit" {
			fmt.Println("Exiting...")
			err := utils.SaveEncryptedData(time.Now().Format("2006-01-02_15-04-05"), record, password)
			if err != nil {
				log.Fatalf("Failed to save encrypted data: %v\n", err)
			}
			fmt.Println("Chat record encrypted and saved.")
			break
		}
		record += "Client:" + text + "\n"
		_, err := conn.Write([]byte(text))
		if err != nil {
			log.Fatalf("client: write: %s", err)
		}
	}

	//尝试读取聊天记录
	fmt.Println("[请输入口令]:")
	var pass string
	fmt.Scanln(&pass)
	if pass == password {
		// 指定记录目录路径
		directory := "record"
		dir, err := os.Open(directory)
		if err != nil {
			log.Fatalf("Failed to open directory: %v\n", err)
		}
		defer dir.Close()

		// 获取目录中的文件信息列表
		files, err := dir.Readdirnames(0)
		if err != nil {
			log.Fatalf("Failed to read directory names: %v\n", err)
		}

		// 按序号输出文件名
		fmt.Println("Files in directory:", directory)
		for i, file := range files {
			fmt.Printf("%d: %s\n", i+1, file)
		}

		for {
			// 提示用户选择文件
			fmt.Print("[输入查看的聊天记录序号]:")
			var index string
			fmt.Scanln(&index)

			// 转换输入的文件序号为整数
			fileIndex, err := strconv.Atoi(index)
			if err != nil || fileIndex < 1 || fileIndex > len(files) {
				log.Fatalf("Invalid file selection\n")
			}

			// 获取用户选择的文件名
			selectedFile := files[fileIndex-1]
			filePath := directory + "/" + selectedFile

			// 加载并解密数据
			decryptedData, err := utils.LoadEncryptedData(filePath, password)
			if err != nil {
				log.Fatalf("Failed to decrypt file: %v\n", err)
			}

			// 输出解密内容
			fmt.Println("[聊天记录]:")
			fmt.Println(decryptedData)
		}
	} else {
		fmt.Println("口令错误")
	}

}
