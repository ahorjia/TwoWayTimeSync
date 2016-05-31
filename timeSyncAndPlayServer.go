package main

import "net"
import "fmt"
import "bufio"
import "time"

//import "math"

//import "math"
import "strings"
import "strconv"

func main() {
	var a Audio
	a.initialize()
	defer a.close()

	fmt.Println("Launching server...")
	//	numTries := 3

	ln, _ := net.Listen("tcp", ":8081")
	conn, _ := ln.Accept()

	//	for index := 0; index < numTries; index++ {
	for {
		clientTimeMessage, readErr := bufio.NewReader(conn).ReadString('\n')
		receivedMessageTime := time.Now()
		if readErr != nil {
			fmt.Printf("Error reading connection: %v\n", readErr)
			break
		}

		if strings.Contains(clientTimeMessage, "PLAYTIME") {
			break
		}

		clientTimeMessage = strings.TrimSpace(clientTimeMessage)

		message := fmt.Sprintf("%v,%v,%v\n", clientTimeMessage,
			strconv.FormatInt(receivedMessageTime.UnixNano(), 10),
			strconv.FormatInt(time.Now().UnixNano(), 10))

		fmt.Println(message)
		_, err := fmt.Fprintf(conn, message)
		if err != nil {
			fmt.Println("Error")
		}
	}

	fmt.Println("READY TO PLAY")
	//	a.play(440, 2, math.Pi)
}
