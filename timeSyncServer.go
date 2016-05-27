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
	//	numTries := 10

	ln, _ := net.Listen("tcp", ":8081")
	conn, _ := ln.Accept()

	//	receiveTimeSlice := make([]time.Time, numTries)
	//	messageTimeSlice := make([]time.Time, numTries)

	//	for index := 0; index < numTries; index++ {
	for {
		beforeMessageTime := time.Now()
		clientTimeMessage, _ := bufio.NewReader(conn).ReadString('\n')
		clientTimeMessage = strings.TrimSpace(clientTimeMessage)

		//		serverCurrenTimeText := strings.ToUpper(serverCurrentTime.Format("02:Jan:2006:15:04:05.999999"))
		//		fmt.Printf("Local time: %v, received client time: %v", serverCurrenTimeText, message)

		message := fmt.Sprintf("%v,%v,%v\n", clientTimeMessage,
			strconv.Itoa(beforeMessageTime.Nanosecond()),
			strconv.Itoa(time.Now().Nanosecond()))

		//		fmt.Println(message)
		fmt.Fprintf(conn, message)
	}
}
