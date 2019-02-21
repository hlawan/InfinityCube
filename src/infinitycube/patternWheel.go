package main

import (
	"time"
)

type BinaryWheel struct {
	Effect
	startLed     int
	size         int
	stepDuration time.Duration
	stepStart    time.Time
}

func NewBinaryWheel(disp Display, cg ColorGenerator, size int, duration time.Duration) *BinaryWheel {
	ef := NewEffect(disp, 0.5, 0.0)

	e := &BinaryWheel{
		Effect:       ef,
		size:         size,
		stepDuration: duration,
		stepStart:    time.Now(),
	}

	e.Painter = cg
	return e
}

func (e *BinaryWheel) Update() {
	for i := (0 + e.OffsetPar); i < (e.OffsetPar + e.LengthPar); i++ {
		e.Leds[i].setV(0)
		if i >= e.startLed && i < e.startLed+e.size {
			e.Leds[i].setV(0.8)
		}
		if time.Since(e.stepStart) > e.stepDuration {
			e.startLed = (e.startLed + e.size) % e.LengthPar
			e.stepStart = time.Now()
		}

	}

	// every update function of an effect ends with this snippet
	e.Painter.Update()
	e.Leds = e.Painter.Colorize(e.Leds)
	e.myDisplay.AddEffect(e.Effect)
}
