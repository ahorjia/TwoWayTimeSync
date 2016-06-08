package main

import "net"
import "fmt"
import "bufio"
import "time"
import "portaudio"
import "os"
import "math"
import "strings"
import "strconv"
import "log"

import "precisesleep"
func main() {
	fmt.Println("Initializing Audio");
	outputFrequency := 100.0;
        var err error
        if(len(os.Args) >=2) {
           outputFrequency, err = strconv.ParseFloat(os.Args[1],64)
           if(err != nil){
               errorString := fmt.Sprintf("error parsing frequency argument: %v\n",err);
               log.Fatal(errorString)
           }
        }
    portaudio.Initialize()
	defer portaudio.Terminate()
    s := newStereoSine(outputFrequency ,outputFrequency, sampleRate,0)
	defer s.Close()
    fmt.Println("Launching server...")
	ln, _ := net.Listen("tcp", ":8081")
	conn, _ := ln.Accept()

	for {
		
		clientTimeMessage, readErr := bufio.NewReader(conn).ReadString('\n')
        receivedMessageTime := time.Now()
		if(readErr != nil){
            fmt.Printf("Error reading connection: %v\n",readErr);
            break;
        }
        clientTimeMessage = strings.TrimSpace(clientTimeMessage)

		//		serverCurrenTimeText := strings.ToUpper(serverCurrentTime.Format("02:Jan:2006:15:04:05.999999"))
		//		fmt.Printf("Local time: %v, received client time: %v", serverCurrenTimeText, message)
        currentTime := time.Now();
        playTime := currentTime.Add( time.Duration(500*time.Millisecond));
		message := fmt.Sprintf("%v,%v,%v,%v\n", clientTimeMessage,
			strconv.FormatInt(receivedMessageTime.UnixNano(),10),
			strconv.FormatInt(time.Now().UnixNano(),10),
            strconv.FormatInt(playTime.UnixNano(),10));

		fmt.Println(message)
		fmt.Fprintf(conn, message)
        go s.playAt(playTime.UnixNano());
	}
}


// audio processing code
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
    precisesleep.SleepUntil64(unixNano);
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
