package main


import (
    "github.com/lucasb-eyer/go-colorful"
)

type Led struct {
    Color colorful.Color
    //opacity (toDo)
    //position (maybe nice to say something like "light up all corners...")
}

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
