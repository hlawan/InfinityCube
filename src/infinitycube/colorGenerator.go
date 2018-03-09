// colorGenerator
package main

import (
	//"fmt"
	"image/color"
)

type ColorGenerator interface {
	Colorize(*[]Led) /*(uint8, uint8)*/
}

type ConstantColor struct {
	myEffect *Effect
	ColorPar color.NYCbCrA
}

func NewConstantColor(eff *Effect) *ConstantColor {

	var col color.NYCbCrA
	col.Y = 255
	col.Cb = 0
	col.Cr = 0
	col.A = 255

	cc := &ConstantColor{
		myEffect: eff,
		ColorPar: col}

	return cc
}
func (cc *ConstantColor) Colorize(leds *[]Led) /* (uint8, uint8)*/ {
	for _, led := range *leds {

		led.Color.Cr = cc.ColorPar.YCbCr.Cr
		led.Color.Cb = cc.ColorPar.YCbCr.Cb

		//fmt.Println(led.Color)
	}
	return //cc.ColorPar.YCbCr.Y, cc.ColorPar.YCbCr.Y
}
