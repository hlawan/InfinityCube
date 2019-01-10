package main

import (
	// "fmt"
	//	"image/color"
	"log"
	//	"math"
	"time"

	"github.com/kellydunn/go-opc"
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
	//fmt.Println(leds)

	for i, led := range leds {
		m.SetPixelColor(i, led.R, led.G, led.B)
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

	// start from black
	for i, _ := range c.finalCube {
		c.finalCube[i].reset()
	}

	// merge all generated patterns
	for i := range c.Patterns {
		for p := 0; p < LEDS; p++ {
			c.finalCube[p] = blendLeds(c.finalCube[p], c.Patterns[i].leds[p])
		}
	}
}

func blendLeds(col1, col2 Led) Led {
	var newCol Led

	if col1.Off() {
		newCol = col2
		newCol.R, newCol.G, newCol.B = col2.RGB()
	} else if col2.Off() {
		newCol = col1
		newCol.R, newCol.G, newCol.B = col1.RGB()
	} else {
		r1, g1, b1 := col1.RGB()
		r2, g2, b2 := col2.RGB()

		// merge
		nR := uint8((uint16(r1) + uint16(r2)) / 2)
		nG := uint8((uint16(g1) + uint16(g2)) / 2)
		nB := uint8((uint16(b1) + uint16(b2)) / 2)

		// set new color
		newCol.FromRGB(nR, nG, nB)
	}

	return newCol
}
