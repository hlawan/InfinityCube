package main

import (
	"time"
  "math/rand"
)

type CellularAutomata struct {
	Consumer
	//Offset       int
	ColorOpacity float64
	BlackOpacity float64
  Rule         int
  SecsPerGen   float64
  lastUpdate   time.Time
	Leds         [LEDS]Led
}

func NewCellularAutomata(colorOpacity, blackOpacity float64, rule int, secsPerGen float64) *CellularAutomata {
	cA := &CellularAutomata{
		ColorOpacity: colorOpacity,
		BlackOpacity: blackOpacity,
    Rule:         rule,
    SecsPerGen:   secsPerGen}

    rand.Seed(42)
    cA.lastUpdate = time.Now()

    for i := 0; i < LEDS; i++ {
      if rand.Float64() < 0.3 {
        cA.Leds[i].Color = red
      }else{
        cA.Leds[i].Color = black
      }
    }
    return cA
}

func (cA *CellularAutomata) Update(){
  var a,b,c bool
  var nextGen [LEDS]Led
  if (time.Since(cA.lastUpdate) > (150 * time.Millisecond)){
    for i := 0; i < LEDS; i++ {
      if i == 0 {
        a = cA.Leds[LEDS-1].OnOrOff()
        b = cA.Leds[i].OnOrOff()
        c = cA.Leds[i+1].OnOrOff()
      }else if i == LEDS-1 {
        a = cA.Leds[i-1].OnOrOff()
        b = cA.Leds[i].OnOrOff()
        c = cA.Leds[0].OnOrOff()
      }else{
        a = cA.Leds[i-1].OnOrOff()
        b = cA.Leds[i].OnOrOff()
        c = cA.Leds[i+1].OnOrOff()
      }

      if r150(a,b,c) == true{
        nextGen[i].Color = red;
      } else {
        nextGen[i].Color = black;
      }
    }
    cA.Leds = nextGen
    cA.Consumer.AddPreCube(cA.Leds, cA.ColorOpacity, cA.BlackOpacity)
    cA.lastUpdate = time.Now()
  }
}


func r150(a, b, c bool) bool{
  if ((a && b && c) || (a && !b && !c) || (!a && b && !c) || (!a && !b && c)) {
    return true
  }
  return false
}
