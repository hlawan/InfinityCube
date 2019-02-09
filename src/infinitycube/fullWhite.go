package main

import (
	"github.com/lucasb-eyer/go-colorful"
)

type FullWhite struct {
	Effect
}

func NewFullWhite(disp Display) *FullWhite {
	ef := NewEffect(disp, 0.5, 0.0)

	r := &FullWhite{
		Effect: ef}

	for i := r.OffsetPar; i < (r.LengthPar + r.OffsetPar); i++ {
		r.Leds[i].Color = colorful.Color{1, 1, 1}
	}
	return r
}

func (r *FullWhite) Update() {
	for i := r.OffsetPar; i < (r.LengthPar + r.OffsetPar); i++ {
		r.Leds[i].Color = colorful.Color{1, 1, 1}
	}
	r.myDisplay.AddPattern(r.Leds, r.ColorOpacity, r.BlackOpacity)
}
