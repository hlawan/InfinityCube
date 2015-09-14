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
	ColorDifference float64
	TimeFullFade    time.Duration
	Leds            [LEDS]Led
  loop            int
}

func NewRgbFader() *RgbFader {
	r := &RgbFader{
		Length:          LEDS,
		ColorDifference: 0.1,
		TimeFullFade:    10 * time.Second,
    loop:            0,
	}
  color := colorful.FastHappyColor()
	for i, _ := range r.Leds {
		r.Leds[i].SetColor(color)
  }
  if DEBUG_LVL >= 1 {fmt.Println("RgbFader initialized all leds")}
	return r
}

func (r *RgbFader) Tick(start time.Time, o interface{}) {
  var color colorful.Color
  var  h, c, l float64
  // duration := time.Since(start)
	// nrOfSteps := int((H_MAX - H_MIN) / r.ColorDifference)
	// frequency := r.TimeFullFade / time.Duration(nrOfSteps)

  //fmt.Print(frequency)
  //fmt.Println("   duration % frequency == ", duration % frequency)
  //fmt.Print("   Led1 ", "is ", r.Leds[1].CIELCH.H, "   colorDifference is ", r.ColorDifference, " NewValue is ")

	if  true {
		for i, _ := range r.Leds {
      if r.Leds[i].CIELCH.H + float32(r.ColorDifference) >= H_MAX {
       	r.Leds[i].CIELCH.H = H_MIN + (H_MAX - r.Leds[i].CIELCH.H)
      }
			color = r.Leds[i].Color()
      h, c, l = color.Hsv()
      color = colorful.Hsv(h + r.ColorDifference, c, l)
      r.Leds[i].SetColor(color)

		}
	}
  r.loop ++
//  fmt.Println(r.Leds[1].CIELCH.H, " loop: ", r.loop)
  r.Consumer.Tick(time.Since(start), r.Leds[:])
}
