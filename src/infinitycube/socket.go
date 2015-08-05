package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

func startSocketComunication(myCube *Cube) {
	fmt.Println(myCube.parseLEDstatus())
	socketCon, err := net.Dial("unix", "/tmp/so")

	if err != nil {
		log.Fatalf("startSocketComunication failed %v", err)
	} else {
		if DEBUG_LVL >= 1 {
			fmt.Println("dialed unix socket...")
		}
	}

	buf := new(bytes.Buffer)
	go func() {
		binary.Write(buf, binary.LittleEndian, myCube.parseLEDstatus())
		socketCon.Write(buf.Bytes())
		time.Sleep(30 * time.Millisecond)
	}()
}

/*parseLEDstatus creates an int array of the rgb values (time three) of every
*	single led. This arry is passed through the unix domain socket and used by the
*	c code that outouts it to the led strip*/
func (my *Cube) parseLEDstatus() (leds [3 * EDGE_LENGTH * EDGES_PER_SIDE * NR_OF_SIDES]int) {
	time.Sleep(500 * time.Millisecond)
	for h := 0; h < ((3 * EDGE_LENGTH * EDGES_PER_SIDE * NR_OF_SIDES) - 2); h++ {
		if h%3 == 0 {
			for i := 0; i < NR_OF_SIDES; i++ {
				for o := 0; o < EDGES_PER_SIDE; o++ {
					for p := 0; p < EDGE_LENGTH; p++ {
						leds[h] = int(my.side[i].edge[o].led[p].Red)
						leds[h+1] = int(my.side[i].edge[o].led[p].Green)
						leds[h+2] = int(my.side[i].edge[o].led[p].Blue)
					}
				}
			}
		}
	}
	return
}
