package main

import (
	"log"
	"math"

	"github.com/kellydunn/go-opc"
)

const MAX_POWER_SINGLE_LED = 0.05 * 5 // 50 mA * 5V
const MAX_POWER_WATTS = 80
const PARALLEL_STRIPES = 1

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
	c.scaleFinalPattern(MAX_POWER_WATTS, PARALLEL_STRIPES)
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

// scaleFinalPattern estimates the requested power by the lightning pattern, stored in c.finalCube
// given the max allowed power maxPowerWatts and the number of parallelly controlles stripes parallelStripes
// the function downscales (darkens) the pattern, if the requested power is too high.
func (c *Cube) scaleFinalPattern(maxPowerWatts float64, parallelStripes uint32) {
	// calculate requested power
	requestedPowerWatts := 0.0
	powerOuputFactor := 1.0
	leds := c.finalCube
	for i, _ := range leds {
		power := (float64(leds[i].R)/255.0 + float64(leds[i].G)/255.0 + float64(leds[i].B)/255.0) / 3.0 * MAX_POWER_SINGLE_LED
		requestedPowerWatts += power
	}
	requestedPowerWatts *= float64(parallelStripes)
	// scale pattern
	if requestedPowerWatts > maxPowerWatts {
		powerOuputFactor = maxPowerWatts / requestedPowerWatts
		log.Printf("Requested Power: %v, max allowed: %v\n", requestedPowerWatts, maxPowerWatts)
		log.Printf("Downscale power with factor %v.\n", powerOuputFactor)
		for i, _ := range leds {
			c.finalCube[i].R = uint8(powerOuputFactor * float64(leds[i].R))
			c.finalCube[i].G = uint8(powerOuputFactor * float64(leds[i].G))
			c.finalCube[i].B = uint8(powerOuputFactor * float64(leds[i].B))
		}
	}
}
