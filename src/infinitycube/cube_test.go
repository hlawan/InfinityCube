package main

import (
	"math"
	"testing"
)

func almosteq(v1, v2 float64) bool {
	var eps float64 = 1e-6
	return math.Abs(v1-v2) < eps
}

func TestScaleFinalPatternDownscale(t *testing.T) {
	// setup
	c := &Cube{
		finalCube: make([]Led, 60)}
	for i := range c.finalCube {
		c.finalCube[i].Color.R = 1.0
		c.finalCube[i].Color.G = 1.0
		c.finalCube[i].Color.B = 1.0
	}
	// execute
	var maxPowerWatts float64 = 45.0
	var parallelStripes uint32 = 6
	var expectedValue float64 = 0.5
	c.scaleFinalPattern(maxPowerWatts, parallelStripes)
	// check results
	for i := range c.finalCube {
		if !almosteq(c.finalCube[i].Color.R, expectedValue) ||
			!almosteq(c.finalCube[i].Color.G, expectedValue) ||
			!almosteq(c.finalCube[i].Color.B, expectedValue) {
			t.Errorf("Expect [%v, %v, %v], but got [%v, %v, %v]",
				expectedValue, expectedValue, expectedValue,
				c.finalCube[i].Color.R, c.finalCube[i].Color.G, c.finalCube[i].Color.B)
		}
	}
}

func TestScaleFinalPatternNoScale(t *testing.T) {
	// setup
	c := &Cube{
		finalCube: make([]Led, 60)}
	for i := range c.finalCube {
		c.finalCube[i].Color.R = float64(i) / float64(len(c.finalCube))
		c.finalCube[i].Color.G = 1.0 - (float64(i) / float64(len(c.finalCube)))
		c.finalCube[i].Color.B = 0.5
	}
	// execute
	var maxPowerWatts float64 = 90.0
	var parallelStripes uint32 = 6
	c.scaleFinalPattern(maxPowerWatts, parallelStripes)
	// check results
	for i := range c.finalCube {
		if !almosteq(c.finalCube[i].Color.R, float64(i)/float64(len(c.finalCube))) ||
			!almosteq(c.finalCube[i].Color.G, 1.0-(float64(i)/float64(len(c.finalCube)))) ||
			!almosteq(c.finalCube[i].Color.B, 0.5) {
			t.Errorf("Expect [%v, %v, %v], but got [%v, %v, %v]",
				i/len(c.finalCube), 1.0-(i/len(c.finalCube)), 0.5,
				c.finalCube[i].Color.R, c.finalCube[i].Color.G, c.finalCube[i].Color.B)
		}
	}
}
