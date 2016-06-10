package main

import (
  "log"
	"github.com/kellydunn/go-opc"
	"github.com/lucasb-eyer/go-colorful"
)

type Cube struct {
	fadeCandy *opc.Client
	preCubes  []PreCube    //every effect generator adds a PreCube
	finalCube [LEDS]Led    //all PreCubes are merged to one finalCube, which then will be sent to the real cube
	buffer    [3 * EDGE_LENGTH * EDGES_PER_SIDE * NR_OF_SIDES]byte
}

func NewCube(server string) (c *Cube, err error) {

	// Create a client
	oc := opc.NewClient()
	err = oc.Connect("tcp", server)
	c = &Cube{fadeCandy: oc}
	if err != nil {
		log.Fatal("Could not connect to Fadecandy server", err)
	}

	return
}

func (c *Cube) renderCube() {
	c.MergePreCubes()
	leds := c.finalCube

	// send pixel data
	m := opc.NewMessage(0)
	m.SetLength(uint16(len(leds) * 3))

	for i := range leds {
		r, g, b := leds[i].Color.RGB255()
		m.SetPixelColor(i, r, g, b)
	}

	err := c.fadeCandy.Send(m)
	if err != nil {
		log.Println("couldn't send color", err)
	}

}

type PreCube struct {
	leds         [LEDS]Led
	colorOpacity float64
	blackOpacity float64
}

func NewPreCube(newLeds [LEDS]Led, cOp, bOp float64) (pc *PreCube) {
	pc = &PreCube{leds: newLeds, colorOpacity: cOp, blackOpacity: bOp}
	return
}

func (c *Cube) AddPreCube(leds [LEDS]Led, colorOpacity float64, blackOpacity float64) {
	pc := NewPreCube(leds, colorOpacity, blackOpacity)
	c.preCubes = append(c.preCubes, *pc)
}

func (c *Cube) resetPreCubes() {
	c.preCubes = nil
}

func (c *Cube) MergePreCubes() {
	black := colorful.Color{0, 0, 0}
	for i := range c.preCubes {
		for p := 0; p < LEDS; p++ {
			if i == 0 { //we dont want to merge the first PreCube with the still black "finalCube"
				c.finalCube[p] = c.preCubes[i].leds[p]
			} else { //and later we merge all folowing PreCubes
				if c.preCubes[i].leds[p].Color == black {
					c.finalCube[p].Color = c.finalCube[p].Color.BlendRgb(c.preCubes[i].leds[p].Color, c.preCubes[i].blackOpacity)
				} else {
					c.finalCube[p].Color = c.finalCube[p].Color.BlendRgb(c.preCubes[i].leds[p].Color, c.preCubes[i].colorOpacity)
				}
			}
		}
	}
}
