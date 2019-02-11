// colorGenerator
package main

import (
	"fmt"
	"math"
)

type ColorGenerator interface {
	Colorize([]Color) []Color
	Update()
}

func correctGamma(inputValue float64) float64 {
	return math.Pow(inputValue, 0.5)
}

type ConstantColor struct {
	ColorPar Color
}

func NewConstantColor(Saturation float64, Hue uint16) *ConstantColor {

	var col Color
	col.S = Saturation
	col.H = Hue

	cc := &ConstantColor{ColorPar: col}

	return cc
}

func (cc *ConstantColor) Update() {

}

func (cc *ConstantColor) Colorize(leds []Color) []Color {

	var colLeds []Color

	for _, led := range leds {

		var colLed Color
		colLed.S = cc.ColorPar.S
		colLed.H = cc.ColorPar.H

		// Gamma correction
		colLed.V = correctGamma(led.V)

		colLeds = append(colLeds, colLed)
	}

	return colLeds
}

type HsvFade struct {
	TimeFullFade float64
	angle        float64
	delta        float64
}

func NewHsvFade(timeFullFade, initAngle float64) *HsvFade {

	hsv := &HsvFade{
		TimeFullFade: timeFullFade,
		angle:        initAngle}

	hsv.delta = 360.0 / (hsv.TimeFullFade * float64(fpsTarget))

	return hsv
}

func (hsv *HsvFade) Update() {
	hsv.angle += hsv.delta
	hsv.angle = math.Mod(hsv.angle, 360)
}

func (hsv *HsvFade) Colorize(leds []Color) []Color {

	var colLeds []Color

	for _, led := range leds {

		var colLed Color
		colLed.S = 1.0
		colLed.H = uint16(math.Round(hsv.angle))

		// Gamma correction
		colLed.V = correctGamma(led.V)

		colLeds = append(colLeds, colLed)
	}

	return colLeds
}

type ColorGradient struct {
	Elements []Color
	length   int
}

func NewColorGradient(colorElements []Color, gradientLength int) *ColorGradient {

	if len(colorElements) < 2 {
		fmt.Println("not enough colors for color gradient")
		return nil
	}

	cg := &ColorGradient{
		Elements: colorElements,
		length:   gradientLength}

	return cg
}

func (cg *ColorGradient) Update() {

}

func (grad *ColorGradient) Colorize(leds []Color) []Color {

	var colLeds []Color

	for i, led := range leds {

		colLed := linearGradient(i%EDGE_LENGTH, EDGE_LENGTH, grad.Elements[0], grad.Elements[1])

		// Gamma correction
		colLed.V = correctGamma(led.V)

		colLeds = append(colLeds, colLed)
	}

	return colLeds
}

func linearGradient(ledPosition, length int, col1, col2 Color) Color {

	x := float64(ledPosition)
	xMax := float64(length)

	r1, g1, b1 := col1.RGB()
	r2, g2, b2 := col2.RGB()

	r1f := float64(r1)
	g1f := float64(g1)
	b1f := float64(b1)
	r2f := float64(r2)
	g2f := float64(g2)
	b2f := float64(b2)

	nR := uint8(r1f + (x/xMax)*(r2f-r1f))
	nG := uint8(g1f + (x/xMax)*(g2f-g1f))
	nB := uint8(b1f + (x/xMax)*(b2f-b1f))

	var newCol Color
	newCol.FromRGB(nR, nG, nB)

	return newCol
}
