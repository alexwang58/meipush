package main

import (
	"bufio"
	"fmt"
	"net"
	//"os"
	"math/rand"
	"strconv"
	"time"
)

var quitSemaphore chan bool

const HB_STR = "%0_#\n"

type ClientConn struct {
	Conn   *net.TCPConn
	Status bool // false连接已断开
}

func main() {
	for i := 1; i <= 10; i++ {
		var tcpAddr *net.TCPAddr
		tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:9999")

		conn, _ := net.DialTCP("tcp", nil, tcpAddr)
		ClientConn := &ClientConn{Conn: conn, Status: true}
		//conn.SetDeadline(0 * time.Duration)
		defer conn.Close()
		fmt.Println("connected!" + strconv.Itoa(i))
		go onMessageRecieved(ClientConn)
	}

	<-quitSemaphore
}

func onMessageRecieved(client *ClientConn) {
	reader := bufio.NewReader(client.Conn)
	closeInfo := make(chan byte)
	go sendHeartBeat(client, closeInfo)
	for {
		msg, err := reader.ReadString('\n')
		//fmt.Println(msg)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("链接已关闭")
			} else {
				fmt.Println(err.Error())
			}
			client.Conn.Close()
			client.Status = false
			// quitSemaphore <- true
			//os.Exit(1)
			break
		}

		fmt.Println(msg)
	}
}

func sendHeartBeat(client *ClientConn, closeInfo chan byte) {
	for {
		if !client.Status {
			break
		}
		b := []byte(HB_STR)
		client.Conn.Write(b)
		slp := rand.Intn(10)

		fmt.Printf("RandSleep: %d\n", slp)
		time.Sleep(time.Duration(slp) * time.Second)
	}
}

func mainThreadQuit() {
	quitSemaphore <- true
}
