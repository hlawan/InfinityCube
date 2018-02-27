package main

import (
	//"fmt"
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

type RunningLight struct {
	Effect
	colorful.Color
	Position    float64
	IntervalPar float64
	delta       float64
	Bounce      bool
	Direction   bool
	ModePar     int8
}

func NewRunningLight(disp Display) *RunningLight {
	ef := NewEffect(disp, 0.5, 0.0)

	r := &RunningLight{
		Effect:      ef,
		Color:       red,
		IntervalPar: 10,
		Bounce:      false,
		Direction:   true,
		ModePar:     2,
	}

	r.SetDelta(r.IntervalPar)
	return r
}

func (r *RunningLight) SetDelta(sec float64) {
	r.delta = (float64(1) / float64(sec*fpsTarget))
}

func (r *RunningLight) Update() {

	// update Position of the Light on the Display [scaled: (0.0 ... 1.0)]
	r.moveLightPosition()

	// update the LightPattern depending on the choosen runningLight mode
	switch r.ModePar {
	case 0:
		r.updateSimple()
	case 1:
		r.updateStride()
	case 2:
		r.updateGauß()
	default:
		r.updateSimple()
	}

	r.myDisplay.AddPattern(r.Leds, r.ColorOpacity, r.BlackOpacity)
}

func (r *RunningLight) moveLightPosition() {

	if r.Bounce {
		// Running light runs back an forth
		// way there
		if r.Position >= 1 && r.Direction {
			r.Direction = false
		}
		// way back
		if r.Position <= 0 && !r.Direction {
			r.Direction = true
		}

		if r.Direction {
			r.Position += r.delta
		} else {
			r.Position -= r.delta
		}

	} else {
		// one way runningLight and at the end start from begining
		r.Position += r.delta
		if r.Position > 1 {
			r.Position = 0
		}
	}
}

func (r *RunningLight) updateSimple() {

	// calc ledNr from scaled position (0.0 ... 1.0)
	pos := r.Position * float64(r.LengthPar-1) // "-1" -> starting to count at 0
	ledNr := int(math.Round(pos))

	for i := 0; i < r.LengthPar; i++ {
		// all LEDs Black, only LED at current postion colored
		if i != ledNr {
			r.Leds[i+r.OffsetPar].Color = BLACK
		} else {
			r.Leds[i+r.OffsetPar].Color = r.Color
		}
	}
}

func (r *RunningLight) updateStride() {

	// calc light position on real display
	pos := r.Position * float64(r.LengthPar-1) // "-1" -> starting to count at 0

	for i := 0; i < r.LengthPar; i++ {
		// calculate the color of every LED based on the distance to the current position of the Light
		r.Leds[i+r.OffsetPar].Color = BLACK.BlendRgb(r.Color, 1-math.Min(1, dist(pos, float64(i))))
	}
}

func (r *RunningLight) updateGauß() {

	// calc light position on real display
	pos := r.Position * float64(r.LengthPar-1) // "-1" -> starting to count at 0

	for i := 0; i < r.LengthPar; i++ {
		distance := dist(pos, float64(i))
		gaus := (1 / (math.Sqrt(math.Pi / 3))) * math.Exp(-(1)*math.Pow(distance, float64(2)))
		r.Leds[i+r.OffsetPar].Color = BLACK.BlendRgb(r.Color, gaus)
	}
}

var BLACK = colorful.Color{0, 0, 0}

func dist(a, b float64) float64 {
	return math.Abs(a - b)
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
		IntervalPar:       5,
		ledsPerDisplayPar: 2}

	// for every edge
	for i := 0; i < NR_OF_SIDES*EDGES_PER_SIDE; i++ {
		shift := i * EDGE_LENGTH

		// runningLights per Display
		for o := 0; o < r.ledsPerDisplayPar; o++ {
			rl := NewRunningLight(r.myDisplay)
			rl.OffsetPar = shift
			rl.LengthPar = EDGE_LENGTH
			gap := 1.0 / float64(r.ledsPerDisplayPar-1)
			rl.Position = float64(o) * gap
			rl.Color = colorful.Color{1, 0, 0}
			rl.Bounce = true
			rl.SetDelta(r.IntervalPar)
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
