package main

import (
	"bufio"
	"fmt"
	//	"io"
	//	"io/ioutil"
	"math"
	"os"

	"math/cmplx"
	"sort"

	"github.com/mjibson/go-dsp/dsputils"
	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/wav"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func convertTo64(ar []float32) []float64 {
	newar := make([]float64, len(ar))
	var v float32
	var i int
	for i, v = range ar {
		newar[i] = float64(v)
	}
	return newar
}

func Max(slice []float64) int {
	retVal := 0
	max := slice[0]
	for index := 1; index < len(slice); index++ {
		if slice[index] > max {
			max = slice[index]
			retVal = index
		}
	}
	return retVal
}

func main() {
	numSamples := 880

	// Reference: https://play.golang.org/p/2kEKXq-kUV
	f, err := os.Open("/home/agah/Dropbox/Classes/09 - CSE237B/Lab3/Lab3/output2.wav")
	check(err)
	_ = bufio.NewReader(f)

	w, err2 := wav.New(f)
	check(err2)

	numSamples = w.Samples
	samples, err3 := w.ReadFloats(numSamples)
	check(err3)

	// Equation 3-10.
	x := func(n int) float64 {
		val := float64(samples[n])
		return val

		//		slice32 := make([]float32, 1000)
		//		slice64 := convertTo64(slice32)

		//		wave0 := math.Sin(2.0 * math.Pi * float64(n) / 440.0)
		//		wave1 := 0.5 * math.Sin(2*math.Pi*float64(n)/220.0+3.0*math.Pi/220.0)
		//		return wave0 + wave1
		//		return wave0
	}

	// Discretize our function by sampling at 8 points.
	a := make([]float64, numSamples)
	for i := 0; i < numSamples; i++ {
		a[i] = x(i)
	}

	X := fft.FFTReal(a)

	//	fmt.Println(X)
	// Print the magnitude and phase at each frequency.
	for i := 0; i < numSamples; i++ {
		r, θ := cmplx.Polar(X[i])
		θ *= 360.0 / (2 * math.Pi)
		if dsputils.Float64Equal(r, 0) {
			θ = 0 // (When the magnitude is close to 0, the angle is meaningless)
		}

		if r > 0.1 { // !dsputils.Float64Equal(r, 0) {
			fmt.Printf("X(%d) = %.1f ∠ %.1f°\n", i, r, θ)
		}
	}

	mags := make([]float64, numSamples)
	for i := 0; i < numSamples; i++ {
		mags[i] = cmplx.Abs(X[i])
	}

	sort.Float64s(mags)
	//	maxVal := Max(mags)
	for i := numSamples; i > numSamples-20; i-- {
		fmt.Printf("X(%d) = %.1f\n", i, mags[i-1])
	}
}
