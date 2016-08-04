package main

import (
	"github.com/lucasb-eyer/go-colorful"
)

type HsvFade struct {
	Effect
	ColorDifference float64
	TimeFullFadePar int //in Seconds
	fpsTarget       int
}

func NewHsvFade(disp Display, fps int) *HsvFade {
	ef := NewEffect(disp, 0.5, 0.0)

	r := &HsvFade{
		Effect:          ef,
		TimeFullFadePar: 20,
		fpsTarget:       fps}

	for i := r.OffsetPar; i < (r.LengthPar + r.OffsetPar); i++ {
		r.Leds[i].Color = colorful.Color{0, 255, 0}
	}
	r.ColorDifference = (float64(H_MAX-H_MIN) / float64(r.TimeFullFadePar*r.fpsTarget))
	return r
}

func (r *HsvFade) Update() {
	r.ColorDifference = (float64(H_MAX-H_MIN) / float64(r.TimeFullFadePar*r.fpsTarget))

	var h float64
	for i := r.OffsetPar; i < (r.LengthPar + r.OffsetPar); i++ {
		h, _, _ = r.Leds[i].Color.Hsv()
		r.Leds[i].Color = colorful.Hsv(h+r.ColorDifference, 1, 1)
		r.Leds[i].CheckColor()
	}

	r.myDisplay.AddPattern(r.Leds, r.ColorOpacity, r.BlackOpacity)
}
