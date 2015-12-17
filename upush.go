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

var CONNECTOR_NUM = 5000
var MESSAGE_PER_CON = 10

var quitSemaphore chan bool
var ccctop chan bool

func main() {
	for i := 0; i < CONNECTOR_NUM; i++ {
		go runConnect()
	}
	<-quitSemaphore
	/*
		for {
			sendHeartBeat(conn)
			fmt.Println("Wait")
			input, err := inputReader.ReadString('\n')
			if err != nil {
				fmt.Printf(err.Error())
			}

			//time.Sleep(time.Second)
			b := []byte(input)
			conn.Write(b)

			reader := bufio.NewReader(conn)
			msg, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Print(msg)
			}
		}
	*/
}

var connectOCu = 0

func runConnect() {
	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:9999")

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		//fmt.Println("DialTCP")
		//log.Fatalln(err)
		log.Println(err.Error())
		return
	}
	defer conn.Close()
	connectOCu++
	fmt.Println(connectOCu)

	//var str string
	//fmt.Scanln(&str)
	//fmt.Printf("INPUT :%s\n", str)
	//var inputReader *bufio.Reader
	//inputReader = bufio.NewReader(os.Stdin)
	go sendHeartBeat(conn)
	//go sendMessage(conn)

	//go sendRandMessage(conn)
	<-ccctop
}

func sendHeartBeat(conn *net.TCPConn) {
	/*
		beat := "%0_#"
		b := []byte(beat)
		conn.Write(b)
	*/
	for {
		//fmt.Println("Beat")
		conn.Write([]byte("%0_#\n"))
		sleepTimer := time.NewTimer(time.Second * 3)
		<-sleepTimer.C
	}
}

func sendMessage(conn *net.TCPConn) {
	inputReader := bufio.NewReader(os.Stdin)
	for {
		input, err := inputReader.ReadString('\n')

		if err != nil {
			fmt.Println("sendMessage_ReadString")
			log.Fatalln(err)
		}
		fmt.Print("INPUT:" + input)
		//time.Sleep(time.Second)
		b := []byte(input)
		conn.Write(b)

		getResponse(conn)
	}
}

func getResponse(conn *net.TCPConn) (string, bool) {
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

func sendRandMessage(conn *net.TCPConn) {
	average := 0
	for i := 0; i < MESSAGE_PER_CON; i++ {
		sucCount := 0
		failCount := 0
		connFaildCount := 0
		startT := time.Now().Unix()
		var diff int64
		for {
			input := randInput() + "\n"
			b := []byte(string(input))
			conn.Write(b)

			resp, success := getResponse(conn)
			if !success {
				fmt.Println("Something wrong")
				connFaildCount++
				break
			}

			if resp != string(input) {
				fmt.Printf("%s<-->%s\n", string(input), resp)
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
		fmt.Printf("Time used %d Seconds\nSuccess: %d\nConnection Faild: %d\nResponse Error: %d\nTotal messages: %d\n", diff, sucCount, connFaildCount, failCount, (connFaildCount + connFaildCount + failCount))
		fmt.Println("------------------")
	}

	fmt.Printf("Average: %0.2f\n", (float32(average) / float32(MESSAGE_PER_CON)))
}

func randInput() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.Itoa(r.Intn(100))
}
