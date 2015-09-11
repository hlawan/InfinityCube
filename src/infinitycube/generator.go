package main


import (
    "net"
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

type Consumer interface {
    Tick(time.Duration, interface{})
}

const LEDS = EDGE_LENGTH * EDGES_PER_SIDE * NR_OF_SIDES

type Generator struct {
    Consumer
    Length int
    Leds [LEDS*2]Led
}


func NewGenerator() *Generator {
    g := &Generator{Length: EDGE_LENGTH}
    color := colorful.FastHappyColor()
    h, c, l := color.Hcl()
    for i, _ := range g.Leds {
        g.Leds[i].CIELCH = CIELCH{float32(h), float32(c), float32(l)}
        if i % g.Length != 0 {
            g.Leds[i].CIELCH.L = 0
        }
    }
    return g
}

func (g *Generator) Tick(d time.Duration, o interface{}) {
    duration := 10 * time.Second
    offset := int(float32(d % duration) / float32(duration) * float32(g.Length))
    g.Consumer.Tick(d, g.Leds[offset:offset+LEDS])
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


const (
    fps_target = 60
    fps_duration = time.Second / fps_target
)

func MakeWorld() (err error) {
    g := NewGenerator()
    c, err := NewCubeX()
    if err != nil {
        fmt.Print(err)
        return
    }

    g.Consumer = c
    starttime := time.Now()
    for {
        a := time.Now()
        g.Tick(a.Sub(starttime), nil)
        b := time.Now()
        elapsed := b.Sub(a)
        time.Sleep(fps_duration - elapsed)
    }
}
