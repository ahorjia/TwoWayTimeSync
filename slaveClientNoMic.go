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
        "portaudio"
        "precisesleep"
        "math"
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
    data2 int64
}

func slaveStateMachine(audioChan chan slaveMessage, syncChan chan slaveMessage, syncControlChan chan slaveMessage, outputFrequency float64) {
    myState := STATE_SLAVE_CALIBRATION;
    nextState := myState;
    calibrationDone := false
    var currentMessage slaveMessage;
    audioAmplitudeMasterSlave := 1.0;
    audioAmplitudeSlaveSlave := 1.0;
    audioTimeMasterSlave := 0.0;
    audioTimeSlaveSlave := 0.0;
    var timeDelta int64 = 0
    var nextPlayTimeServer int64 = 0
    var nextPlayTimeClient int64 = 0
    s := newStereoSine(outputFrequency,outputFrequency, sampleRate,0.5)
	defer s.Close()

    for {
        switch{
            case myState == STATE_SLAVE_CALIBRATION:
                nextState = myState;
                select{
                case currentMessage = <-audioChan:
                    fmt.Printf("State %v, audioChan Received: %v,%v\n",myState,currentMessage.data1,currentMessage.data2);
                case currentMessage = <-syncChan:
                    fmt.Printf("State %v, Unexpected syncChan Received: %v,%v\n",myState,currentMessage.data1,currentMessage.data2);
                    timeDelta = int64(currentMessage.data1);
                    nextPlayTimeServer = currentMessage.data2;
                    calibrationDone = true
                case <- time.After(500*time.Millisecond):
                    fmt.Printf("Slave State: %v, Time Delta %v, Ams %v, Tms %v, Ass %v, Tss %v\n", myState, timeDelta, audioAmplitudeMasterSlave, audioTimeMasterSlave, audioAmplitudeSlaveSlave, audioTimeSlaveSlave);
               
                }
                
                if(calibrationDone){
                   syncBeginMessage := slaveMessage{MESSAGE_SYNCHRONIZATION_BEGIN,0,0}
                   syncControlChan <-syncBeginMessage;
                   nextState =  STATE_SLAVE_CANCELLATION;
                }
            case myState == STATE_SLAVE_MEASUREMENT:
                nextState = myState;
                select{
                case currentMessage = <-audioChan:
                    fmt.Printf("State %v, audioChan Received: %v,%v\n",myState,currentMessage.data1,currentMessage.data2);
                    nextState = STATE_SLAVE_CANCELLATION;
                case currentMessage = <-syncChan:
                    fmt.Printf("State %v, syncChan Received: %v,%v\n",myState,currentMessage.data1,currentMessage.data2);
                    timeDelta = int64(currentMessage.data1);
                    nextPlayTimeServer = currentMessage.data2;
                case <- time.After(500*time.Millisecond):
                    fmt.Printf("Slave State: %v, Time Delta %v, Ams %v, Tms %v, Ass %v, Tss %v\n", myState, timeDelta, audioAmplitudeMasterSlave, audioTimeMasterSlave, audioAmplitudeSlaveSlave, audioTimeSlaveSlave);
               
                }
            case myState == STATE_SLAVE_CANCELLATION:
                nextState = myState;
                select{
                case currentMessage = <-audioChan:
                    fmt.Printf("State %v, audioChan Received: %v,%v\n",myState,currentMessage.data1,currentMessage.data2);
                    nextState =  STATE_SLAVE_CANCELLATION;
                case currentMessage = <-syncChan:
                    fmt.Printf("State %v, syncChan Received: %v,%v\n",myState,currentMessage.data1,currentMessage.data2);
                    timeDelta = int64(currentMessage.data1);
                    nextPlayTimeServer = currentMessage.data2;
                    nextPlayTimeClient = nextPlayTimeServer - timeDelta;
                    s.playAt(nextPlayTimeClient);
                case <- time.After(500*time.Millisecond):
                    fmt.Printf("Slave State: %v, Time Delta %v, Ams %v, Tms %v, Ass %v, Tss %v\n", myState, timeDelta, audioAmplitudeMasterSlave, audioTimeMasterSlave, audioAmplitudeSlaveSlave, audioTimeSlaveSlave);
               
                }
            default:
                errorMessage := fmt.Sprintf("Unknown Slave State %v\n:", myState);
                log.Fatal(errorMessage);
        }
        myState = nextState;
   }
}

func synchronizerThread(inChan chan slaveMessage, outChan chan slaveMessage,d time.Duration, f func(time.Time, net.Conn) (int64,int64), conn net.Conn) {
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
            retVal,playTimeServer := f(x, conn)
            counter += 1
            delta = float64(retVal);
            sum += delta
            currentMean := sum/counter;              
            fmt.Printf("current delta: %v, mean:%v\n",retVal,currentMean)
            outMessage := slaveMessage{MESSAGE_SYNCHRONIZATION_INFO,delta,playTimeServer}
            
            outChan <- outMessage;
        case inCommand = <-inChan:
            fmt.Printf("message received: %v\n", inCommand.messageType);
        }
	}
}

func sendReceiveToServer(t time.Time, conn net.Conn) (int64,int64) {
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
    playTimeServer, _ := strconv.ParseInt(times[3],10,64)

    A := t2 - t1;
    B := t4 - t3;
    delta := (A - B)/ 2.0
	return delta,playTimeServer
}

// The main fucntion sets up a connection to a server, creates the channels, and runs the three threads (audio,sync, stateMachine)

func main() {
    portaudio.Initialize()
	defer portaudio.Terminate()
    outputFrequency := 100.0
    // conn, _ := net.Dial("tcp", "[2001::fa1a]:8081")
    numArgs := len(os.Args);
    address :=  "127.0.0.1" 
    if(numArgs >= 2){
       address = os.Args[1]
    }
    if(numArgs >= 3){
       frequency,err:= strconv.ParseFloat(os.Args[2],64);
       if(err != nil){
        errorMessage  := fmt.Sprintf("Error parsing frequency %v: %v\n",os.Args[2],err);
        log.Fatal(errorMessage)
       }else{
         outputFrequency = frequency;
       }
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
        
        slaveStateMachine(audioChan,syncChan,syncControlChan,outputFrequency);
      
    }
}

// audio processing code
const sampleRate = 48000

type stereoSine struct {
	*portaudio.Stream
	stepL, phaseL float64
	stepR, phaseR float64
    playDuration time.Duration 
    initialPhase float64
}

func newStereoSine(freqL, freqR, sampleRate float64, phase float64) *stereoSine {
	s := &stereoSine{nil, freqL / sampleRate, phase, freqR / sampleRate, phase,time.Duration(250*time.Millisecond),phase}
    
	var err error
	s.Stream, err = portaudio.OpenDefaultStream(0, 2, sampleRate, 0, s.processAudio)
	chk(err)
	return s
}


func (g *stereoSine) processAudio(out [][]float32) {
	for i := range out[0] {
		out[0][i] = float32(math.Sin(2 * math.Pi * g.phaseL))
		_, g.phaseL = math.Modf(g.phaseL + g.stepL)
		out[1][i] = float32(math.Sin(2 * math.Pi * g.phaseR))
		_, g.phaseR = math.Modf(g.phaseR + g.stepR)
	}
}

func (g * stereoSine) playAt(unixNano int64) {
    precisesleep.SleepUntil64(unixNano)
    g.phaseL = g.initialPhase;
    g.phaseR = g.initialPhase;
    g.Start();
    time.Sleep(g.playDuration)
    g.Stop();
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
