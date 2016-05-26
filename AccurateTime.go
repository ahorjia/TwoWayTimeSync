package main

import "syscall"

import "time"
import "fmt"
import "math"

func doTheMath(numTries int, timeNanoSlice []int64) {
	// Calculate the diff
	numDiffs := numTries - 1
	diffSlice := make([]float64, numDiffs)

	// Calculate the total diffs incrementally
	var totalDiffs float64 = 0

	for index := 0; index < numDiffs; index++ {
		i64Diff := timeNanoSlice[index+1] - timeNanoSlice[index]
		floatDiff := float64(i64Diff)
		diffSlice[index] = floatDiff
		totalDiffs += floatDiff
	}

	var meanDiff float64 = totalDiffs / float64(numDiffs)
	var sumVariance float64 = 0
	var variance float64 = 0
	var stdDev float64 = 0

	deltaSlice := make([]float64, numDiffs)
	varianceSlice := make([]float64, numDiffs)

	for index := 0; index < (numTries - 1); index++ {
		fDelta := diffSlice[index] - meanDiff
		fVariance := fDelta * fDelta
		deltaSlice[index] = fDelta
		varianceSlice[index] = fVariance
		sumVariance += fVariance
	}
	variance = sumVariance / float64(numDiffs-1)
	stdDev = math.Sqrt(variance)
	fmt.Printf("Mean: %f\nVariance: %f\nstdDev: %f\n", meanDiff, variance, stdDev)
}

func main() {

	fmt.Println("time.Now Calls...")
	numTries := 10000000
	timeSlices := make([]int64, numTries)

	for i := 0; i < numTries; i++ {
		timeSlices[i] = time.Now().UnixNano()
	}

	doTheMath(numTries, timeSlices)

	fmt.Println("Syscall:")
	timeSlices2 := make([]int64, numTries)

	for i := 0; i < numTries; i++ {
		var tv syscall.Timeval
		syscall.Gettimeofday(&tv)
		timeSlices2[i] = tv.Nano()
	}

	doTheMath(numTries, timeSlices2)
}
