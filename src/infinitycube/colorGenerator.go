// colorGenerator
package main

type ColorGenerator interface {
	Colorize([]Led) []Led
}

type ConstantColor struct {
	myEffect *Effect
	ColorPar Led
}

func NewConstantColor(eff *Effect) *ConstantColor {

	var col Led
	col.S = 1
	col.H = 60

	cc := &ConstantColor{
		myEffect: eff,
		ColorPar: col}

	return cc
}
func (cc *ConstantColor) Colorize(leds []Led) []Led {

	var colLeds []Led

	for _, led := range leds {

		var colLed Led
		colLed.S = cc.ColorPar.S
		colLed.H = cc.ColorPar.H

		// dont touch brightness/value
		colLed.V = led.V

		colLeds = append(colLeds, colLed)
	}

	return colLeds
}
