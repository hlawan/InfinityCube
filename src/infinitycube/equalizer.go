package main

import (
	"github.com/lucasb-eyer/go-colorful"
	//"fmt"
)

// Effect: White Spectrum
type WhiteSpectrum struct {
	*Effect
	sound *ProcessedAudio
}

func NewWhiteSpectrum(disp Display, s *ProcessedAudio) *WhiteSpectrum {
	ef := NewEffect(disp, 0.5, 0.0)

	e := &WhiteSpectrum{
		Effect: ef,
		sound:  s,
	}
	return e
}

func (e *WhiteSpectrum) Update() {
	for i := (0 + e.Offset); i < (e.Offset + e.Length); i++ {
		e.Leds[i].Color = colorful.Color{
			e.sound.spektralDensity[(i%EDGE_LENGTH)+1],
			e.sound.spektralDensity[(i%EDGE_LENGTH)+1],
			e.sound.spektralDensity[(i%EDGE_LENGTH)+1]}
	}
	e.myDisplay.AddPattern(e.Leds, e.ColorOpacity, e.BlackOpacity)
}

// Effect: White Edge Spectrum
type WhiteEdgeSpectrum struct {
	*Effect
	sound *ProcessedAudio
}

func NewWhiteEdgeSpectrum(disp Display, s *ProcessedAudio) *WhiteEdgeSpectrum {
	ef := NewEffect(disp, 0.5, 0.0)

	e := &WhiteEdgeSpectrum{
		Effect: ef,
		sound:  s,
	}
	return e
}

func (e *WhiteEdgeSpectrum) Update() {
	for i := (0 + e.Offset); i < (e.Offset + e.Length); i++ {
		if i%EDGE_LENGTH < (EDGE_LENGTH / 2) {
			e.Leds[i].Color = colorful.Color{
				e.sound.spektralDensity[(i%EDGE_LENGTH)+1],
				e.sound.spektralDensity[(i%EDGE_LENGTH)+10],
				e.sound.spektralDensity[(i%EDGE_LENGTH)+10]}
		} else {
			e.Leds[i].Color = colorful.Color{
				e.sound.spektralDensity[EDGE_LENGTH-(i%EDGE_LENGTH)+1],
				e.sound.spektralDensity[EDGE_LENGTH-(i%EDGE_LENGTH)+10],
				e.sound.spektralDensity[EDGE_LENGTH-(i%EDGE_LENGTH)+10]}
		}
	}
	e.myDisplay.AddPattern(e.Leds, e.ColorOpacity, e.BlackOpacity)
}

// Effect: Edge Volume
type EdgeVolume struct {
	*Effect
	sound *ProcessedAudio
}

func NewEdgeVolume(disp Display, s *ProcessedAudio) *EdgeVolume {
	ef := NewEffect(disp, 0.5, 0.0)

	e := &EdgeVolume{
		Effect: ef,
		sound:  s,
	}
	return e
}

func (e *EdgeVolume) Update() {
	for i := 0; i < e.Length; i++ {

		r := float64(i%EDGE_LENGTH) * 1.0 / float64(EDGE_LENGTH)
		g := notNegative(1 - (float64(i%EDGE_LENGTH) * 1.0 / float64(EDGE_LENGTH)) - 0.2)
		b := 0.0

		effectIndex := (i % (2 * EDGE_LENGTH))

		if effectIndex < (EDGE_LENGTH) {
			p := (e.Length + e.Offset) - (i) - EDGE_LENGTH
			e.Leds[p].Color = colorful.Color{r, g, b}
			if float64(effectIndex) > e.sound.maxPeak*EDGE_LENGTH {
				e.Leds[p].Color = e.Leds[p].Color.BlendRgb(black, 1.0)
			}
		} else {
			k := e.Offset + i
			e.Leds[k].Color = colorful.Color{r, g, b}
			if float64(effectIndex-EDGE_LENGTH) > e.sound.maxPeak*EDGE_LENGTH {
				e.Leds[k].Color = e.Leds[k].Color.BlendRgb(black, 1.0)
			}
		}
	}

	e.myDisplay.AddPattern(e.Leds, e.ColorOpacity, e.BlackOpacity)
}

func notNegative(nr float64) float64 {
	if nr < 0 {
		return 0
	} else {
		return nr
	}
}
