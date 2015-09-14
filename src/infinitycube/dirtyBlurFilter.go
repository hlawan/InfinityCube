package main

import (
    "time"
)

type DirtyBlurFilter struct {
    Consumer
    Leds [LEDS]Led
}

func idx(i, o int) int {
    i += o
    if i < 0 {
        i += LEDS
    }
    if i >= LEDS {
        i -= LEDS
    }
    return i
}

func (f *DirtyBlurFilter) Tick(d time.Duration, o interface{}) {
    leds := o.([]Led)

    for i, _ := range leds {
        s := .02
        c := leds[i].Color
        c = c.BlendHsv(leds[idx(i, -2)].Color, s/4)
        c = c.BlendHsv(leds[idx(i, -1)].Color, s)
        c = c.BlendHsv(leds[idx(i,  1)].Color, s)
        c = c.BlendHsv(leds[idx(i,  2)].Color, s/4)
        f.Leds[i].Color = c
    }
    f.Consumer.Tick(d, f.Leds[:])
}
