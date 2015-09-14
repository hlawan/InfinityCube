package main

import (
	"github.com/lucasb-eyer/go-colorful"
	"time"
  "fmt"
)

const H_MAX = 360
const H_MIN = 0

type HsvFader struct {
	Consumer
	Length          int
	ColorDifference float64
	TimeFullFade    time.Duration
	Leds            [LEDS]Led
  loop            int
}

func NewHsvFader() *HsvFader {
	r := &HsvFader{
		Length:          LEDS,
		ColorDifference: 0.5,
		TimeFullFade:    10 * time.Second,
    loop:            0,
	}

	for i, _ := range r.Leds {
		r.Leds[i].Color = colorful.Color{255,0,0}
  }
  if DEBUG_LVL >= 1 {fmt.Println("HsvFader initialized all leds")}
	return r
}

func (r *HsvFader) Tick(start time.Time, o interface{}) {
  var  h, s float64
  // duration := time.Since(start)
	// nrOfSteps := int((H_MAX - H_MIN) / r.ColorDifference)
	// frequency := r.TimeFullFade / time.Duration(nrOfSteps)

  //fmt.Print(frequency)
  //fmt.Println("   duration % frequency == ", duration % frequency)
  //fmt.Print("   Led1 ", "is ", r.Leds[1].CIELCH.H, "   colorDifference is ", r.ColorDifference, " NewValue is ")

	if  true {
		for i, _ := range r.Leds {
      h, s, _ = r.Leds[i].Color.Hsv()

      r.Leds[i].Color = colorful.Hsv(h + r.ColorDifference, s, 1)
      r.Leds[i].CheckColor()
		}
	}
  r.loop ++
  fmt.Println(r.Leds[1].Color, " loop: ", r.loop)
  r.Consumer.Tick(time.Since(start), r.Leds[:])
}
