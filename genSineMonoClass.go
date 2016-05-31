package main

import (
	//	"fmt"
	"math"
	//	"os"
	"portaudio"
	//	"strconv"
	"time"
)

const sampleRate = 44100

type Audio struct{}

func (h Audio) initialize() {
	portaudio.Initialize()
}

func (h Audio) close() {
	portaudio.Terminate()
}

func (h Audio) play(freq float64, duration int64, phase float64) {
	s := newStereoSine(freq, sampleRate, phase)
	defer s.Close()
	chk(s.Start())
	time.Sleep(time.Duration(duration) * time.Second)
	chk(s.Stop())
}

type stereoSine struct {
	*portaudio.Stream
	step, phase float64
}

func newStereoSine(freq, sampleRate float64, phase float64) *stereoSine {
	s := &stereoSine{nil, freq / sampleRate, phase}
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
