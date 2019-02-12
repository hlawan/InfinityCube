package main

func StaticGradientGreenRed(eH *EffectHandler, playTime time.Duration) map[Effector]time.Duration {

	green := Color{R: 0, G: 255, B: 0}
	red := Color{R: 255, G: 0, B: 0}

	colors := make([]Color, 2)
	colors[0] = green
	colors[1] = red

	colGrad := NewColorGradient(colors, EDGE_LENGTH)
	solid := NewSolidBrightness(eH.myDisplay, colGrad, 1.0)

	effectMap := map[Effector]time.Duration{}
	effectMap[solid] = playTime

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
