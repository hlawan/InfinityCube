package main

import (
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
	"math"
)

type RunningLight struct {
	*Effect
	colorful.Color
	Position float64
	DeltaPar float64
}

func NewRunningLight(disp Display) *RunningLight {
	ef := NewEffect(disp, 0.5, 0.0)

	r := &RunningLight{
		Effect:   ef,
		Color:    red,
		DeltaPar: 0.0001}

	return r
}

var BLACK = colorful.Color{0, 0, 0}

func max(a, b float64) float64 {
	if a < b {
		return b
	}
	return a
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func dist(a, b float64) float64 {
	return math.Abs(a - b)
}

func (r *RunningLight) Update() {

	//if advance { // need a new system to define speed here
	r.Position += r.DeltaPar
	if r.Position > 1 {
		r.Position -= 1
	}

	pos := r.Position * float64(r.Length)
	for i, _ := range r.Leds {
		j := i % r.Length
		r.Leds[i].Color = BLACK.BlendRgb(r.Color, 1-min(1, dist(pos, float64(j))))
	}
	//}

	r.myDisplay.AddPattern(r.Leds, r.ColorOpacity, r.BlackOpacity)
}

//-----------------------------------------------------------------------------
type GausRunningLight struct {
	*Effect
	colorful.Color
	Position    float64
	Delta       float64
	IntervalPar float64
	fpsTarget   int
}

func NewGausRunningLight(disp Display, fps int) *GausRunningLight {
	ef := NewEffect(disp, 0.5, 0.0)

	r := &GausRunningLight{
		Effect:      ef,
		Color:       blue,
		fpsTarget:   fps,
		IntervalPar: 30}

	r.Delta = (float64(1) / float64(r.IntervalPar*fpsTarget))
	return r
}

func (r *GausRunningLight) Update() {
	r.Delta = (float64(1) / float64(r.IntervalPar*fpsTarget))
	//	if advance {
	r.Position += r.Delta
	if r.Position > 1 {
		r.Position -= 1
	}
	pos := r.Position * float64(r.Length)
	for i, _ := range r.Leds {
		j := i % r.Length
		distance := dist(pos, float64(j))
		gaus := (1 / (math.Sqrt(math.Pi / 3))) * math.Exp(-(1)*math.Pow(distance, float64(2)))
		r.Leds[i].Color = BLACK.BlendRgb(r.Color, gaus)
	}
	//	}
	fmt.Println(r.IntervalPar)
	r.myDisplay.AddPattern(r.Leds, r.ColorOpacity, r.BlackOpacity)
}
