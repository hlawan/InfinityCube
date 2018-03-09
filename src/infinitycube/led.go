package main

import (
	"image/color"
	//"github.com/lucasb-eyer/go-colorful"
)

type Led struct {
	Color color.NYCbCrA
	//opacity (toDo)
	//position (maybe nice to say something like "light up all corners...")
}

func (a *Led) OnOrOff() bool {
	if a.Color.Y < 1 {
		return false
	} else {
		return true
	}
}
