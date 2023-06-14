package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	go udpSendProc()
	go udpRecvProc()
}

func udpSendProc() {
	udpAddr, err := net.ResolveUDPAddr("udp", ":8888")
	if err != nil {
		fmt.Println("Err resolve UDP address: ", err)
		return
	}

	serverConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("ListenUDP error: ", err)
		return
	}

	var ticker = time.Tick(time.Second * 2) // 每隔2秒钟发送一个数据

	for {
		for _ = range ticker {
			var buff [512]byte
			n, rAddr, err := serverConn.ReadFromUDP(buff[0:])
			if err != nil {
				fmt.Println("Read error: ", err)
				break
			}
			fmt.Println("Read from client: ", string(buff[:n]))
			// 如果使用Write，本地测试时客户端接收不到信息
			serverConn.WriteToUDP([]byte("Hello client"), rAddr)
		}
	}
}

func udpRecvProc() {
	udpAddr, err := net.ResolveUDPAddr("udp", ":8888")
	if err != nil {
		fmt.Println("Err resolve UDP address: ", err)
		return
	}

	serverConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("ListenUDP error: ", err)
		return
	}

	var ticker = time.Tick(time.Second * 2) // 每隔2秒钟发送一个数据

	for {
		for _ = range ticker {
			var buff [512]byte
			n, rAddr, err := serverConn.ReadFromUDP(buff[0:])
			if err != nil {
				fmt.Println("Read error: ", err)
				break
			}
			fmt.Println("Read from client: ", string(buff[:n]))
			// 如果使用Write，本地测试时客户端接收不到信息
			serverConn.WriteToUDP([]byte("Hello client"), rAddr)
		}
	}
}
