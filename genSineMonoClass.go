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

	initialize()
	play(arg, duration)
	close()

}

func initialize() {
	portaudio.Initialize()
}

func close() {
	portaudio.Terminate()
}

func play(freq float64, duration int64) {
	s := newStereoSine(freq, sampleRate)
	defer s.Close()
	chk(s.Start())
	time.Sleep(time.Duration(duration) * time.Second)
	chk(s.Stop())
}

type stereoSine struct {
	*portaudio.Stream
	step, phase float64
}

func newStereoSine(freq, sampleRate float64) *stereoSine {
	s := &stereoSine{nil, freq / sampleRate, 0}
	var err error
	s.Stream, err = portaudio.OpenDefaultStream(0, 1, sampleRate, 0, s.processAudio)
	chk(err)
	return s
}

func (g *stereoSine) processAudio(out [][]float32) {
	for i := range out[0] {
		out[0][i] = float32(math.Sin(2 * math.Pi * g.phase))
		_, g.phase = math.Modf(g.phase + g.step)
	}
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
