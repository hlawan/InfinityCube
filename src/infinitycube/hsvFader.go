package main

import (
	"github.com/lucasb-eyer/go-colorful"
	"time"
)

type HsvFader struct {
	Consumer
	FirstLed        int
	Length          int  //Nr of Leds
	ColorDifference float64
	TimeFullFade    int  //in Seconds
	ColorOpacity	float64
	BlackOpacity	float64
	Leds            [LEDS]Led
}

func NewHsvFader(firstLed, length, timeFullFade int, colorOpacity, blackOpacity float64) *HsvFader {
	r := &HsvFader{
		FirstLed: firstLed,
		Length: length,
		TimeFullFade: timeFullFade,
		ColorOpacity: colorOpacity,
		BlackOpacity: blackOpacity}

	for i := r.FirstLed; i < (r.Length + r.FirstLed); i++ {
		r.Leds[i].Color = colorful.Color{0,255,0}
	}
  	r.ColorDifference = (float64(H_MAX - H_MIN) / float64(r.TimeFullFade * fps_target))
	return r
}

func (r *HsvFader) Tick(start time.Time, o interface{}) {
	var h float64
	for i := r.FirstLed; i < (r.Length + r.FirstLed); i++ {
    	h, _, _ = r.Leds[i].Color.Hsv()
    	r.Leds[i].Color = colorful.Hsv(h + r.ColorDifference, 1, 1)
    	r.Leds[i].CheckColor()
	}
  	//r.Consumer.Tick(time.Since(start), r.Leds[:])
	r.Consumer.AddPreCube(r.Leds, r.ColorOpacity, r.BlackOpacity)
}
