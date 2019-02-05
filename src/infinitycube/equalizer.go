package main

import (
	//"fmt"
	"time"
)

//Effect: Linear Spectrum
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
	e.myDisplay.AddPattern(e.Leds, e.ColorOpacity, e.BlackOpacity)
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

	// every update function of an effect ends with this snippet
	e.Painter.Update()
	e.Leds = e.Painter.Colorize(e.Leds)
	e.myDisplay.AddPattern(e.Leds, e.ColorOpacity, e.BlackOpacity)
}

func LinearEdgeSpectrumMonochrome(eH *EffectHandler, playTime time.Duration) map[Effector]time.Duration {
	cc1 := NewConstantColor(0, 0)
	spec := NewLinearEdgeSpectrum(eH.myDisplay, cc1, eH.audio)

	effectMap := map[Effector]time.Duration{}
	effectMap[spec] = playTime

	return effectMap
}

//Effect: Edge Volume
// type EdgeVolume struct {
// 	Effect
// 	sound *ProcessedAudio
// }

// func NewEdgeVolume(disp Display, cg ColorGenerators *ProcessedAudio) *EdgeVolume {
// 	ef := NewEffect(disp, 0.5, 0.0)

// 	e := &EdgeVolume{
// 		Effect: ef,
// 		sound:  s,
// 	}
// 	return e
// }

// func (e *EdgeVolume) Update() {
// 	for i := 0; i < e.LengthPar; i++ {

// 		r := float64(i%EDGE_LENGTH) * 1.0 / float64(EDGE_LENGTH)
// 		g := notNegative(1 - (float64(i%EDGE_LENGTH) * 1.0 / float64(EDGE_LENGTH)) - 0.2)
// 		b := 0.0

// 		effectIndex := (i % (2 * EDGE_LENGTH))

// 		if effectIndex < (EDGE_LENGTH) {
// 			p := (e.LengthPar + e.OffsetPar) - (i) - EDGE_LENGTH
// 			e.Leds[p].Color = colorful.Color{r, g, b}
// 			if float64(effectIndex) > e.sound.maxPeak*EDGE_LENGTH {
// 				e.Leds[p].Color = e.Leds[p].Color.BlendRgb(black, 1.0)
// 			}
// 		} else {
// 			k := e.OffsetPar + i
// 			e.Leds[k].Color = colorful.Color{r, g, b}
// 			if float64(effectIndex-EDGE_LENGTH) > e.sound.maxPeak*EDGE_LENGTH {
// 				e.Leds[k].Color = e.Leds[k].Color.BlendRgb(black, 1.0)
// 			}
// 		}
// 	}

// 	e.myDisplay.AddPattern(e.Leds, e.ColorOpacity, e.BlackOpacity)
// }

// func notNegative(nr float64) float64 {
// 	if nr < 0 {
// 		return 0
// 	} else {
// 		return nr
// 	}
// }
