package main

import (
	"github.com/lucasb-eyer/go-colorful"
  //"fmt"
)

type Equalizer struct {
	Consumer
	Offset       int
	Length       int
	ColorOpacity float64
	BlackOpacity float64
	sound        *processedAudio
	Leds         [LEDS]Led
}

func NewEqualizer(offset, length int, colorOpacity, blackOpacity float64, s *processedAudio) *Equalizer {
	e := &Equalizer{
		Offset:       offset,
		Length:       length,
		ColorOpacity: colorOpacity,
		BlackOpacity: blackOpacity,
		sound:        s,
	}
	return e
}

func (e *Equalizer) WhiteSpectrum() {
	for i := (0 + e.Offset); i < (e.Offset + e.Length); i++ {
		e.Leds[i].Color = colorful.Color{
			e.sound.spektralDensity[(i%EDGE_LENGTH)+1],
			e.sound.spektralDensity[(i%EDGE_LENGTH)+1],
			e.sound.spektralDensity[(i%EDGE_LENGTH)+1]}
	}
	e.Consumer.AddPreCube(e.Leds, e.ColorOpacity, e.BlackOpacity)
}

func (e *Equalizer) WhiteEdgeSpectrum() {
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
	e.Consumer.AddPreCube(e.Leds, e.ColorOpacity, e.BlackOpacity)
}

func (e *Equalizer) EdgeVolume() {
	for i := 0; i < e.Length; i++ {

    r := float64(i % EDGE_LENGTH) * 1.0/float64(EDGE_LENGTH)
    g := notNegative(1-(float64(i % EDGE_LENGTH) * 1.0/float64(EDGE_LENGTH))-0.2)
    b := 0.0

    effectIndex := (i % (2 * EDGE_LENGTH))

		if effectIndex < (EDGE_LENGTH) {
      p := (e.Length + e.Offset) - (i) - EDGE_LENGTH
      e.Leds[p].Color = colorful.Color{r, g, b}
      if float64(effectIndex) > e.sound.maxPeak * EDGE_LENGTH{
        e.Leds[p].Color = e.Leds[p].Color.BlendRgb(black, 1.0)
      }
		} else {
      k := e.Offset + i
      e.Leds[k].Color = colorful.Color{r, g, b}
      if float64(effectIndex - EDGE_LENGTH) > e.sound.maxPeak * EDGE_LENGTH{
        e.Leds[k].Color = e.Leds[k].Color.BlendRgb(black, 1.0)
      }
    }
	}

	e.Consumer.AddPreCube(e.Leds, e.ColorOpacity, e.BlackOpacity)
}


func notNegative(nr float64) float64{
  if nr < 0 {
    return 0
  }else{
    return nr
  }
}
