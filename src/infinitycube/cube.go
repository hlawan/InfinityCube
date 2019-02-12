package main

import (
	"log"
	"math"

	"github.com/kellydunn/go-opc"
)

type Cube struct {
	fadeCandy *opc.Client
	Effects   []Effect
	finalCube []Color //all Effects are merged to one finalCube, which then will be sent to the real cube
}

func NewCube(server string, nrOfLeds int) (c *Cube, err error) {
	oc := opc.NewClient()
	err = oc.Connect("tcp", server)
	if err != nil {
		log.Fatal("Could not connect to Fadecandy server", err)
	}

	c = &Cube{
		fadeCandy: oc,
		finalCube: make([]Color, nrOfLeds)}

	return
}

func (c *Cube) NrOfLeds() (nrOfLeds int) {
	return len(c.finalCube)
}

func (c *Cube) render() {
	c.MergeEffects()
	c.resetEffects()
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

func (c *Cube) AddEffect(newEffect Effect) {
	c.Effects = append(c.Effects, newEffect)
}

func (c *Cube) resetEffects() {
	c.Effects = nil
}

func (c *Cube) MergeEffects() {

	// start from black
	for i, _ := range c.finalCube {
		c.finalCube[i].reset()
	}

	// merge all generated patterns
	for i := range c.Effects {
		for p := 0; p < LEDS; p++ {
			c.finalCube[p] = blendLeds(c.finalCube[p], c.Effects[i].Leds[p])
		}
	}
}

func blendLeds(col1, col2 Color) Color {
	var newCol Color

	r1, g1, b1 := col1.RGBfromHSV()
	r2, g2, b2 := col2.RGBfromHSV()

	// merge
	nR := uint8((float64(r1) + float64(r2)) / 2.0)
	nG := uint8((float64(g1) + float64(g2)) / 2.0)
	nB := uint8((float64(b1) + float64(b2)) / 2.0)

	// set new color
	newCol.FromRGB(nR, nG, nB)

	// restore brightness
	nV := math.Max(col1.V, col2.V)
	newCol.setV(nV)

	return newCol
}
