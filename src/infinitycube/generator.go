package main


import (
    "net"
    "math/rand"
    "io"
    "fmt"
    "time"
    "github.com/lucasb-eyer/go-colorful"
)

type Led struct {
    Color colorful.Color
}

type Consumer interface {
    Tick(time.Duration, interface{})
}

const LEDS = EDGE_LENGTH * EDGES_PER_SIDE * NR_OF_SIDES

type Generator struct {
    Consumer
    Offset int
    Length int
    Direction int
    Leds [LEDS*2]Led
}


func NewGenerator() *Generator {
    g := &Generator{Length: EDGE_LENGTH * 2, Direction:1}
    for i, _ := range g.Leds {
        if i % g.Length == 0 {
            g.Leds[i].Color = colorful.FastHappyColor()
        } else {
            g.Leds[i].Color = colorful.Color{0, 0, 0}
        }
    }
    return g
}

func (g *Generator) Tick(d time.Duration, o interface{}) {
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

type CubeX struct {
    io.ReadWriter
    buffer [3 * EDGE_LENGTH * EDGES_PER_SIDE * NR_OF_SIDES]byte
}

func NewCubeX() (c *CubeX, err error) {
	socketCon, err := net.Dial("tcp", "192.168.1.222:12345")
	if err != nil {
		return
    }
    c = &CubeX{ReadWriter: socketCon}
    return
}

func (c *CubeX) Tick(d time.Duration, o interface{}) {
    leds := o.([]Led)

    h := 0
    for i, _ := range leds {
        c.buffer[h+0], c.buffer[h+1], c.buffer[h+2] = leds[i].Color.RGB255()
        h += 3
    }

	var startByte [1]byte
    n, _ := c.Read(startByte[:])
    if(n == 1) {
        c.Write(c.buffer[:])
    }
}

type RandomTicker struct {
    Consumer
    Threshold float32
}

func (r *RandomTicker) Tick(d time.Duration, o interface{}) {
    v := rand.Float32()
    r.Consumer.Tick(d, v < r.Threshold)
}

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

const (
    fps_target = 30
    fps_duration = time.Second / fps_target
)

func MakeWorld() (err error) {
    //g := NewGenerator()
    //r := &RandomTicker{Threshold: .05}
    //i := &IntervalTicker{Interval: 1 * time.Second / 2 / EDGE_LENGTH}
    myRgbFader := NewRgbFader()
    bf := &DirtyBlurFilter{}
    c, err := NewCubeX()
    if err != nil {
        fmt.Print(err)
        return
    }

    //r.Consumer = g
    //i.Consumer = g
    myRgbFader.Consumer = bf
    bf.Consumer = c

    starttime := time.Now()
    for {
        a := time.Now()

        //i.Tick(a.Sub(starttime), true)
        myRgbFader.Tick(starttime, nil)


        b := time.Now()
        elapsed := b.Sub(a)
        time.Sleep(fps_duration - elapsed)
    }
}
