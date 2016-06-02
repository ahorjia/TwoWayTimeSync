package main

import (
	"fmt"
	"math"
	"os"
	"portaudio"
	"strconv"
	"time"
)

const sampleRate = 44100

func main() {
	arg, err := strconv.ParseFloat(os.Args[1], 32)

	if err != nil {
		fmt.Printf("Error")
	}

	duration, err := strconv.Atoi(os.Args[2])

	if err != nil {
		fmt.Printf("Error")
	}
    
    repeatCount, err := strconv.Atoi(os.Args[3])
    
    if err != nil {
		fmt.Printf("Error")
	}
    
    

	portaudio.Initialize()
	defer portaudio.Terminate()
	s := newStereoSine(arg, arg, sampleRate,0)
	defer s.Close()
    
    for i:= 0; i< repeatCount;i++{
	chk(s.Start())
	   playTime := time.Now().Add( time.Duration(250*time.Millisecond));
       go s.playAt(playTime.UnixNano());
       time.Sleep(time.Duration(duration)*time.Second);
    }
}

type stereoSine struct {
	*portaudio.Stream
	stepL, phaseL float64
	stepR, phaseR float64
    playDuration time.Duration 
}

func newStereoSine(freqL, freqR, sampleRate float64, phase float64) *stereoSine {
	s := &stereoSine{nil, freqL / sampleRate, phase, freqR / sampleRate, phase,time.Duration(250*time.Millisecond)}
    
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
    deadline := time.Unix(0,unixNano)
    delay := deadline.Sub( time.Now())
    <-time.After(delay)
    g.Start();
    time.Sleep(g.playDuration)
    g.Stop();
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
