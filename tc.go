package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	addr            = "127.0.0.1:3333"
	MESSAGE_PER_CON = 1
)

var quitSemaphore chan bool

func main() {

	for j := 0; j < 2000; j++ {
		go func() {
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				fmt.Println("连接服务端失败:", err.Error())
				fmt.Println(j)
				quitSemaphore <- true
				os.Exit(1)
				return
			}
			fmt.Println("已连接服务器")
			defer conn.Close()
			Client(conn)
		}()
	}

	<-quitSemaphore
}

func Client(conn net.Conn) {
	//sms := make([]byte, 128)
	//inputReader := bufio.NewReader(os.Stdin)
	//fmt.Println(i)
	sendRandMessage(conn)
	/*
		fmt.Print("请输入要发送的消息:")
		//_, err := fmt.Scan(&sms)
		input, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Println("数据输入异常:", err.Error())
		}
		fmt.Println(input)
		b := []byte(input)
		//conn.Write(input)
		conn.Write(b)
		buf := make([]byte, 128)
		c, err := conn.Read(buf)
		if err != nil {
			fmt.Println("读取服务器数据异常:", err.Error())
		}
		fmt.Println(string(buf[0:c]))
	*/

}

func sendRandMessage(conn net.Conn) {
	average := 0
	for i := 0; i < MESSAGE_PER_CON; i++ {
		sucCount := 0
		failCount := 0
		connFaildCount := 0
		startT := time.Now().Unix()
		var diff int64
		for {
			input := randInput() + "\n"
			//fmt.Print(input)
			b := []byte(string(input))
			conn.Write(b)

			resp, success := getResponse(conn)
			if !success {
				fmt.Println("Something wrong")
				connFaildCount++
				break
			}

			if resp != string(input) {
				//fmt.Printf("%s<-->%s\n", string(input), resp)
				failCount++
			} else {
				sucCount++
			}
			diff = time.Now().Unix() - startT
			if diff >= 1 {
				break
			}
		}
		average += sucCount
		//fmt.Printf("Time used %d Seconds\nSuccess: %d\nConnection Faild: %d\nResponse Error: %d\nTotal messages: %d\n", diff, sucCount, connFaildCount, failCount, (connFaildCount + connFaildCount + failCount))
		//fmt.Println("------------------")
	}

	//fmt.Printf("Average: %0.2f\n", (float32(average) / float32(MESSAGE_PER_CON)))
}

func randInput() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.Itoa(r.Intn(100))
}

func getResponse(conn net.Conn) (string, bool) {
	reader := bufio.NewReader(conn)
	msg, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("sendMessage_getResponse")
		log.Fatalln(err)
		return "", false
	} else {
		return msg, true
	}
}
