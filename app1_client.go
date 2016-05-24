package main

import "net"
import "fmt"
import "time"
import "strings"
import "bufio"

func tick(d time.Duration, f func(time.Time, net.Conn), conn net.Conn) {
	for x := range time.Tick(d) {
		f(x, conn)
	}
}

func printtime(t time.Time, conn net.Conn) {
	t0 := time.Now()
	text := strings.ToUpper(t0.Format("02:Jan:2006:15:04:05.999999"))
	fmt.Fprintf(conn, text+"\n")

	message, _ := bufio.NewReader(conn).ReadString('\n')
	t1 := time.Now()
	text1 := strings.ToUpper(t1.Format("02:Jan:2006:15:04:05.999999"))
	fmt.Printf("Local time: %v; received server time: %v", text1, message)
}

func main() {
	// conn, _ := net.Dial("tcp", "[2001::fa1a]:8081")
	conn, _ := net.Dial("tcp", "127.0.0.1:8081")
	tick(time.Second, printtime, conn)
}
