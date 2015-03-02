package main

import (
    "fmt"
)

func (my *edge) setEdge(red, green, blue uint8) {
	if(DEBUG_LVL >=3) {fmt.Println("setEdge called with Arguments", red, green, blue)}
	for i := 0; i < EDGE_LENGTH; i++ {
		my.led[i].setRGB(red, green, blue)
	}
}