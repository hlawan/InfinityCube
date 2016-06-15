package main

import (
	"github.com/kellydunn/go-opc"
	"github.com/lucasb-eyer/go-colorful"
	"log"
	"time"
)

type Cube struct {
	fadeCandy *opc.Client
	Patterns  []Pattern //every effect generator adds a Pattern
	finalCube []Led     //all Patterns are merged to one finalCube, which then will be sent to the real cube
}

type Consumer interface {
	Tick(time.Duration, interface{})
	AddPattern([LEDS]Led, float64, float64)
}

func NewCube(server string, nrOfLeds int) (c *Cube, err error) {
	// Create a client
	oc := opc.NewClient()
	err = oc.Connect("tcp", server)
	c = &Cube{
		fadeCandy: oc,
		finalCube: make([]Led, nrOfLeds)}

	if err != nil {
		log.Fatal("Could not connect to Fadecandy server", err)
	}
	return
}

func (c *Cube) NrOfLeds() (nrOfLeds int) {
	return len(c.finalCube)
}

func (c *Cube) render() {
	c.MergePatterns()
	c.resetPatterns()
}

func (c *Cube) Show() {
  c.render()
	leds := c.finalCube

	// send pixel data
	m := opc.NewMessage(0)
	m.SetLength(uint16(len(leds) * 3)) // *3 -> r, g, b

	for i := range leds {
		r, g, b := leds[i].Color.RGB255()
		m.SetPixelColor(i, r, g, b)
	}

	err := c.fadeCandy.Send(m)
	if err != nil {
		log.Println("couldn't send color", err)
	}
}

func (c *Cube) Tick(d time.Duration, o interface{}) {
	//moved stuff to renderCube()
}

type Pattern struct {
	leds         []Led
	colorOpacity float64
	blackOpacity float64
}

func NewPattern(newLeds []Led, cOp, bOp float64) (pc *Pattern) {
	pc = &Pattern{
		leds:         newLeds,
		colorOpacity: cOp,
		blackOpacity: bOp}
	return
}

func (c *Cube) AddPattern(leds []Led, colorOpacity float64, blackOpacity float64) {
	pc := NewPattern(leds, colorOpacity, blackOpacity)
	c.Patterns = append(c.Patterns, *pc)
}

func (c *Cube) resetPatterns() {
	c.Patterns = nil
}

func (c *Cube) MergePatterns() {
	black := colorful.Color{0, 0, 0}
	for i := range c.Patterns {
		for p := 0; p < LEDS; p++ {
			if i == 0 { //we dont want to merge the first Pattern with the still black "finalCube"
				c.finalCube[p] = c.Patterns[i].leds[p]
			} else { //and later we merge all folowing Patterns
				if c.Patterns[i].leds[p].Color == black {
					c.finalCube[p].Color = c.finalCube[p].Color.BlendRgb(c.Patterns[i].leds[p].Color, c.Patterns[i].blackOpacity)
				} else {
					c.finalCube[p].Color = c.finalCube[p].Color.BlendRgb(c.Patterns[i].leds[p].Color, c.Patterns[i].colorOpacity)
				}
			}
		}
	}
}
