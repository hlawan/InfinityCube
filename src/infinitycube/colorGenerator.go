// colorGenerator
package main

import (
	"math"
)

type ColorGenerator interface {
	Colorize([]Led) []Led
	Update()
}

func correctGamma(inputValue float64) float64 {
	return math.Pow(inputValue, 0.5)
}

type ConstantColor struct {
	ColorPar Led
}

func NewConstantColor(Saturation float64, Hue uint16) *ConstantColor {

	var col Led
	col.S = Saturation
	col.H = Hue

	cc := &ConstantColor{ColorPar: col}

	return cc
}

func (cc *ConstantColor) Update() {

}

func (cc *ConstantColor) Colorize(leds []Led) []Led {

	var colLeds []Led

	for _, led := range leds {

		var colLed Led
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

func (hsv *HsvFade) Colorize(leds []Led) []Led {

	var colLeds []Led

	for _, led := range leds {

		var colLed Led
		colLed.S = 1.0
		colLed.H = uint16(math.Round(hsv.angle))

		// Gamma correction
		colLed.V = correctGamma(led.V)

		colLeds = append(colLeds, colLed)
	}

	return colLeds
}