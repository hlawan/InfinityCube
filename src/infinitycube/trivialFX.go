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

	sb.Painter.Update()
	sb.Leds = sb.Painter.Colorize(sb.Leds)
	sb.myDisplay.AddPattern(sb.Leds, sb.ColorOpacity, sb.BlackOpacity)
}

func StaticGradient(eH *EffectHandler, playTime time.Duration) map[Effector]time.Duration {

	green := Color{R: 0, G: 255, B: 0}
	red := Color{R: 255, G: 0, B: 0}

	colors := make([]Color, 2)
	colors[0] = green
	colors[1] = red

	colGrad := NewColorGradient(colors, EDGE_LENGTH)
	solid := NewSolidBrightness(eH.myDisplay, colGrad, 1.0)

	effectMap := map[Effector]time.Duration{}
	effectMap[solid] = playTime

	return effectMap
}
