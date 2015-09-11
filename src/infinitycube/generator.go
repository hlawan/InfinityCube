package main


import (
    "net"
    "math/rand"
    "io"
    "fmt"
    "time"
    "github.com/lucasb-eyer/go-colorful"
)

type CIELCH struct {
    H float32
    C float32
    L float32
}

type Led struct {
    CIELCH
}

func (led *Led) SetColor(color colorful.Color) {
    h, c, l := color.Hcl()
    led.CIELCH = CIELCH{float32(h), float32(c), float32(l)}
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
    g := &Generator{Length: EDGE_LENGTH, Direction:1}
    color := colorful.FastHappyColor()
    color = colorful.Color{0, 1, 0}
    for i, _ := range g.Leds {
        if i % g.Length == 0 {
            g.Leds[i].SetColor(color)
        } else {
            g.Leds[i].SetColor(colorful.Color{0, 0, 0})
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
    for _, l := range leds {
        color := colorful.Hcl(float64(l.CIELCH.H), float64(l.CIELCH.C), float64(l.CIELCH.L))
        c.buffer[h+0] = uint8(255 * color.R)
        c.buffer[h+1] = uint8(255 * color.G)
        c.buffer[h+2] = uint8(255 * color.B)
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

const (
    fps_target = 30
    fps_duration = time.Second / fps_target
)

func MakeWorld() (err error) {
    g := NewGenerator()
    r := &RandomTicker{Threshold: .05}
    i := &IntervalTicker{Interval: 2 * time.Second / EDGE_LENGTH}
    c, err := NewCubeX()
    if err != nil {
        fmt.Print(err)
        return
    }

    r.Consumer = g
    i.Consumer = g
    g.Consumer = c

    starttime := time.Now()
    for {
        a := time.Now()

        i.Tick(a.Sub(starttime), true)

        b := time.Now()
        elapsed := b.Sub(a)
        time.Sleep(fps_duration - elapsed)
    }
}
