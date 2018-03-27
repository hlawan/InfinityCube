package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
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
	for i := range leds {
		fmt.Println(leds[i].Color)
		//y := leds[i].Color.Y
		cb := leds[i].Color.Cb
		cr := leds[i].Color.Cr
		r, g, b := color.YCbCrToRGB(0, cb, cr)
		fmt.Println(r, g, b)
		//r8 := uint8(255 * (r / (2 ^ 32)))
		//g8 := uint8(255 * (g / (2 ^ 32)))
		//b8 := uint8(255 * (b / (2 ^ 32)))

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
	for i := range c.Patterns {
		for p := 0; p < LEDS; p++ {
			if i == 0 { //we dont want to merge the first Pattern with the still black "finalCube"
				c.finalCube[p] = c.Patterns[i].leds[p]
			} else { //and later we merge all folowing Patterns
				//				if c.Patterns[i].leds[p].OnOrOff() {
				c.finalCube[p].Color = blendYCbCr(c.finalCube[p].Color, c.Patterns[i].leds[p].Color)
				//				} else {
				//					c.finalCube[p].Color = c.finalCube[p].Color.BlendRgb(c.Patterns[i].leds[p].Color, c.Patterns[i].blackOpacity)
				//				}
			}
		}
	}
}

func blendYCbCr(col1, col2 color.NYCbCrA) color.NYCbCrA {
	//	nY = (col1.Y + col2.Y) / 2
	nY := math.Max(float64(col1.Y), float64(col2.Y))
	nCb := (float64(col1.Cb) + float64(col2.Cb)) / 2
	nCr := (float64(col1.Cr) + float64(col2.Cr)) / 2
	nA := (float64(col1.A) + float64(col2.A)) / 2

	var col color.NYCbCrA
	col.Y = uint8(nY)
	col.Cb = uint8(nCb)
	col.Cr = uint8(nCr)
	col.A = uint8(nA)

	return col
}
