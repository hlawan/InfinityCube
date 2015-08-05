package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func startSocketComunication() {
	socketCon, err := net.Dial("unix", "/tmp/so")
	if err != nil {
		log.Fatalf("startSocketComunication failed %v", err)
	} else {
		if DEBUG_LVL >= 1 {
			fmt.Println("dialed unix socket...")
		}
	}
	for {
		socketCon.Write([]byte("hi\n"))
		time.Sleep(1e9)
	}
}

func (my *Cube) parseLEDstatus(leds *[EDGE_LENGTH * EDGES_PER_SIDE * NR_OF_SIDES]int) {
	return
}
