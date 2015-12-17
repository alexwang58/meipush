package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const (
	//绑定IP地址
	ip = ""
	//绑定端口号
	port = 3333
)

func main() {
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(ip), port, ""})
	if err != nil {
		fmt.Println("监听端口失败:", err.Error())
		return
	}
	fmt.Println("已初始化连接，等待客户端连接...")
	Server(listen)
}

var count = 0

func Server(listen *net.TCPListener) {
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			fmt.Println("接受客户端连接异常:", err.Error())
			continue
		}
		fmt.Println("客户端连接来自:", conn.RemoteAddr().String())
		count++
		defer conn.Close()
		go func() {
			//data := make([]byte, 3)
			reader := bufio.NewReader(conn)
			for {
				_, err := reader.ReadString('\n')
				//i, err := conn.Read(data)
				//fmt.Println(i)
				//fmt.Println("客户端发来数据:", string(data[0:i]))
				//fmt.Println("客户端发来数据:", message)
				if err != nil {
					fmt.Println("读取客户端数据错误:", err.Error())
					fmt.Println(count)
					os.Exit(1)
					break
				}
				conn.Write([]byte{'f', 'i', 'n', 'i', 's', 'h', '\n'})
			}

		}()
	}
}
