package main

import "net"
import "fmt"
import "bufio"
import "time"
import "math"
import "strings"

func main() {

	fmt.Println("Launching server...")
	numTries := 10

	ln, _ := net.Listen("tcp", ":8081")

	conn, _ := ln.Accept()

	//	receiveTimeSlice := make([]time.Time, numTries)
	messageTimeSlice := make([]time.Time, numTries)

	//	for index := 0; index < numTries; index++ {
	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		currentTime := time.Now()
		//		receiveTimeSlice[index] = currentTime
		//		messageTime, err := time.Parse("02:Jan:2006:15:04:05.999999\n", message)
		//		if err == nil {
		//			messageTimeSlice[index] = messageTime
		//		} else {
		//			fmt.Println(err)
		//		}

		//		if err != nil {
		//			fmt.Println(err)
		//		}

		timeText := strings.ToUpper(currentTime.Format("02:Jan:2006:15:04:05.999999"))

		fmt.Printf("Local time: %v, received client time: %v", timeText, message)

		fmt.Fprintf(conn, timeText+"\n")
		//conn.Write([]byte(newmessage + "\n"))
	}

	fmt.Println("Calculating statistics...")
	// Calculate the time to the nanosecond level
	timeNanoSlice := make([]int64, numTries)
	for index := 0; index < numTries; index++ {
		timeNanoSlice[index] = messageTimeSlice[index].UnixNano()
	}

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
	// Calculate the deltas and the variances
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
