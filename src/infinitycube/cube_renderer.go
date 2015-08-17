package main

import (
	"time"
    "fmt"
)

func (my *Cube) RGBiteration() {
    const speed = 1000 //effekt speed: RGB-cycle time in ms

	for i := 0; i < NR_OF_SIDES; i++ {
        my.side[i].renderer = my.side[i].setSide
    }

	go func() {
        if(DEBUG_LVL >=1) {fmt.Println("RGBiteration go routine started")}
    	for {
        	for i := 0; i < 3; i++ {
        		switch i {
   					case 0:
        				for i := 0; i < NR_OF_SIDES; i++ {my.side[i].renderer(255,10,10)}
    				case 1:
        				for i := 0; i < NR_OF_SIDES; i++ {my.side[i].renderer(10,255,10)}
    				case 2:
        				for i := 0; i < NR_OF_SIDES; i++ {my.side[i].renderer(10,10,255)}
    			}
    			time.Sleep(speed * time.Millisecond)
        	}
        }
 	}()
}


func (my *Cube) simpleRunningLight(red, green, blue uint8) {
   for i := 0; i < NR_OF_SIDES; i++ {
        for o := 0; o < EDGES_PER_SIDE; o++ {
            my.side[i].edge[o].simpleRunningLight(red, green, blue)
        }
    }
}    

func (my *Cube) growingRunningLight(red, green, blue uint8) {
    const speed = 333 //effekt speed: delay before start on next edge
    go func() {
        for{
            for i := 0; i < NR_OF_SIDES; i++ {
                for o := 0; o < EDGES_PER_SIDE; o++ {
                    my.side[i].edge[o].simpleRunningLight(red, green, blue)
                    time.Sleep(speed * time.Millisecond)
                }
            }
        }
    }()
}

