package main

import(
  "io"
  "net"
  "time"
)

type Cube struct {
    io.ReadWriter
    buffer [3 * EDGE_LENGTH * EDGES_PER_SIDE * NR_OF_SIDES]byte
}

func NewCube(addr string) (c *Cube, err error) {
	socketCon, err := net.Dial("tcp", addr)
	if err != nil {
		return
    }
    c = &Cube{ReadWriter: socketCon}
    return
}

func (c *Cube) Tick(d time.Duration, o interface{}) {
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
