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
	col.S = 255
	col.H = 180

	cc := &ConstantColor{
		myEffect: eff,
		ColorPar: col}

	return cc
}
func (cc *ConstantColor) Colorize(leds []Led) []Led {

	var colLeds []Led

	for _, led := range leds {

		var colLed Led
		colLed.S = 1
		colLed.H = 60

		// dont touch brightness/value
		colLed.V = led.V

		colLeds = append(colLeds, colLed)
	}

	return colLeds
}
