package main

// Imports
import (
        "net"
        "fmt"
        "os"
        "log"
        "time"
        "bufio"
        "strconv"
        "strings"
        )

// Struct and enum declarations
// State Enum for the Client

const (
     STATE_SLAVE_CALIBRATION = iota
     STATE_SLAVE_MEASUREMENT
     STATE_SLAVE_CANCELLATION
)

const (
     MESSAGE_SYNCHRONIZATION_INFO = iota
     MESSAGE_AUDIO_RECEIVED
     MESSAGE_ERROR_INDICATED
     MESSAGE_SYNCHRONIZATION_BEGIN

)

type slaveMessage struct{
    messageType int
    data1 float64
    data2 float64
}

func slaveStateMachine(audioChan chan slaveMessage, syncChan chan slaveMessage, syncControlChan chan slaveMessage) {
    myState := STATE_SLAVE_CALIBRATION;
    calibrationDone := false
    var currentMessage slaveMessage;
    audioAmplitudeMasterSlave := 1.0;
    audioAmplitudeSlaveSlave := 1.0;
    audioTimeMasterSlave := 0.0;
    audioTimeSlaveSlave := 0.0;
    var timeDelta int64 = 0

    for {
        switch{
            case myState == STATE_SLAVE_CALIBRATION:
                select{
                case currentMessage = <-audioChan:
                    fmt.Printf("State %v, audioChan Received: %v,%v\n",myState,currentMessage.data1,currentMessage.data2);
                case currentMessage = <-syncChan:
                    fmt.Printf("State %v, Unexpected syncChan Received: %v,%v\n",myState,currentMessage.data1,currentMessage.data2);
                    timeDelta = int64(currentMessage.data1);
                case <- time.After(500*time.Millisecond):
                    fmt.Printf("Slave State: %v, Time Delta %v, Ams %v, Tms %v, Ass %v, Tss %v\n", myState, timeDelta, audioAmplitudeMasterSlave, audioTimeMasterSlave, audioAmplitudeSlaveSlave, audioTimeSlaveSlave);
               
                }
                
                if(calibrationDone){
                   syncBeginMessage := slaveMessage{MESSAGE_SYNCHRONIZATION_BEGIN,0,0}
                   syncControlChan <-syncBeginMessage;
                }
            case myState == STATE_SLAVE_MEASUREMENT:
                select{
                case currentMessage = <-audioChan:
                    fmt.Printf("State %v, audioChan Received: %v,%v\n",myState,currentMessage.data1,currentMessage.data2);
                case currentMessage = <-syncChan:
                    fmt.Printf("State %v, syncChan Received: %v,%v\n",myState,currentMessage.data1,currentMessage.data2);
                    timeDelta = int64(currentMessage.data1);
                case <- time.After(500*time.Millisecond):
                    fmt.Printf("Slave State: %v, Time Delta %v, Ams %v, Tms %v, Ass %v, Tss %v\n", myState, timeDelta, audioAmplitudeMasterSlave, audioTimeMasterSlave, audioAmplitudeSlaveSlave, audioTimeSlaveSlave);
               
                }
            case myState == STATE_SLAVE_CANCELLATION:
                select{
                case currentMessage = <-audioChan:
                    fmt.Printf("State %v, audioChan Received: %v,%v\n",myState,currentMessage.data1,currentMessage.data2);
                case currentMessage = <-syncChan:
                    fmt.Printf("State %v, syncChan Received: %v,%v\n",myState,currentMessage.data1,currentMessage.data2);
                    timeDelta = int64(currentMessage.data1);
                case <- time.After(500*time.Millisecond):
                    fmt.Printf("Slave State: %v, Time Delta %v, Ams %v, Tms %v, Ass %v, Tss %v\n", myState, timeDelta, audioAmplitudeMasterSlave, audioTimeMasterSlave, audioAmplitudeSlaveSlave, audioTimeSlaveSlave);
               
                }
            default:
                errorMessage := fmt.Sprintf("Unknown Slave State %v\n:", myState);
                log.Fatal(errorMessage);
        }
   }
}

func synchronizerThread(inChan chan slaveMessage, outChan chan slaveMessage,d time.Duration, f func(time.Time, net.Conn) int64, conn net.Conn) {
	counter := 0.0
    var sum float64
 	sum = 0
    var inCommand slaveMessage;
    // wait for an inputCommand to begin
    fmt.Println("Synchronizer waiting for signal to begin...");
    for {
        inCommand = <-inChan
        if(inCommand.messageType == MESSAGE_SYNCHRONIZATION_BEGIN){
            break;
        }
    }
    
    ticker := time.Tick(d);
    var delta float64;
	for {
        select{
        case x := <-ticker:
            retVal := f(x, conn)
            counter += 1
            delta = float64(retVal);
            sum += delta
            currentMean := sum/counter;              
            fmt.Printf("current delta: %v, mean:%v\n",retVal,currentMean)
            outMessage := slaveMessage{MESSAGE_SYNCHRONIZATION_INFO,delta,0.0}
            
            outChan <- outMessage;
        case inCommand = <-inChan:
            fmt.Printf("message received: %v\n", inCommand.messageType);
        }
	}
}

func sendReceiveToServer(t time.Time, conn net.Conn) int64 {
	text := strconv.FormatInt(time.Now().UnixNano(),10)
	fmt.Fprintf(conn, text+"\n")

	message, connErr := bufio.NewReader(conn).ReadString('\n')
    if(connErr != nil){
        quitString := fmt.Sprintf("Error connecting to server: %v\n",connErr);
        log.Fatal(quitString);
    }
	message = strings.TrimSpace(message)

	times := strings.Split(message, ",")
	t1, _ := strconv.ParseInt(times[0],10,64)
	t2, _ := strconv.ParseInt(times[1],10,64)
	t3, _ := strconv.ParseInt(times[2],10,64)
	t4 := time.Now().UnixNano()

    A := t2 - t1;
    B := t4 - t3;
    delta := (A - B)/ 2.0
	return delta
}

// The main fucntion sets up a connection to a server, creates the channels, and runs the three threads (audio,sync, stateMachine)

func main() {
    
    // conn, _ := net.Dial("tcp", "[2001::fa1a]:8081")
    numArgs := len(os.Args);
    address :=  "127.0.0.1" 
    if(numArgs >= 2){
       address = os.Args[1]
    }
    fmt.Printf("Connecting to server address %s\n",address);
    conn, err := net.Dial("tcp", address + ":8081" )
    if(err != nil){
        errorMessage  := fmt.Sprintf("Error connecting to server %v: %v\n",address,err);
        log.Fatal(errorMessage)
    }else{
        audioChan       := make(chan slaveMessage,5);
        syncChan        := make(chan slaveMessage,5);
        syncControlChan := make(chan slaveMessage,5);

        go synchronizerThread(syncControlChan,syncChan,time.Duration(1000 * time.Millisecond),sendReceiveToServer,conn);

        //FIXME: start the synchronization channel now, this should be handled later by the slaveStateMachine after calibration is complete
        syncBeginMessage := slaveMessage{MESSAGE_SYNCHRONIZATION_BEGIN,0,0}
        syncControlChan <-syncBeginMessage;
        
        slaveStateMachine(audioChan,syncChan,syncControlChan);
      
    }
}
