package main

import (
	"fmt"
	"math/rand"
	"time"
)

type CellularAutomata struct {
	Effect
	currentGeneration []bool
	nextGeneration    []bool
	ledVs             []float64
	secsPerGenPar     float64
	delta             float64
	generationChange  time.Time
}

func NewCellularAutomata(newDisp Display, cg ColorGenerator) *CellularAutomata {
	ef := NewEffect(newDisp, 0.5, 0.0)

	cA := &CellularAutomata{
		Effect:            ef,
		currentGeneration: make([]bool, len(ef.Leds)),
		nextGeneration:    make([]bool, len(ef.Leds)),
		ledVs:             make([]float64, len(ef.Leds)),
		secsPerGenPar:     2}

	cA.SetDelta(cA.secsPerGenPar)
	cA.Painter = cg
	cA.generationChange = time.Now()
	cA.sprinkle()
	return cA
}

func (cA *CellularAutomata) SetDelta(sec float64) {
	cA.delta = (float64(1) / (sec * float64(fpsTarget)))
	fmt.Println(cA.delta)
}

func (cA *CellularAutomata) SetSecsPerGen(sec float64) {
	cA.secsPerGenPar = sec
	cA.SetDelta(sec)
}

func (cA *CellularAutomata) sprinkle() {
	rand.Seed(int64(time.Now().Nanosecond()))

	for i := 0; i < len(cA.Leds); i++ {
		if rand.Float64() < 0.3 {
			cA.currentGeneration[i] = true
			cA.ledVs[i] = 1.0
		} else {
			cA.currentGeneration[i] = false
			cA.ledVs[i] = 0.0
		}
	}

	cA.calcNextGeneration()
}

func (cA *CellularAutomata) Update() {

	if time.Since(cA.generationChange) > (time.Duration(cA.secsPerGenPar) * time.Second) { //(time.Duration(int(1000*cA.SecsPerGenPar)) * time.Millisecond) {
		copy(cA.currentGeneration, cA.nextGeneration)
		cA.calcNextGeneration()
	}

	for i, led := range cA.ledVs {
		v := led

		if cA.nextGeneration[i] == true {
			v += cA.delta

		} else {
			v -= cA.delta

		}

		if v > 1.0 {
			v = 1.0
		}
		if v < 0.0 {
			v = 0.0
		}

		cA.ledVs[i] = v
		cA.Leds[i].setV(v)
	}

	// every update function of an effect ends with this snippet
	cA.Painter.Update()
	cA.Leds = cA.Painter.Colorize(cA.Leds)

	cA.myDisplay.AddEffect(cA.Effect)

}

func (cA *CellularAutomata) calcNextGeneration() {
	var a, b, c bool

	// wrap leds to a circle
	for i := 0; i < len(cA.Leds); i++ {
		if i == 0 {
			// first led => check end
			a = cA.currentGeneration[LEDS-1]
			b = cA.currentGeneration[i]
			c = cA.currentGeneration[i+1]
		} else if i == LEDS-1 {
			// last led => check beginning
			a = cA.currentGeneration[i-1]
			b = cA.currentGeneration[i]
			c = cA.currentGeneration[0]
		} else {
			// default => check neigbors
			a = cA.currentGeneration[i-1]
			b = cA.currentGeneration[i]
			c = cA.currentGeneration[i+1]
		}

		cA.nextGeneration[i] = r150(a, b, c)
	}

	cA.generationChange = time.Now()
}

func r150(a, b, c bool) bool {
	if (a && b && c) || (a && !b && !c) || (!a && b && !c) || (!a && !b && c) {
		return true
	}
	return false
}

func CellularAutomataMonochrome(eH *EffectHandler, playTime time.Duration) map[Effector]time.Duration {
	cc1 := NewConstantColor(1, 0)
	ca := NewCellularAutomata(eH.myDisplay, cc1)

	automata := map[Effector]time.Duration{}
	automata[ca] = playTime

	return automata
}
