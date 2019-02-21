package main

import (
	"time"
)

func GreenBinaryWheel(eH *EffectHandler, playTime time.Duration) map[Effector]time.Duration {
	cc1 := NewConstantColor(0.8, 130)
	sd := NewBinaryWheel(eH.myDisplay, cc1, EDGE_LENGTH, 300*time.Millisecond)

	effectMap := map[Effector]time.Duration{}
	effectMap[sd] = playTime

	return effectMap
}

func GoldenStarDust(eH *EffectHandler, playTime time.Duration) map[Effector]time.Duration {
	cc1 := NewConstantColor(0.77, 47)
	sd := NewStarDust(eH.myDisplay, cc1, eH.audio)

	effectMap := map[Effector]time.Duration{}
	effectMap[sd] = playTime

	return effectMap
}

func StaticWhite(eH *EffectHandler, playTime time.Duration) map[Effector]time.Duration {
	cc1 := NewConstantColor(0, 0)
	solid := NewSolidBrightness(eH.myDisplay, cc1, 1.0)

	effectMap := map[Effector]time.Duration{}
	effectMap[solid] = playTime

	return effectMap
}

func RedSunsetStarDust(eH *EffectHandler, playTime time.Duration) map[Effector]time.Duration {

	greenish := Color{R: 52, G: 92, B: 125} //(53, 92, 125)
	blueish := Color{R: 108, G: 91, B: 123} //(108, 91, 123)
	redish := Color{R: 192, G: 108, B: 132} // (192, 108, 132)

	colors := make([]Color, 5)
	colors[0] = redish
	colors[1] = blueish
	colors[2] = greenish
	colors[3] = blueish
	colors[4] = redish

	colGrad := NewColorGradient(colors, 2*EDGE_LENGTH)
	solid := NewSolidBrightness(eH.myDisplay, colGrad, 1.0)

	cc1 := NewConstantColor(0.60, 187)
	sd := NewStarDust(eH.myDisplay, cc1, eH.audio)

	effectMap := map[Effector]time.Duration{}
	effectMap[solid] = playTime
	effectMap[sd] = playTime

	return effectMap
}

func MagmaPlasma(eH *EffectHandler, playTime time.Duration) map[Effector]time.Duration {
	// sine 1
	cc1 := NewConstantColor(1, 0)
	sine1 := NewSine(eH.myDisplay, cc1)
	sine1.Frequency = 2 * NR_OF_SIDES
	sine1.SetLoopTime(5 * NR_OF_SIDES)

	// sine 2
	cc2 := NewConstantColor(1, 30)
	sine2 := NewSine(eH.myDisplay, cc2)
	sine2.Frequency = 3 * NR_OF_SIDES
	sine2.SetLoopTime(7 * NR_OF_SIDES)

	// multi running light
	cc3 := NewConstantColor(1, 60)
	mrl := NewMultiRunningLight(eH.myDisplay, cc3)

	magmaPlasma := map[Effector]time.Duration{}

	magmaPlasma[sine1] = playTime
	magmaPlasma[sine2] = playTime
	magmaPlasma[mrl] = playTime

	return magmaPlasma
}

func CellularAutomataMonochrome(eH *EffectHandler, playTime time.Duration) map[Effector]time.Duration {
	cc1 := NewConstantColor(1, 0)
	ca := NewCellularAutomata(eH.myDisplay, cc1)

	automata := map[Effector]time.Duration{}
	automata[ca] = playTime

	return automata
}

func CellularAutomatagGradient(eH *EffectHandler, playTime time.Duration) map[Effector]time.Duration {
	blue := Color{R: 63, G: 43, B: 150}
	white := Color{R: 168, G: 192, B: 255}

	colors := make([]Color, 2)
	colors[0] = blue
	colors[1] = white

	colGrad := NewColorGradient(colors, 2*EDGE_LENGTH)
	ca := NewCellularAutomata(eH.myDisplay, colGrad)

	automata := map[Effector]time.Duration{}
	automata[ca] = playTime

	return automata
}

func LinearSpectrumMonochrome(eH *EffectHandler, playTime time.Duration) map[Effector]time.Duration {
	cc1 := NewConstantColor(1, 0)
	spec := NewLinearSpectrum(eH.myDisplay, cc1, eH.audio)

	effectMap := map[Effector]time.Duration{}
	effectMap[spec] = playTime

	return effectMap
}

func LinearEdgeSpectrumMonochrome(eH *EffectHandler, playTime time.Duration) map[Effector]time.Duration {
	cc1 := NewConstantColor(0, 0)
	spec := NewLinearEdgeSpectrum(eH.myDisplay, cc1, eH.audio)

	effectMap := map[Effector]time.Duration{}
	effectMap[spec] = playTime

	return effectMap
}

func EdgeVolumeRedGreen(eH *EffectHandler, playTime time.Duration) map[Effector]time.Duration {

	green := Color{R: 0, G: 255, B: 0}
	red := Color{R: 255, G: 0, B: 0}

	colors := make([]Color, 2)
	colors[0] = green
	colors[1] = red

	colGrad := NewColorGradient(colors, EDGE_LENGTH)
	vol := NewEdgeVolume(eH.myDisplay, colGrad, eH.audio)

	effectMap := map[Effector]time.Duration{}
	effectMap[vol] = playTime

	return effectMap
}

func MultiRunningLightHSV(eH *EffectHandler, playTime time.Duration) map[Effector]time.Duration {
	hsv := NewHsvFade(60, 0)
	mrl := NewMultiRunningLight(eH.myDisplay, hsv)

	effectMap := map[Effector]time.Duration{}
	effectMap[mrl] = playTime

	return effectMap
}
