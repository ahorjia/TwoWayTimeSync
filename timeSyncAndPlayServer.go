package main

import "net"
import "fmt"
import "bufio"
import "time"

//import "math"
import "strings"
import "strconv"

func main() {
	fmt.Println("Launching server...")
	numTries := 10

	ln, _ := net.Listen("tcp", ":8081")
	conn, _ := ln.Accept()

	for index := 0; index < numTries; index++ {
		//	for {

		clientTimeMessage, readErr := bufio.NewReader(conn).ReadString('\n')
		receivedMessageTime := time.Now()
		if readErr != nil {
			fmt.Printf("Error reading connection: %v\n", readErr)
			break
		}
		clientTimeMessage = strings.TrimSpace(clientTimeMessage)

		message := fmt.Sprintf("%v,%v,%v\n", clientTimeMessage,
			strconv.FormatInt(receivedMessageTime.UnixNano(), 10),
			strconv.FormatInt(time.Now().UnixNano(), 10))

		fmt.Println(message)
		fmt.Fprintf(conn, message)
	}
}
