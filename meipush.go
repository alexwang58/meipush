package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	// "runtime"
	"log"
	"strconv"
	"time"
)

var (
	HEART_BEAT_HOLD     = 5
	quitSemaphore       chan bool
	HB_EXPIRED          int64 = 60
	HB_CHECK                  = make(map[string]bool)
	MAX_CONNECTION_POOL       = 256
	CONNECTION_POOL           = 0
	FULL_INFO_SET             = false
)

type Connections struct {
	Conn     *net.TCPConn
	LastBeat int64
}

var HB_HASHMAP = make(map[string]Connections)

func main() {
	//runtime.GOMAXPROCS(runtime.NumCPU())
	//for i := 0; i < 10; i++ {
	var tcpAddr *net.TCPAddr
	//var lport = "127.0.0.1:" + strconv.Itoa(9980+i)
	lport := "127.0.0.1:9999"
	fmt.Println(lport)
	tcpAddr, err := net.ResolveTCPAddr("tcp", lport)
	if err != nil {
		log.Fatalln(err)
	}
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	defer tcpListener.Close()
	//fmt.Printf("ConnectInstance:%d\n", i)
	connectionInstance(tcpListener)
	//}

	//<-quitSemaphore
}

func connectionInstance(tcpListener *net.TCPListener) {
	for {
		if CONNECTION_POOL >= MAX_CONNECTION_POOL {
			if !FULL_INFO_SET {
				fmt.Println("Connection Full: " + strconv.Itoa(CONNECTION_POOL))
				FULL_INFO_SET = true
			}
			continue
		}
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			//continue
			log.Fatalln(err)
		}
		CONNECTION_POOL++
		fmt.Println("A client connected : " + tcpConn.RemoteAddr().String())
		go tcpPipe(tcpConn)
	}
}

func tcpPipe(conn *net.TCPConn) {
	ipStr := conn.RemoteAddr().String()
	defer func() {
		fmt.Println("disconnected :" + ipStr)
		conn.Close()
	}()

	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		messageHandler(conn, message)
	}
}

func messageHandler(conn *net.TCPConn, message string) {
	var response string
	//fmt.Println("REV: " + message)
	now := time.Now().Unix()
	if message == "%0_#\n" {
		lastCns, ok := HB_HASHMAP[conn.RemoteAddr().String()]
		if ok && ((now - lastCns.LastBeat) > HB_EXPIRED) {
			//diff := time.Now().Unix() - HB_HASHMAP[conn.RemoteAddr().String()].LastBeat
			cons := HB_HASHMAP[conn.RemoteAddr().String()]
			cons.LastBeat = now
		} else {
			HB_HASHMAP[conn.RemoteAddr().String()] = Connections{Conn: conn, LastBeat: now}
		}

		//fmt.Printf("HB@%s\n", conn.RemoteAddr().String())
	} else {
		response = "[" + time.Now().String() + "] --> " + message
		response = message
		//fmt.Println(response)
		b := []byte(response)
		conn.Write(b)
	}
}

func heartBeatCheck(conn *net.TCPConn) {
	for client, check := range HB_CHECK {
		diff := time.Now().Unix() - HB_HASHMAP[client].LastBeat
		if check {

		}
		if diff > HB_EXPIRED {
			HB_HASHMAP[conn.RemoteAddr().String()].Conn.Close()
		}
	}
}
