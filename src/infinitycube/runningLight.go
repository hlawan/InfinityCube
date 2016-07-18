package main

import (
	"github.com/lucasb-eyer/go-colorful"
	"math"
	"time"
)

type RunningLight struct {
	Consumer
	colorful.Color
	Position     float64
	Delta        float64
	Length       int
	ColorOpacity float64
	BlackOpacity float64
	Leds         [LEDS]Led
}

func NewRunningLight(color colorful.Color, length int, delta, colorOpacity, blackOpacity float64) *RunningLight {
	r := &RunningLight{
		Color:        color,
		Length:       length,
		Delta:        delta,
		ColorOpacity: colorOpacity,
		BlackOpacity: blackOpacity,
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

func (r *RunningLight) Tick(d time.Duration, o interface{}) {
	advance := o.(bool)

	if advance {
		r.Position += r.Delta
		if r.Position > 1 {
			r.Position -= 1
		}

		pos := r.Position * float64(r.Length)
		for i, _ := range r.Leds {
			j := i % r.Length
			r.Leds[i].Color = BLACK.BlendRgb(r.Color, 1-min(1, dist(pos, float64(j))))
		}
	}

	r.Consumer.AddPattern(r.Leds, r.ColorOpacity, r.BlackOpacity)
}

//-----------------------------------------------------------------------------
type BinaryRunningLight struct {
	Consumer
	Offset       int
	Length       int
	Direction    int
	ColorOpacity float64
	BlackOpacity float64
	Leds         [LEDS * 2]Led
}

func NewBinaryRunningLight(length, direction int, colorOpacity, blackOpacity float64) *BinaryRunningLight {
	g := &BinaryRunningLight{
		Length:       length,
		Direction:    direction,
		ColorOpacity: colorOpacity,
		BlackOpacity: blackOpacity,
	}
	for i, _ := range g.Leds {
		if i%g.Length == 0 {
			g.Leds[i].Color = colorful.Color{1, 0.4, 0}
		} else {
			g.Leds[i].Color = colorful.Color{0, 0, 0}
		}
	}
	return g
}

func (g *BinaryRunningLight) Tick(d time.Duration, o interface{}) {
	advance := o.(bool)
	if advance {
		g.Offset += g.Direction
		if g.Offset < 0 {
			g.Offset += g.Length
		}
		if g.Offset > g.Length {
			g.Offset -= g.Length
		}
	}

	var cube [LEDS]Led
	for i, v := range g.Leds[g.Offset : g.Offset+LEDS] {
		cube[i] = v
	}
	//g.Consumer.Tick(d, g.Leds[g.Offset:g.Offset+LEDS])
	g.Consumer.AddPattern(cube, g.ColorOpacity, g.BlackOpacity)
}

//-----------------------------------------------------------------------------
type GausRunningLight struct {
	Consumer
	colorful.Color
	Position     float64
	Delta        float64
	Length       int
	Interval     float64
	ColorOpacity float64
	BlackOpacity float64
	Leds         [LEDS]Led
}

func NewGausRunningLight(color colorful.Color, length int, interval, colorOpacity, blackOpacity float64) *GausRunningLight {
	r := &GausRunningLight{
		Color:        color,
		Length:       length,
		Interval:     interval,
		ColorOpacity: colorOpacity,
		BlackOpacity: blackOpacity,
	}
	r.Delta = (float64(1) / float64(r.Interval*fpsTarget))
	return r
}

func (r *GausRunningLight) Tick(d time.Duration, o interface{}) {
	advance := o.(bool)

	if advance {
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
	}

	r.Consumer.AddPattern(r.Leds, r.ColorOpacity, r.BlackOpacity)
}
