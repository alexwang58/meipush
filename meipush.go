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
	HEART_BEAT_TOLERATE int8 = 2
	quitSemaphore       chan bool
	HB_EXPIRED          int64 = 10
	HB_CHECK                  = make(map[string]bool)
	MAX_CONNECTION_POOL       = 256
	CONNECTION_POOL           = 0
	FULL_INFO_SET             = false
)

const (
	HB_STR = "%0_#\n"
	//HB_TOLERATE_DURATION = HEART_BEAT_TOLERATE * HB_EXPIRED
)

type Connections struct {
	Conn     *net.TCPConn
	LastBeat int64
	MissBeat int8
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
	go heartBeatToleratePolling()
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
		_, ok := HB_HASHMAP[ipStr]
		// 未初始化链接池
		if !ok {
			now := time.Now().Unix()
			HB_HASHMAP[ipStr] = Connections{Conn: conn, LastBeat: now, MissBeat: 0}
		}

		messageHandler(conn, message)
	}
}

func messageHandler(conn *net.TCPConn, message string) {
	var response string
	isHeartBeat := updateHeartBeat(conn, message)
	if isHeartBeat {
		return
	}

	//conn.RemoteAddr().String()

	//response = "[" + time.Now().String() + "] --> " + message
	response = "[" + conn.RemoteAddr().String() + "] --> " + message
	//response = message
	fmt.Println(response)
	b := []byte(response)
	conn.Write(b)
}

// 更新心跳信息
func updateHeartBeat(conn *net.TCPConn, message string) (isHeartBeat bool) {

	isHeartBeat = false
	if message != HB_STR {
		// 不是心跳
		return
	}
	//fmt.Print("[REV: " + message + "--" + HB_STR + "]")

	isHeartBeat = true
	connTag := conn.RemoteAddr().String()
	lastCns, ok := HB_HASHMAP[connTag]
	if !ok {
		// Connection不在Map中，当做心跳处理
		return true
	}
	now := time.Now().Unix()
	lastCns.LastBeat = now
	lastCns.MissBeat = 0
	HB_HASHMAP[connTag] = lastCns
	/*
		// 心跳超时
		fmt.Println(now, lastCns.LastBeat, (now - lastCns.LastBeat))
		if (now - lastCns.LastBeat) > HB_EXPIRED {
			fmt.Println("错过")
			lastCns.MissBeat += 1
		}
		lastCns.LastBeat = now
		// 说好的指针呢
		HB_HASHMAP[connTag] = lastCns
	*/
	fmt.Printf("HB@%s\n", connTag)
	return
}

func heartBeatToleratePolling() {
HB_CHECKRECHECK:
	now := time.Now().Unix()
	for connTag, connStruct := range HB_HASHMAP {
		// 连续错过的心跳超过最大阈值
		if (now - connStruct.LastBeat) > HB_EXPIRED {
			fmt.Println("错过")
			connStruct.MissBeat += 1
			HB_HASHMAP[connTag] = connStruct
		}
		if connStruct.MissBeat >= HEART_BEAT_TOLERATE {
			fmt.Printf("TOLERATE: %d\n", connStruct.MissBeat)
			onHeartBeatExpired(connStruct.Conn)
		}
	}

	time.Sleep(1 * time.Second)
	goto HB_CHECKRECHECK
}

func onHeartBeatExpired(conn *net.TCPConn) {
	conn.Close()
	// 将关闭链接，从心跳监测中移除
	delete(HB_HASHMAP, conn.RemoteAddr().String())
	fmt.Printf("Closed %s\n", conn.RemoteAddr().String())
}
