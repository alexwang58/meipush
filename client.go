package main

import (
	"bufio"
	"fmt"
	"net"
	//"os"
	"strconv"
	"time"
)

var quitSemaphore chan bool

const HEART_BEAT_STR = "%0_#\n"

func main() {
	for i := 0; i < 10000; i++ {
		var tcpAddr *net.TCPAddr
		tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:9999")

		conn, _ := net.DialTCP("tcp", nil, tcpAddr)
		//conn.SetDeadline(0 * time.Duration)
		defer conn.Close()
		fmt.Println("connected!" + strconv.Itoa(i))
		go onMessageRecieved(conn)
	}

	// <-quitSemaphore
}

func onMessageRecieved(conn *net.TCPConn) {
	reader := bufio.NewReader(conn)
	for {
		go sendHeartBeat(conn)
		msg, err := reader.ReadString('\n')
		//fmt.Println(msg)
		if err != nil {
			//fmt.Println(err.Error())
			// quitSemaphore <- true
			//os.Exit(1)
			break
		}

		fmt.Println(msg)
		//time.Sleep(time.Second)
		//b := []byte(msg)
		//conn.Write(b)
	}
}

func sendHeartBeat(conn *net.TCPConn) {
	for {
		b := []byte(HEART_BEAT_STR)
		conn.Write(b)
		time.Sleep(3 * time.Second)
	}
}
