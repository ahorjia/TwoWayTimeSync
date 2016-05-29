package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	ServerAddr2, _ := net.ResolveUDPAddr("udp", ":10002")
	ServerConn2, _ := net.ListenUDP("udp", ServerAddr2)
	defer ServerConn2.Close()
	read_buf := make([]byte, 1024)

	ServerAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:10001")
	LocalAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	Conn, _ := net.DialUDP("udp", LocalAddr, ServerAddr)

	defer Conn.Close()
	i := 0
	for {
		msg := time.Now().Format(time.RFC3339Nano)
		i++
		buf := []byte(msg)
		Conn.Write(buf)

		n, _, _ := ServerConn2.ReadFromUDP(read_buf)
		fmt.Println("App1 Reads from App2:" + string(read_buf[0:n]))
		time.Sleep(time.Second * 1)
	}
}
