package main

import (
	//	"bufio"
	"fmt"
	//	"io"
	//	"io/ioutil"
	//	"os"
	"math"

	"math/cmplx"

	"github.com/mjibson/go-dsp/dsputils"
	"github.com/mjibson/go-dsp/fft"
	//	"github.com/mjibson/go-dsp/wav"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	numSamples := 880

	// Equation 3-10.
	x := func(n int) float64 {
		wave0 := math.Sin(2.0 * math.Pi * float64(n) / 440.0)
		wave1 := 0.5 * math.Sin(2*math.Pi*float64(n)/220.0+3.0*math.Pi/220.0)
		return wave0 + wave1
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

		if !dsputils.Float64Equal(r, 0) {
			fmt.Printf("X(%d) = %.1f ∠ %.1f°\n", i, r, θ)
		}
	}
}
