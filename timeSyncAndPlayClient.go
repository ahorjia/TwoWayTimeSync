package main

import "net"
import "fmt"
import "time"
import "strings"
import "strconv"
import "bufio"
import "os"
import "log"

func sync(d time.Duration, f func(time.Time, net.Conn) int64, conn net.Conn) {
	counter := 0.0
	var sum float64
	sum = 0
	for x := range time.Tick(d) {
		retVal := f(x, conn)
		counter += 1
		sum += float64(retVal)
		currentMean := sum / counter
		fmt.Printf("current delta: %v, mean:%v\n", retVal, currentMean)
	}
}

func printtime(t time.Time, conn net.Conn) int64 {
	text := strconv.FormatInt(time.Now().UnixNano(), 10)
	fmt.Fprintf(conn, text+"\n")

	message, connErr := bufio.NewReader(conn).ReadString('\n')
	if connErr != nil {
		quitString := fmt.Sprintf("Error connecting to server: %v\n", connErr)
		log.Fatal(quitString)
	}
	message = strings.TrimSpace(message)

	times := strings.Split(message, ",")
	t1, _ := strconv.ParseInt(times[0], 10, 64)
	t2, _ := strconv.ParseInt(times[1], 10, 64)
	t3, _ := strconv.ParseInt(times[2], 10, 64)
	t4 := time.Now().UnixNano()

	fmt.Printf("%v,%v,%v,%v\n", t1, t2, t3, t4)
	A := t2 - t1
	B := t4 - t3
	delta := (A - B) / 2.0
	return delta
}

func main() {
	numArgs := len(os.Args)
	address := "127.0.0.1"
	if numArgs >= 2 {
		address = os.Args[1]
	}
	fmt.Printf("Connecting to address %s\n", address)
	conn, err := net.Dial("tcp", address+":8081")
	if err == nil {
		sync(time.Second*2, printtime, conn)
	} else {
		fmt.Printf("Error connecting to server %v: %v\n", address, err)
	}
}
