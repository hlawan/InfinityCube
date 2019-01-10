package main

import (
	//"fmt"
	"math"
)

type Sine struct {
	Effect
	Position  float64
	Frequency float64 // defines how many waves are displayed per edge
	offset    float64 // moves the waves along the edge
	delta     float64 // step size the offset increases per loop
	loopTime  float64 // [s] time the waves need to move one edge length
	Direction bool
}

func NewSine(disp Display) *Sine {
	ef := NewEffect(disp, 0.5, 0.0)

	s := &Sine{
		Effect:    ef,
		Frequency: 2,
		offset:    0,
		loopTime:  1,
		Direction: true,
	}

	s.Painter = NewConstantColor(1, 0)
	s.SetDelta(s.loopTime)
	s.LengthPar = EDGE_LENGTH

	return s
}

func (s *Sine) SetDelta(sec float64) {
	s.delta = (float64(1) / (sec * float64(fpsTarget))) * 2 * math.Pi * s.Frequency
}

func (s *Sine) SetLoopTime(sec float64) {
	s.loopTime = sec
	s.SetDelta(sec)
}

func (s *Sine) Update() {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.offset = s.offset + s.delta
	s.offset = math.Mod(s.offset, 2*math.Pi)

	for i := 0; i < s.LengthPar; i++ {
		distance := float64(i) * ((2.0 * math.Pi) / float64(s.LengthPar))
		sine := 0.8 + math.Sin((s.Frequency*distance)+s.offset)
		s.Leds[i+s.OffsetPar].V = sine

	}

	s.Leds = s.Painter.Colorize(s.Leds)

	s.myDisplay.AddPattern(s.Leds, s.ColorOpacity, s.BlackOpacity)
}
