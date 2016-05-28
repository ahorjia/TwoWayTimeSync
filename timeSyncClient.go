package main

import "net"
import "fmt"
import "time"
import "strings"
import "strconv"
import "bufio"
import "os"
import "log"

func tick(d time.Duration, f func(time.Time, net.Conn) int64, conn net.Conn) {
	counter := 0.0
        var sum float64
 	sum = 0
	for x := range time.Tick(d) {
		retVal := f(x, conn)
		counter += 1
		sum += float64(retVal)
                currentMean := sum/counter;              
		fmt.Printf("current delta: %v, mean:%v\n",retVal,currentMean)
		//		fmt.Println(counter)
		//		fmt.Println(sum)
	}
}

func printtime(t time.Time, conn net.Conn) int64 {
	//	t0 := time.Now()
	//	text := strings.ToUpper(t0.Format("02:Jan:2006:15:04:05.999999"))
	text := strconv.FormatInt(time.Now().UnixNano(),10)
	fmt.Fprintf(conn, text+"\n")

	message, connErr := bufio.NewReader(conn).ReadString('\n')
        if(connErr != nil){
            quitString := fmt.Sprintf("Error connecting to server: %v\n",connErr);
            log.Fatal(quitString);
        }
	message = strings.TrimSpace(message)

	//	t1 := time.Now()
	//	text1 := strings.ToUpper(t1.Format("02:Jan:2006:15:04:05.999999"))
	//	fmt.Printf("Local time: %v; received server time: %v", text1, message)
	times := strings.Split(message, ",")
	t1, _ := strconv.ParseInt(times[0],10,64)
	t2, _ := strconv.ParseInt(times[1],10,64)
	t3, _ := strconv.ParseInt(times[2],10,64)
	t4 := time.Now().UnixNano()

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

