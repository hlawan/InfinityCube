package main


import (
    "github.com/lucasb-eyer/go-colorful"
)

type Led struct {
    Color colorful.Color
    //opacity (toDo)
    //position (maybe nice to say something like "light up all corners...")
}

var (
    white = colorful.Color{1, 1, 1}
    violett = colorful.Color{.5, 0, .5}
    redish = colorful.Color{.8, .1, .3}
    red =  colorful.Color{.3, 0, 0}
    black= colorful.Color{0, 0, 0}
    blue = colorful.Color{0, 0, 1}
    green = colorful.Color{0, 1, 0}
)

func (l *Led) CheckColor() {
  if l.Color.R <= 1e-5{
    l.Color.R = 0
  }
  if l.Color.G <= 1e-5{
    l.Color.G = 0
  }
  if l.Color.B <= 1e-5{
    l.Color.B = 0
  }
}

func (a *Led) OnOrOff() bool{
  if (a.Color.R + a.Color.G + a.Color.B) < 0.0001 {
    return false
  }else{
    return true
  }
}
