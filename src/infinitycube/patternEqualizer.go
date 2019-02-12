package main

import (
	"time"
)

// Effect: Linear Spectrum
type LinearSpectrum struct {
	Effect
	sound *ProcessedAudio
}

func NewLinearSpectrum(disp Display, cg ColorGenerator, s *ProcessedAudio) *LinearSpectrum {
	ef := NewEffect(disp, 0.5, 0.0)

	e := &LinearSpectrum{
		Effect: ef,
		sound:  s,
	}

	e.Painter = cg
	return e
}

func (e *LinearSpectrum) Update() {
	for i := (0 + e.OffsetPar); i < (e.OffsetPar + e.LengthPar); i++ {
		e.Leds[i].V = e.sound.spektralDensity[i%EDGE_LENGTH]
	}

	// every update function of an effect ends with this snippet
	e.Painter.Update()
	e.Leds = e.Painter.Colorize(e.Leds)
	e.myDisplay.AddEffect(e.Effect)
}

func LinearSpectrumMonochrome(eH *EffectHandler, playTime time.Duration) map[Effector]time.Duration {
	cc1 := NewConstantColor(1, 0)
	spec := NewLinearSpectrum(eH.myDisplay, cc1, eH.audio)

	effectMap := map[Effector]time.Duration{}
	effectMap[spec] = playTime

	return effectMap
}

// Effect: Linear Edge Spectrum
type LinearEdgeSpectrum struct {
	Effect
	sound *ProcessedAudio
}

func NewLinearEdgeSpectrum(disp Display, cg ColorGenerator, s *ProcessedAudio) *LinearEdgeSpectrum {
	ef := NewEffect(disp, 0.5, 0.0)

	e := &LinearEdgeSpectrum{
		Effect: ef,
		sound:  s,
	}
	e.Painter = cg
	return e
}

func (e *LinearEdgeSpectrum) Update() {
	for i := (0 + e.OffsetPar); i < (e.OffsetPar + e.LengthPar); i++ {
		if i%EDGE_LENGTH < (EDGE_LENGTH / 2) {
			e.Leds[i].V = e.sound.spektralDensity[(i % EDGE_LENGTH)]
		} else {
			e.Leds[i].V = e.sound.spektralDensity[EDGE_LENGTH-(i%EDGE_LENGTH)]
		}
	}

	//every update function of an effect ends with this snippet
	e.Painter.Update()
	e.Leds = e.Painter.Colorize(e.Leds)
	e.myDisplay.AddEffect(e.Effect)
}

func LinearEdgeSpectrumMonochrome(eH *EffectHandler, playTime time.Duration) map[Effector]time.Duration {
	cc1 := NewConstantColor(0, 0)
	spec := NewLinearEdgeSpectrum(eH.myDisplay, cc1, eH.audio)

	effectMap := map[Effector]time.Duration{}
	effectMap[spec] = playTime

	return effectMap
}

//Effect: Edge Volume
type EdgeVolume struct {
	Effect
	sound *ProcessedAudio
}

func NewEdgeVolume(disp Display, cg ColorGenerator, s *ProcessedAudio) *EdgeVolume {
	ef := NewEffect(disp, 0.5, 0.0)

	e := &EdgeVolume{
		Effect: ef,
		sound:  s,
	}

	e.Painter = cg

	return e
}

func (e *EdgeVolume) Update() {

	for i := range e.Leds {

		effectIndex := i % EDGE_LENGTH

		if float64(effectIndex) < e.sound.maxPeak*EDGE_LENGTH {
			e.Leds[i].setV(1.0)
		} else {
			e.Leds[i].setV(0.0)
		}

	}

	// every update function of an effect ends with this snippet
	e.Painter.Update()
	e.Leds = e.Painter.Colorize(e.Leds)
	e.myDisplay.AddEffect(e.Effect)
}

func EdgeVolumeRedGreen(eH *EffectHandler, playTime time.Duration) map[Effector]time.Duration {

	green := Color{R: 0, G: 255, B: 0}
	red := Color{R: 255, G: 0, B: 0}

	colors := make([]Color, 2)
	colors[0] = green
	colors[1] = red

	colGrad := NewColorGradient(colors, EDGE_LENGTH)
	vol := NewEdgeVolume(eH.myDisplay, colGrad, eH.audio)

	effectMap := map[Effector]time.Duration{}
	effectMap[vol] = playTime

	return effectMap
}
