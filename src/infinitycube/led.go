package main

import (
	//"fmt"
	"math"
)

type Led struct {
	H       uint16  /*0..360*/
	S, V    float64 /*0..1*/
	R, G, B uint8   /*0..255*/
}

func (a *Led) RGB() (r, g, b uint8) {
	// Direct implementation of the graph in this image:
	// https://en.wikipedia.org/wiki/HSL_and_HSV#/media/File:HSV-RGB-comparison.svg

	a.fixRanges()

	C := a.V * a.S
	segment := float64(a.H) / 60.0
	X := C * (1 - math.Abs(math.Mod(segment, 2)-1))

	var r1, g1, b1 float64

	switch uint8(segment) {
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

	a.R = uint8(r1 * 255)
	a.G = uint8(g1 * 255)
	a.B = uint8(b1 * 255)

	return a.R, a.G, a.B
}

func (a *Led) FromRGB(rInt, gInt, bInt uint8) {
	// from https://www.rapidtables.com/convert/color/rgb-to-hsv.html

	a.R = rInt
	a.G = gInt
	a.B = bInt

	r := float64(rInt) / 255.0
	g := float64(gInt) / 255.0
	b := float64(bInt) / 255.0

	cMax := math.Max(r, math.Max(g, b))
	cMin := math.Min(r, math.Min(g, b))

	delta := cMax - cMin

	// calc Hue
	nH := 0.0

	if delta > 0 {
		if cMax == r {
			nH = (g - b) / delta
		} else if cMax == g {
			nH = ((b - r) / delta) + 2
		} else if cMax == b {
			nH = ((r - g) / delta) + 4
		}
	}
	a.H = uint16(60*nH) % 360

	// calc Saturation
	if cMax > 0 {
		a.S = delta / cMax
	} else {
		a.S = 0
	}

	// calc Value
	a.V = cMax
}

func (a *Led) Off() bool {
	if a.V > 1e-3 {
		return false
	} else {
		return true
	}
}

func (a *Led) fixRanges() {
	if a.V > 1 {
		a.V = 1
	}
	if a.V < 0 {
		a.V = 0
	}

	if a.S > 1 {
		a.S = 1
	}
	if a.S < 0 {
		a.S = 0
	}

	a.H = a.H % 360

}

func (a *Led) reset() {
	a.H = 0
	a.S = 0
	a.V = 0
	a.R = 0
	a.G = 0
	a.B = 0
}

func (a *Led) setV(nV float64) {
	a.V = nV

	// after Setting V: RGB has to be updated as well
	a.R, a.G, a.B = a.RGB()
}
