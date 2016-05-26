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

	duration, err := strconv.ParseInt(os.Args[2], 10, 32)

	if err != nil {
		fmt.Printf("Error")
	}

	portaudio.Initialize()
	defer portaudio.Terminate()
	s := newStereoSine(arg, arg, sampleRate)
	defer s.Close()
	chk(s.Start())
	time.Sleep(time.Duration(duration) * time.Second)
	chk(s.Stop())
}

type stereoSine struct {
	*portaudio.Stream
	stepL, phaseL float64
	stepR, phaseR float64
}

func newStereoSine(freqL, freqR, sampleRate float64) *stereoSine {
	s := &stereoSine{nil, freqL / sampleRate, 0, freqR / sampleRate, 0}
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

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
