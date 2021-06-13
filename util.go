package main

import "math"

func max(x, y float64) float64 {
	if x > y {
		return x
	}
	return y
}

func min(x, y float64) float64 {
	if x < y {
		return x
	}
	return y
}

func rnd(x float64) int {
	return int(math.Round(x))
}
