package main

import (
	//"fmt"
	"math"
	"time"
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

func NewSine(disp Display, cg ColorGenerator) *Sine {
	ef := NewEffect(disp, 0.5, 0.0)

	s := &Sine{
		Effect:    ef,
		Frequency: 2,
		offset:    0,
		loopTime:  1,
		Direction: true,
	}

	s.Painter = cg
	s.SetDelta(s.loopTime)
	s.LengthPar = disp.NrOfLeds()

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

	// every update function of an effect ends with this snippet
	s.Painter.Update()
	s.Leds = s.Painter.Colorize(s.Leds)
	s.myDisplay.AddPattern(s.Leds, s.ColorOpacity, s.BlackOpacity)
}

func MagmaPlasma(eH *EffectHandler, playTime time.Duration) map[Effector]time.Duration {
	// sine 1
	cc1 := NewConstantColor(1, 0)
	sine1 := NewSine(eH.myDisplay, cc1)
	sine1.Frequency = 2 * NR_OF_SIDES
	sine1.SetLoopTime(5 * NR_OF_SIDES)

	// sine 2
	cc2 := NewConstantColor(1, 30)
	sine2 := NewSine(eH.myDisplay, cc2)
	sine2.Frequency = 3 * NR_OF_SIDES
	sine2.SetLoopTime(7 * NR_OF_SIDES)

	// multi running light
	cc3 := NewConstantColor(1, 60)
	mrl := NewMultiRunningLight(eH.myDisplay, cc3)

	magmaPlasma := map[Effector]time.Duration{}

	magmaPlasma[sine1] = playTime
	magmaPlasma[sine2] = playTime
	magmaPlasma[mrl] = playTime

	return magmaPlasma
}
