package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	listenAddress, _ := net.ResolveUDPAddr("udp", ":10001")
	listenConnection, _ := net.ListenUDP("udp", listenAddress)
	defer listenConnection.Close()

	sendAddressListener, _ := net.ResolveUDPAddr("udp", "127.0.0.1:10002")
	sendAddressLocal, _ := net.ResolveUDPAddr("udp", "127.0.0.1:10003")
	sendConnection, _ := net.DialUDP("udp", sendAddressLocal, sendAddressListener)
	defer sendConnection.Close()

	buf := make([]byte, 1024)

	for {
		// Receive
		n, _, _ := listenConnection.ReadFromUDP(buf)
		fmt.Println("App2 Reads from App1:" + string(buf[0:n]))

		// Send
		msg := time.Now().Format(time.RFC3339Nano)
		buf := []byte(msg)
		sendConnection.Write(buf)
	}
}
