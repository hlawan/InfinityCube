package main

import (
	"math"
)

type Led struct {
	H       uint16  /*0..360*/
	S, V    float64 /*0..255*/
	R, G, B uint8   /*0..255*/
}

func (a *Led) RGB() (r, g, b uint8) {
	// Direct implementation of the graph in this image:
	// https://en.wikipedia.org/wiki/HSL_and_HSV#/media/File:HSV-RGB-comparison.svg

	C := a.V * a.S
	segment := a.H / 60
	X := C * (1 - math.Abs(float64(segment%2)-1))

	var r1, g1, b1 float64

	switch segment {
	case 0:
		r1 = C
		g1 = X
		b1 = 0
	case 1:
		r1 = X
		g1 = C
		b1 = 0
	case 2:
		r1 = 0
		g1 = C
		b1 = X
	case 3:
		r1 = 0
		g1 = X
		b1 = C
	case 4:
		r1 = X
		g1 = 0
		b1 = C
	case 5:
		r1 = C
		g1 = 0
		b1 = X
	}

	m := a.V - C
	r1 += m
	g1 += m
	b1 += m

	r = uint8(r1 * 255)
	g = uint8(g1 * 255)
	b = uint8(b1 * 255)

	return r, g, b
}

//func (a *Led) RGB() (r, g, b uint8) {
//	// Direct implementation of the graph in this image:
//	// https://en.wikipedia.org/wiki/HSL_and_HSV#/media/File:HSV-RGB-comparison.svg
//	max := uint32(a.V) * 255
//	min := uint32(a.V) * uint32(255-a.S)

//	a.H %= 360
//	segment := a.H / 60
//	offset := uint32(a.H % 60)
//	mid := ((max - min) * offset) / 60

//	switch segment {
//	case 0:
//		return max, min + mid, min
//	case 1:
//		return max - mid, max, min
//	case 2:
//		return min, max, min + mid
//	case 3:
//		return min, max - mid, max
//	case 4:
//		return min + mid, min, max
//	case 5:
//		return max, min, max - mid
//	}

//	return 0, 0, 0
//}

func (a *Led) Off() bool {
	if a.V > 0 {
		return false
	} else {
		return true
	}
}
