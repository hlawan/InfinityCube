// patternSparkle.go
package main

import (
	"math/rand"
)

type StarDust struct {
	Effect
	sound    *ProcessedAudio
	ledVs    []float64
	Chance   float64
	Cooldown float64
}

func NewStarDust(disp Display, cg ColorGenerator, s *ProcessedAudio) *StarDust {
	ef := NewEffect(disp, 0.5, 0.0)

	e := &StarDust{
		Effect:   ef,
		sound:    s,
		ledVs:    make([]float64, len(ef.Leds)),
		Chance:   0.001,
		Cooldown: 1.0 / float64(fpsTarget),
	}

	e.Painter = cg
	return e
}

func (e *StarDust) Update() {
	for i := (0 + e.OffsetPar); i < (e.OffsetPar + e.LengthPar); i++ {
		if e.ledVs[i] > 0 {

			nV := e.ledVs[i] - e.Cooldown

			if nV < 1e-3 {
				nV = 0
			}

			e.ledVs[i] = nV
			e.Leds[i].V = e.ledVs[i]

		} else {
			val := rand.Float64()
			if val < e.Chance {
				e.ledVs[i] = rand.Float64()
				e.Leds[i].V = e.ledVs[i]
			}
		}
	}

	// every update function of an effect ends with this snippet
	e.Painter.Update()
	e.Leds = e.Painter.Colorize(e.Leds)
	e.myDisplay.AddEffect(e.Effect)
}
