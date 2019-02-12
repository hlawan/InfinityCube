package main

import (
	"time"
)

type SolidBrightness struct {
	Effect
	Brightness float64
}

func NewSolidBrightness(disp Display, cg ColorGenerator, brightness float64) *SolidBrightness {
	ef := NewEffect(disp, 0.5, 0.0)

	sb := &SolidBrightness{
		Effect:     ef,
		Brightness: brightness}

	sb.Painter = cg
	return sb
}

func (sb *SolidBrightness) Update() {
	sb.mux.Lock()
	defer sb.mux.Unlock()

	for i, _ := range sb.Leds {
		sb.Leds[i].setV(sb.Brightness)
	}

	// every update function of an effect ends with this snippet
	sb.Painter.Update()
	sb.Leds = sb.Painter.Colorize(sb.Leds)
	sb.myDisplay.AddEffect(sb.Effect)
}
