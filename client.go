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

var (
	quitSemaphore     chan bool
	HeartBeatDuration int = 5
	TestConnections   int = 1
)

const HB_STR = "%0_#\n"

type ClientConn struct {
	Conn   *net.TCPConn
	Status bool // false连接已断开
}

var exampleMessage = []string{
	/*
		"类与结构体在C++中有三点区别。[1] ",
		"（1）class中默认的成员访问权限是private的，而struct中则是public的。",
		"（2）从class继承默认是private继承，而从struct继承默认是public继承。",
		"（3）C++的结构体声明不必有struct关键字，而C语言的结构体声明必须带有关键字（使用typedef别名定义除外）。[1] ",
		"在实际项目中，结构体是大量存在的。研发人员常使用结构体来封装成新的类型。由于C语言内部程序比较简单，研发人员通常使用结构体创造新的“属性”，其目的是简化运算。[1] ",
		"结构体在函数中的作用不是简便，其最主要的作用就是封装。封装的好处就是可以再次利用。让使用者不必关心这个是什么，只要根据定义使用就可以了。[1] ",
		"结构体的大小与内存对齐",
		"结构体的大小不是结构体元素单纯相加就行的，因为我们主流的计算机使用的都是32bit字长的CPU，对这类型的CPU取4个字节的数要比取一个字节要高效，也更方便。所以在结构体中每个成员的首地",
		"址都是4的整数倍的话，取数据元素时就会相对更高效，这就是内存对齐的由来。每个特定平台上的编译器都有自己的默认“对齐系数”(也叫对齐模数)。程序员可以通过预编译命令#pragma pack(n)",
		"，n=1,2,4,8,16来改变这一系数，其中的n就是你要指定的“对齐系数”。[1] ",
		"Alice OK. So we'll find out later if you're right or wrong later on. Now let's listen to Andreas Wilkey, a psychologist at Clarkson University in New York, talking about why we're bad at assessing risk.",
		"INSERT Andreas Wilkey, Psychologist, Clarkson University, Potsdam, New York",
		"People typically fear anything which is small 一些属性来组 probability but it's extremely catastrophic if it were to happen… Think about dying in a plane crash, ",
		"think about a nuclear meltdown from the nearby power plant. Recently we have another increase in these birds' virus outbreaks in South Korea. People ",
		"read about that. And they may pay a lot of attention to that in the news but they may forget to get their flu shot.",
	*/
	"A",
	"F",
	"D",
	"K",
}

var exampleMessageLen int = len(exampleMessage)

func main() {
	for i := 1; i <= TestConnections; i++ {
		var tcpAddr *net.TCPAddr
		tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:9999")

		conn, _ := net.DialTCP("tcp", nil, tcpAddr)
		ClientConn := &ClientConn{Conn: conn, Status: true}
		//conn.SetDeadline(0 * time.Duration)
		defer ClientConn.Conn.Close()
		fmt.Println("connected!" + strconv.Itoa(i))
		go onMessageRecieved(ClientConn)
	}

	<-quitSemaphore
}

func onMessageRecieved(client *ClientConn) {
	reader := bufio.NewReader(client.Conn)
	closeInfo := make(chan byte)
	go sendHeartBeat(client, closeInfo)
	go sendMessage(client)
	for {
		//msg, err := reader.ReadString('\n')
		_, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("链接已关闭")
			} else {
				fmt.Println(err.Error())
			}
			onConnectionClose(client)
			// quitSemaphore <- true
			//os.Exit(1)
			break
		}
		//msg
		//fmt.Println("-------------\n" + msg + "-------------\n\n")
	}
}

func sendMessage(client *ClientConn) {
	for {
		if !client.Status {
			fmt.Println("[SendMsg]Close:" + client.Conn.RemoteAddr().String())
			break
		}
		time.Sleep(time.Duration(3) * time.Second)
		rIdx := rand.Intn(exampleMessageLen)
		message := messageWraper(exampleMessage[rIdx])
		fmt.Println("[" + client.Conn.RemoteAddr().String() + "] --> " + message)
		b := []byte(message)
		client.Conn.Write(b)
	}
}

func messageWraper(message string) string {
	return message + "\n"
}

func sendHeartBeat(client *ClientConn, closeInfo chan byte) {
	for {
		if !client.Status {
			break
		}
		b := []byte(HB_STR)
		client.Conn.Write(b)
		//slp := rand.Intn(10)

		//fmt.Printf("RandSleep: %d\n", slp)
		time.Sleep(time.Duration(HeartBeatDuration) * time.Second)
	}
}

func onConnectionClose(client *ClientConn) {
	client.Conn.Close()
	client.Status = false
}

func mainThreadQuit() {
	quitSemaphore <- true
}
