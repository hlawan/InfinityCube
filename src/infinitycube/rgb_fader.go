package main

import (
	"github.com/lucasb-eyer/go-colorful"
	"time"
  "fmt"
)

const H_MAX = 360
const H_MIN = 0

type RgbFader struct {
	Consumer
	Length          int
	ColorDifference float32
	TimeFullFade    time.Duration
	Leds            [LEDS]Led
}

func NewRgbFader() *RgbFader {
	r := &RgbFader{
		Length:          LEDS,
		ColorDifference: .3,
		TimeFullFade:    10 * time.Second,
	}
  color := colorful.FastHappyColor()
	for i, _ := range r.Leds {
		r.Leds[i].SetColor(color)
    r.Leds[i].CIELCH.H = float32(i)
	  fmt.Println(r.Leds[i].CIELCH.H, " ")
  }
  if DEBUG_LVL >= 1 {fmt.Println("RgbFader initialized all leds")}
	return r
}

func (r *RgbFader) Tick(start time.Time, o interface{}) {
  duration := time.Since(start)
	nrOfSteps := int((H_MAX - H_MIN) / r.ColorDifference)
	frequency := r.TimeFullFade / time.Duration(nrOfSteps)

  //fmt.Print(frequency)
  //fmt.Println("   duration % frequency == ", duration % frequency)

	if  duration % frequency < (1000 * time.Microsecond) {
		for i, _ := range r.Leds {
      //fmt.Print("   Led1 ", "is ", r.Leds[1].CIELCH.H, "   colorDifference is ", r.ColorDifference, " NewValue is ")
			//r.Leds[i].CIELCH.H += r.colorDifference
      //fmt.Println(r.Leds[1].CIELCH.H)
			if r.Leds[i].CIELCH.H >= H_MAX {
				r.Leds[i].CIELCH.H = H_MIN
			}
		}
	}

  r.Consumer.Tick(time.Since(start), r.Leds[:])
}
