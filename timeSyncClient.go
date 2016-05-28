package main

import "net"
import "fmt"
import "time"
import "strings"
import "strconv"
import "bufio"
import "os"

func tick(d time.Duration, f func(time.Time, net.Conn) int, conn net.Conn) {
	counter := 0
	sum := 0
	for x := range time.Tick(d) {
		retVal := f(x, conn)
		counter += 1
		sum += retVal
		fmt.Println(sum / counter)
		//		fmt.Println(counter)
		//		fmt.Println(sum)
	}
}

func printtime(t time.Time, conn net.Conn) int {
	//	t0 := time.Now()
	//	text := strings.ToUpper(t0.Format("02:Jan:2006:15:04:05.999999"))
	text := strconv.Itoa(time.Now().Nanosecond())
	fmt.Fprintf(conn, text+"\n")

	message, _ := bufio.NewReader(conn).ReadString('\n')
	message = strings.TrimSpace(message)

	//	t1 := time.Now()
	//	text1 := strings.ToUpper(t1.Format("02:Jan:2006:15:04:05.999999"))
	//	fmt.Printf("Local time: %v; received server time: %v", text1, message)
	times := strings.Split(message, ",")
	t1, _ := strconv.Atoi(times[0])
	t2, _ := strconv.Atoi(times[1])
	t3, _ := strconv.Atoi(times[2])
	t4 := time.Now().Nanosecond()

	delta := (t4 - t3 + t2 - t1) / 2.0

	return delta
	//	fmt.Println(delta)
	//	message = fmt.Sprintf("%v,%v,%v,%v\n",
	//		strconv.Itoa(t1),
	//		strconv.Itoa(t2),
	//		strconv.Itoa(t3),
	//		strconv.Itoa(time.Now().Nanosecond()))

	//	fmt.Println(message)
}

func main() {
        
	// conn, _ := net.Dial("tcp", "[2001::fa1a]:8081")
        numArgs := len(os.Args);
        address :=  "127.0.0.1" 
        if(numArgs >= 2){
           address = os.Args[1]
        }
        fmt.Printf("Connecting to address %s\n",address);
        conn, err := net.Dial("tcp", address + ":8081" )
        if(err == nil){
          tick(time.Second, printtime, conn)
        }else{
          fmt.Printf("Error connecting to server %v: %v\n",address,err);
        }
}

