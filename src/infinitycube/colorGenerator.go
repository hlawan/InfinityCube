// colorGenerator
package main

import (
	"fmt"

	"github.com/lucasb-eyer/go-colorful"
)

type ColorGenerator interface {
	Colorize(*[]Led)
}

type ConstantColor struct {
	myEffect *Effect
	ColorPar colorful.Color
}

func NewConstantColor(eff *Effect) *ConstantColor {

	cc := &ConstantColor{
		myEffect: eff,
		ColorPar: colorful.Hsv(120, 1, 1)}

	return cc
}
func (cc *ConstantColor) Colorize(leds *[]Led) {
	for _, led := range cc.myEffect.Leds {
		// get value (brightness) from current pattern
		_, _, v := led.Color.Hsv()
		// get saturation and hue (color) to set
		h, s, _ := cc.ColorPar.Hsv()
		// set
		//fmt.Println(h, " ", s, " ", v)
		fmt.Println(colorful.Hsv(h, s, v).RGB255())
		newColor := colorful.Hsv(h, s, v)
		led.Color = newColor
	}
}
