// colorGenerator
package main

import (
	//"fmt"
	"image/color"
)

type ColorGenerator interface {
	Colorize([]Led) []Led
}

type ConstantColor struct {
	myEffect *Effect
	ColorPar color.NYCbCrA
}

func NewConstantColor(eff *Effect) *ConstantColor {

	var col color.NYCbCrA
	col.Y = 255
	col.Cb = 100
	col.Cr = 100
	col.A = 255

	cc := &ConstantColor{
		myEffect: eff,
		ColorPar: col}

	return cc
}
func (cc *ConstantColor) Colorize(leds []Led) []Led {

	var colLeds []Led

	for _, led := range leds {
		var col color.NYCbCrA
		col.Y = led.Color.YCbCr.Y
		col.Cb = cc.ColorPar.YCbCr.Cb
		col.Cr = cc.ColorPar.YCbCr.Cr
		col.A = cc.ColorPar.A

		var colLed Led
		colLed.Color = col

		colLeds = append(colLeds, colLed)
	}

	return colLeds
}
