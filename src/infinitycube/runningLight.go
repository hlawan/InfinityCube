package main

import (
    "math/rand"
    //"fmt"
    "time"
    "github.com/lucasb-eyer/go-colorful"
)

type Consumer interface {
    Tick(time.Duration, interface{})
    AddPreCube([LEDS]Led, float64, float64)
}

type RunningLight struct {
    Consumer
    Offset int
    Length int
    Direction int
    ColorOpacity	float64
    BlackOpacity	float64
    Leds [LEDS*2]Led
}


func NewRunningLight(length, direction int, colorOpacity, blackOpacity  float64) *RunningLight {
    g := &RunningLight {
        Length: length,
        Direction:direction,
        ColorOpacity: colorOpacity,
        BlackOpacity: blackOpacity}
    for i, _ := range g.Leds {
        if i % g.Length == 0 {
            g.Leds[i].Color = colorful.FastHappyColor()
        } else {
            g.Leds[i].Color = colorful.Color{0, 0, 0}
        }
    }
    return g
}

func (g *RunningLight) Tick(d time.Duration, o interface{}) {
    advance := o.(bool)
    if advance {
        g.Offset += g.Direction
        if g.Offset < 0 {
            g.Offset += g.Length
        }
        if g.Offset > g.Length {
            g.Offset -= g.Length
        }
    }

    var cube [LEDS]Led
    for i, v := range g.Leds[g.Offset:g.Offset+LEDS] {
        cube[i] = v
    }
    //g.Consumer.Tick(d, g.Leds[g.Offset:g.Offset+LEDS])
    g.Consumer.AddPreCube(cube, g.ColorOpacity, g.BlackOpacity)
}

type IntervalTicker struct {
    Consumer
    Last time.Duration
    Interval time.Duration
}

func (i *IntervalTicker) Tick(d time.Duration, o interface{}) {
    fire := false
    if d - i.Last > i.Interval {
        fire = true
        i.Last = d
    }
    i.Consumer.Tick(d, fire)
}

type RandomTicker struct {
    Consumer
    Threshold float32
}

func (r *RandomTicker) Tick(d time.Duration, o interface{}) {
    v := rand.Float32()
    r.Consumer.Tick(d, v < r.Threshold)
}
