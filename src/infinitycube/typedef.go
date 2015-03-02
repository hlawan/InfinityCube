/*
Creating different types and methodes in order to divide the led-strip in 
segments. These segments are orgnaized hierarchally. Since my led-strip is 
going to be mounted in a cube The hierachy will be something like: 
cube (has six) sides (has four) edges (has 14) leds...
*/

package main

import "fmt"

const EDGE_LENGTH = 14 //in my setup there are always 14 leds in a row
const EDGES_PER_SIDE = 4 //well for me its a square...so 4
const NR_OF_SIDES = 6 //regular cube => 6 sides

//the smalest segment is one single led
type led struct { 
	Red uint8
	Green uint8
	Blue uint8
}

func (my *led) setRGB(red, green, blue uint8) {			//dot or not?
	if(DEBUG_LVL >=4) {fmt.Println("setRGB called with Arguments", red, green, blue)}
	my.Red = red
	my.Green = green
	my.Blue = blue
}

func (my *led) printRGB() {								//dot or not?
	fmt.Println("RGB values:(", my.Red, my.Green, my.Blue,")")
}

type RGBrenderer func()
type sideRenderer func(uint8,uint8,uint8)
type edgeRenderer func(uint8,uint8,uint8)

//a couple of leds in a row are an edge
type edge struct {
	renderer edgeRenderer
	led [EDGE_LENGTH]led
	/*edgeRenderer()...*/
}

//one side of my cube is framed by four edges...it's a square ;)
type side struct {
	renderer sideRenderer
	edge [EDGES_PER_SIDE]edge
	/*sideRenderer()...*/
}

//and finally a cube is six sides glued together
type Cube struct{
	renderer RGBrenderer 
	side [NR_OF_SIDES]side 
}

func NewCube() (c *Cube){
	c = &Cube{}
	return
}

