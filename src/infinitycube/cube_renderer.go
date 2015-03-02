package main

import (
	"time"
    "fmt"
)

func (my *Cube) RGBiteration() {

	for i := 0; i < NR_OF_SIDES; i++ {
        my.side[i].renderer = my.side[i].setSide
    }

	go func() {
        if(DEBUG_LVL >=1) {fmt.Println("RGBiteration go routine started")}
    	for {
        	for i := 0; i < 3; i++ {
        		switch i {
   					case 0:
        				for i := 0; i < NR_OF_SIDES; i++ {my.side[i].renderer(255,0,0)}
    				case 1:
        				for i := 0; i < NR_OF_SIDES; i++ {my.side[i].renderer(0,255,0)}
    				case 2:
        				for i := 0; i < NR_OF_SIDES; i++ {my.side[i].renderer(0,0,255)}
    			}
    			time.Sleep(1000 * time.Millisecond)
        	}
        }
 	}()
}

