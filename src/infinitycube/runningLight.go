package main

import (
	//"fmt"
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

type RunningLight struct {
	Effect
	colorful.Color
	Position  float64
	DeltaPar  float64
	Bounce    bool
	Direction bool
}

func NewRunningLight(disp Display) *RunningLight {
	ef := NewEffect(disp, 0.5, 0.0)

	r := &RunningLight{
		Effect:    ef,
		Color:     red,
		DeltaPar:  0.0001,
		Bounce:    true,
		Direction: true,
	}

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

	if r.Bounce {
		// way there
		if r.Position >= 1 && r.Direction {
			r.Direction = false
		}
		// way back
		if r.Position <= 0 && !r.Direction {
			r.Direction = true
		}

		if r.Direction {
			r.Position += r.DeltaPar
		} else {
			r.Position -= r.DeltaPar

		}
	} else {
		r.Position += r.DeltaPar
		if r.Position > 1 {
			r.Position -= 1
		}
	}

	pos := r.Position * float64(r.LengthPar)
	//fmt.Println(pos)

	for i := 0; i < r.LengthPar; i++ {
		j := i
		if dist(pos, float64(j)) < 1 {
			//fmt.Println(dist(pos, float64(j)))
		}

		r.Leds[i+r.OffsetPar].Color = BLACK.BlendRgb(r.Color, 1-min(1, dist(pos, float64(j))))
	}
	//fmt.Println(r.Leds)
	// for i, _ := range r.Leds {
	// 	j := i
	// 	r.Leds[i].Color = BLACK.BlendRgb(r.Color, 1-min(1, dist(pos, float64(j))))
	// }

	r.myDisplay.AddPattern(r.Leds, r.ColorOpacity, r.BlackOpacity)
}

//-----------------------------------------------------------------------------
type GausRunningLight struct {
	Effect
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
	pos := r.Position * float64(r.LengthPar)
	for i, _ := range r.Leds {
		j := i % r.LengthPar
		distance := dist(pos, float64(j))
		gaus := (1 / (math.Sqrt(math.Pi / 3))) * math.Exp(-(1)*math.Pow(distance, float64(2)))
		r.Leds[i].Color = BLACK.BlendRgb(r.Color, gaus)
	}
	//	}
	r.myDisplay.AddPattern(r.Leds, r.ColorOpacity, r.BlackOpacity)
}

type MultiRunningLight struct {
	Effect
	runningLights     []Effector
	IntervalPar       float64
	fpsTarget         int
	ledsPerDisplayPar int
}

func NewMultiRunningLight(disp Display, fps int) *MultiRunningLight {
	ef := NewEffect(disp, 0.5, 0.0)

	r := &MultiRunningLight{
		Effect:            ef,
		fpsTarget:         fps,
		IntervalPar:       30,
		ledsPerDisplayPar: 2}

	for i := 0; i < NR_OF_SIDES*EDGES_PER_SIDE; i++ {
		shift := i * EDGE_LENGTH
		//fmt.Println(i)

		for o := 0; o < r.ledsPerDisplayPar; o++ {
			//fmt.Println("newRunningLight")
			//fmt.Println(o)
			rl := NewRunningLight(r.myDisplay)
			rl.OffsetPar = shift
			rl.LengthPar = EDGE_LENGTH
			rl.Position = (float64(EDGE_LENGTH) / float64(r.ledsPerDisplayPar)) * float64(o)
			rl.Color = colorful.Color{5, 0, 0}
			rl.DeltaPar = 0.001
			rl.BlackOpacity = 0
			r.runningLights = append(r.runningLights, rl)
		}
	}

	return r
}

func (r *MultiRunningLight) Update() {
	for _, effect := range r.runningLights {
		effect.Update()
	}
}
