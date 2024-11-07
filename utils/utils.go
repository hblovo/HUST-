package utils

import (
	"fmt"
	"log"
	"net"
)

func ReadLoop(conn net.Conn, role int, messageChan chan<- string) {
	defer conn.Close()

	// 标记名称：服务器或客户端
	label := "Server"
	if role == 1 {
		label = "Client"
	}

	for {
		buf := make([]byte, 256)
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("%s: read error: %s", label, err)
			return
		}
		message := fmt.Sprintf("%s:%s\n", label, string(buf[:n]))
		fmt.Print(message)

		// 仅在服务器模式下（role == 0）将消息发送到 messageChan，用于记录
		if role == 0 {
			messageChan <- message
		}
	}
}
func SetPassword() string {
	var pass1, pass2 string
	for {
		fmt.Print("[请设置口令]：")
		_, err := fmt.Scanln(&pass1)
		if err != nil {
			fmt.Println("Error:", err)
			return ""
		}

		// 获取第二次输入的口令
		fmt.Print("[请再次输入口令]：")
		_, err = fmt.Scanln(&pass2)
		if err != nil {
			fmt.Println("Error:", err)
			return ""
		}

		// 检查两次输入是否一致
		if pass1 == pass2 {
			fmt.Println("[口令设置成功]")
			break // 如果匹配，退出循环
		} else {
			fmt.Println("[两次输入不一致，请重新设置]")
			// 如果不匹配，不退出循环，继续下一次循环让用户重新输入
		}
	}
	return pass1
}
