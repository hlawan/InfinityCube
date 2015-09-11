package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

func startSocketComunication(myCube *Cube) {
	//fmt.Println(myCube.parseLEDstatus())
	socketCon, err := net.Dial("unix", "/tmp/so")

	if err != nil {
		log.Fatalf("startSocketComunication failed %v", err)
	} else {
		if DEBUG_LVL >= 1 {
			fmt.Println("dialed unix socket...")
		}
	}

	startByte := make([]byte, 1)
	go func() {
		for {
			//fmt.Println("Before Read...")
			n, _ := socketCon.Read(startByte)
        		if(n == 1){
				//fmt.Println("Before Write...")
				binary.Write(socketCon, binary.LittleEndian, myCube.parseLEDstatus())
				//fmt.Println(myCube.parseLEDstatus())
				//limit to 30 frames per second
				time.Sleep((1000 * time.Millisecond)/30)
			}
		}
	}()
}

/*parseLEDstatus creates an int array of the rgb values (time three) of every
*	single led. This arry is passed through the unix domain socket and used by the
*	c code that outouts it to the led strip*/
func (my *Cube) parseLEDstatus() (leds [3 * EDGE_LENGTH * EDGES_PER_SIDE * NR_OF_SIDES]byte) {
	//time.Sleep(500 * time.Millisecond)
	//for h := 0; h < ((3 * EDGE_LENGTH * EDGES_PER_SIDE * NR_OF_SIDES) - 2); h++ {
		//if h%3 == 0 {
			h := 0
			for i := 0; i < NR_OF_SIDES; i++ {
				for o := 0; o < EDGES_PER_SIDE; o++ {
					for p := 0; p < EDGE_LENGTH; p++ {
						leds[h] = (my.side[i].edge[o].led[p].Red)
						leds[h+1] = (my.side[i].edge[o].led[p].Green)
						leds[h+2] = (my.side[i].edge[o].led[p].Blue)
						h += 3
					}
				}
			}
		//}
	//}
	return
}
