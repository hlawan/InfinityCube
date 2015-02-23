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
	r uint8
	g uint8
	b uint8
}

func (my led) setRGB(red, green, blue uint8) {
	my.r = red
	my.g = green
	my.b = blue
}

func (my led) printRGB() {
	fmt.Printf("RGB values:(%d,%d,%d)", my.r, my.g, my.b)
}

//a couple of leds in a row are an edge
type edge struct {
	led [EDGE_LENGTH]led
	/*edgeRenderer()...*/
}

//one side of my cube is framed by sour edges...it's a square ;)
type side struct {
	edge [EDGES_PER_SIDE]edge
	/*sideRenderer()...*/
}

//and finally a cube is six sides glued together
type cube struct{
	side [NR_OF_SIDES]side 
	/*cubeRenderer()...*/
}
