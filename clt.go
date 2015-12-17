package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
)

var quitSemaphore chan bool

func main() {
	for i := 0; i < 100; i++ {
		var tcpAddr *net.TCPAddr
		var subport = "999" + strconv.Itoa(i%10)
		subport = "9999"
		fmt.Println("127.0.0.1:" + subport)
		tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:"+subport)

		conn, err := net.DialTCP("tcp", nil, tcpAddr)

		if err != nil {
			log.Fatalln(err)
		}
		defer conn.Close()
		fmt.Println("connected!")

		go onMessageRecived(conn)

		// 控制台聊天功能加入
		// for {
		// 	var msg string
		// 	fmt.Scanln(&msg)
		// 	if msg == "quit" {
		// 		break
		// 	}
		// 	b := []byte(msg + "\n")
		// 	conn.Write(b)
		// }
	}
	<-quitSemaphore

}

func onMessageRecived(conn *net.TCPConn) {
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		fmt.Println(msg)
		if err != nil {
			quitSemaphore <- true
			break
		}
	}
}
