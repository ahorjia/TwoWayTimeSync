package main

import (
	"net"
	"time"
)

func main() {
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
		time.Sleep(time.Second * 1)
	}
}
