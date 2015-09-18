package main

import(
  "io"
  "net"
  "time"
  //"github.com/lucasb-eyer/go-colorful"
  //"fmt"
)

type PreCube struct {
    leds [LEDS]Led
    colorOpacity float64
    blackOpacity float64
}

func NewPreCube(newLeds [LEDS]Led, cOp, bOp float64) (pc *PreCube) {
    pc = &PreCube{leds: newLeds, colorOpacity: cOp, blackOpacity: bOp}
    return
}

type Cube struct {
    io.ReadWriter
    preCubes []PreCube
    finalCube [LEDS]Led
    buffer [3 * EDGE_LENGTH * EDGES_PER_SIDE * NR_OF_SIDES]byte
}

func NewCube() (c *Cube, err error) {
	socketCon, err := net.Dial("tcp", "192.168.1.222:12345")
	if err != nil {
		return
    }
    c = &Cube{ReadWriter: socketCon}
    return
}

func (c *Cube) Tick(d time.Duration, o interface{}) {

}

func (c *Cube) AddPreCube(leds [LEDS]Led, colorOpacity float64, blackOpacity float64) {
    pc := NewPreCube(leds, colorOpacity, blackOpacity)
    c.preCubes = append(c.preCubes, *pc)
}

func (c *Cube) renderCube(){
    c.MergePreCubes()
    //leds := o.([]Led)

    leds := c.finalCube
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

func (c *Cube) resetPreCubes() {
    c.preCubes = nil
}

func (c *Cube) MergePreCubes() {
    for i, _ := range c.preCubes {
        for p := 0; p < LEDS; p++ {
            if i == 0 { //we dont want to blend the first PreCube with the still black "finalCube"
                c.finalCube[p] =  c.preCubes[i].leds[p]
            }else{ //and later we blend all folowing PreCubes in
                if c.preCubes[i].leds[p].Color.R == 0 && c.preCubes[i].leds[p].Color.G == 0 && c.preCubes[i].leds[p].Color.B == 0 {
                    c.finalCube[p].Color = c.finalCube[p].Color.BlendRgb(c.preCubes[i].leds[p].Color, c.preCubes[i].blackOpacity)
                }else{
                    c.finalCube[p].Color = c.finalCube[p].Color.BlendRgb(c.preCubes[i].leds[p].Color, c.preCubes[i].colorOpacity)
                }
            }
        }
    }
}
