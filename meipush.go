package main

import (
	"bufio"
	"fmt"
	mp_util "github.com/alexwang58/meipush/util"
	"log"
	"net"
	"os"
	"runtime"
	//"strconv"
	"strings"
	"time"
)

var (
	HEART_BEAT_TOLERATE int8 = 2
	quitSemaphore       chan bool
	HB_EXPIRED          int64 = 10
	HB_CHECK                  = make(map[string]bool)
	MAX_CONNECTION_POOL       = 3
	FULL_INFO_SET             = false
	READER_DELIM        byte  = '\n'
	limitListener       *mp_util.LimitListener
	//MAINTAIN_CONNECTION chan bool
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

func ENVIRENT_INIT() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	ENVIRENT_INIT()
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
	limitListener = mp_util.NewLimitListener(tcpListener, MAX_CONNECTION_POOL)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	defer tcpListener.Close()
	//fmt.Printf("ConnectInstance:%d\n", i)
	//connectionInstance(tcpListener)
	connectionInstance(limitListener)
	//}

	//<-quitSemaphore
}

//func connectionInstance(tcpListener *net.TCPListener) {
func connectionInstance(limitListener *mp_util.LimitListener) {
	go heartBeatToleratePolling()

	for {

		limitListenerConn, err := limitListener.Accept() //tcpListener.AcceptTCP()
		if err != nil {
			//continue
			log.Fatalln(err)
		}
		fmt.Println("A client connected : " + limitListenerConn.RemoteAddr().String())
		go tcpPipe(limitListenerConn.TCPConn)

	}

}

func tcpPipe(conn *net.TCPConn) {
	ipStr := conn.RemoteAddr().String()
	defer func() {
		fmt.Println("disconnected :" + ipStr)
		closeConnection(conn)
	}()

	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString(READER_DELIM)
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

	if commandHandler(conn, message) {
		return
	}

	response = "[" + conn.RemoteAddr().String() + "] --> " + message
	// fmt.Println(response)
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

	//fmt.Printf("HB@%s\n", connTag)
	return
}

func commandHandler(conn *net.TCPConn, message string) (isCommand bool) {
	prefix := "#MYB<#CMD>"
	isCommand = true
	if !strings.HasPrefix(message, prefix) {
		isCommand = false
		return
	}
	command := message[len(prefix):]
	fmt.Println("[COMMAND]>"+command, "CLOSE_CONN")
	if command == "CLOSE_CONN" {
		fmt.Println("CloseByClient")
		closeConnection(conn)
	}
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
			closeConnection(connStruct.Conn)
		}
	}

	time.Sleep(1 * time.Second)
	goto HB_CHECKRECHECK
}

func closeConnection(conn *net.TCPConn) {
	conn.Close()
	// 将关闭链接，从心跳监测中移除
	delete(HB_HASHMAP, conn.RemoteAddr().String())
	limitListener.Release()
	fmt.Printf("Closed %s\n", conn.RemoteAddr().String())
}
