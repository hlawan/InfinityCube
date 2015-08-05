package main

import (
	"net"
	"time"
)

func startSocketComunication(c *Cube) {
	socketCon, _ := net.Dial("unix", "/tmp/so")
	for {
		socketCon.Write([]byte("hi"))
		time.Sleep(1e9)
	}
}

func (my *Cube) parseLEDstatus(leds *[EDGE_LENGTH * EDGES_PER_SIDE * NR_OF_SIDES]int) {
	return
}
