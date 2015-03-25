package main

import (
    "fmt"
    "time"
)

func (my *edge) setEdge(red, green, blue uint8) {
	if(DEBUG_LVL >=3) {fmt.Println("setEdge() called with Arguments", red, green, blue)}
	for i := 0; i < EDGE_LENGTH; i++ {
		my.led[i].setRGB(red, green, blue)
	}
}

func (my *edge) turnOffEdge() {
	if(DEBUG_LVL >=3) {fmt.Println("turnOffEdge() called")}
	for i := 0; i < EDGE_LENGTH; i++ {
		my.led[i].setRGB(0,0,0)
	}
}

func (my *edge) simpleRunningLight(red, green, blue uint8) {
	if(DEBUG_LVL >=3) {fmt.Println("simpleMovingLight() called with Arguments", red, green, blue)}
	const speed = 1000 //ms from first to last led
	my.turnOffEdge()
	go func() {
		for{
			for i := 0; i < EDGE_LENGTH; i++ {
				my.led[i].setRGB(red, green, blue)
				if (i > 0){
					my.led[i-1].setRGB(0, 0, 0)
				}else{
					my.led[EDGE_LENGTH-1].setRGB(0, 0, 0)
				}
				time.Sleep(speed/EDGE_LENGTH * time.Millisecond)
			}	
		}
	}()
}

//----------------------new stuff------------------//

