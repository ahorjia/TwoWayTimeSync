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
	arg1, err := strconv.ParseFloat(os.Args[1], 32)

	if err != nil {
		fmt.Printf("Error")
	}

	arg2, err := strconv.ParseInt(os.Args[2], 10, 32)

	if err != nil {
		fmt.Printf("Error")
	}

	portaudio.Initialize()
	defer portaudio.Terminate()
	s := newSquare(arg1, sampleRate)
	defer s.Close()
	chk(s.Start())
	time.Sleep(time.Duration(arg2) * time.Second)
	chk(s.Stop())
}

type monoSquare struct {
	*portaudio.Stream
	stepL, phaseL float64
}

func newSquare(freqStart, sampleRate float64) *monoSquare {
	s := &monoSquare{nil, freqStart / sampleRate, 0}
	var err error
	s.Stream, err = portaudio.OpenDefaultStream(0, 1, sampleRate, 0, s.processAudio)
	chk(err)
	return s
}

func (g *monoSquare) processAudio(out [][]float32) {
	for i := range out[0] {
		testVal := float32(math.Sin(2 * math.Pi * g.phaseL))
		out[0][i] = float32(1)
		if testVal < 0 {
			out[0][i] = float32(-1)
		}

		_, g.phaseL = math.Modf(g.phaseL + g.stepL)
	}
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
