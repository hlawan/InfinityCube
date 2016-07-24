package main

import (
	"math/rand"
	"time"
)

type CellularAutomata struct {
	*Effect
	Rule          int
	SecsPerGenPar float64
	lastUpdate    time.Time
}

/*
func NewCellularAutomata(newDisplay Display, colorOpacity, blackOpacity float64, rule int, SecsPerGenPar float64) *CellularAutomata {
	cA := &CellularAutomata{
		ColorOpacity: colorOpacity,
		BlackOpacity: blackOpacity,
		Rule:         rule,
		SecsPerGenPar:   SecsPerGenPar,
		myDisplay:		newDisplay,
		Leds: 				make([]Led, newDisplay.NrOfLeds())}

	cA.lastUpdate = time.Now()
	cA.sprinkle()
	return cA
}*/

func NewCellularAutomata(newDisp Display) *CellularAutomata {
	ef := NewEffect(newDisp, 0.5, 0.0)

	cA := &CellularAutomata{
		Effect:        ef,
		Rule:          152,
		SecsPerGenPar: 0.33}

	cA.lastUpdate = time.Now()
	cA.sprinkle()
	return cA
}

func (cA *CellularAutomata) sprinkle() {
	rand.Seed(int64(time.Now().Nanosecond()))
	for i := 0; i < LEDS; i++ {
		if rand.Float64() < 0.3 {
			cA.Leds[i].Color = red
		} else {
			cA.Leds[i].Color = black
		}
	}
}

func (cA *CellularAutomata) Update() {
	var a, b, c bool
	nextGen := make([]Led, len(cA.Leds))
	if time.Since(cA.lastUpdate) > (150 * time.Millisecond) {
		for i := 0; i < LEDS; i++ {
			if i == 0 {
				a = cA.Leds[LEDS-1].OnOrOff()
				b = cA.Leds[i].OnOrOff()
				c = cA.Leds[i+1].OnOrOff()
			} else if i == LEDS-1 {
				a = cA.Leds[i-1].OnOrOff()
				b = cA.Leds[i].OnOrOff()
				c = cA.Leds[0].OnOrOff()
			} else {
				a = cA.Leds[i-1].OnOrOff()
				b = cA.Leds[i].OnOrOff()
				c = cA.Leds[i+1].OnOrOff()
			}

			if r150(a, b, c) == true {
				nextGen[i].Color = red
			} else {
				nextGen[i].Color = black
			}
		}
		cA.Leds = nextGen
		cA.myDisplay.AddPattern(cA.Leds, cA.ColorOpacity, cA.BlackOpacity)
		cA.lastUpdate = time.Now()
	}
}

func r150(a, b, c bool) bool {
	if (a && b && c) || (a && !b && !c) || (!a && b && !c) || (!a && !b && c) {
		return true
	}
	return false
}
