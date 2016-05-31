package main

import "math"

func main() {
	var a Audio
	a.initialize()
	defer a.close()
	a.play(440, 1, math.Pi)
}
