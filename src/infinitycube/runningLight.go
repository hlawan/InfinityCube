package main

import (
    "math/rand"
    //"fmt"
    "math"
    "time"
    "github.com/lucasb-eyer/go-colorful"
)

type Consumer interface {
    Tick(time.Duration, interface{})
}

type RunningLight struct {
    Consumer
    colorful.Color
    Position float64
    Delta float64
    Length int
    Leds [LEDS]Led
}

func NewRunningLight(color colorful.Color, length int, delta float64) *RunningLight {
    r := &RunningLight{
        Color: color,
        Length: length,
        Delta: delta,
    }
    return r
}

var BLACK = colorful.Color{0, 0, 0}

func max(a, b float64) float64 {
    if a < b {
        return b
    }
    return a
}

func min(a, b float64) float64 {
    if a < b {
        return a
    }
    return b
}

func dist(a, b float64) float64 {
    return math.Abs(a-b)
}

func (r *RunningLight) Tick(d time.Duration, o interface{}) {
    advance := o.(bool)

    if advance {
        r.Position += r.Delta
        if r.Position > 1 {
            r.Position -= 1
        }

        pos := r.Position * float64(r.Length)
        for i, _ := range r.Leds {
            j := i % r.Length
            r.Leds[i].Color = BLACK.BlendRgb(r.Color, 1 - min(1, dist(pos, float64(j))))
        }
    }

    r.Consumer.Tick(d, r.Leds[:])
}

type BinaryRunningLight struct {
    Consumer
    Offset int
    Length int
    Direction int
    Leds [LEDS*2]Led
}


func NewBinaryRunningLight() *BinaryRunningLight {
    g := &BinaryRunningLight{Length: EDGE_LENGTH * 2, Direction:1}
    for i, _ := range g.Leds {
        if i % g.Length == 0 {
            g.Leds[i].Color = colorful.FastHappyColor()
        } else {
            g.Leds[i].Color = colorful.Color{0, 0, 0}
        }
    }
    return g
}

func (g *BinaryRunningLight) Tick(d time.Duration, o interface{}) {
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
    g.Consumer.Tick(d, g.Leds[g.Offset:g.Offset+LEDS])
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
