/*
This package controlls a ws2812b led strip. The idea is to divide the strip
in segments which can be controlled seperatly. This will make it possible to
play (different) light effects on different segments of the strip at the same
time. In my setup the led strip is going to be mounted on a cube. This is why
there are variables like "cube, side, edge"...nevertheless these are basicly just
just better sounding names for "led-strip, segment, smaller segment". So there
should be no trouble adjusting the code for a different led setup.
*/

package main

import (
	"flag"
    "log"
    "net/http"
    "time"
)

const DEBUG_LVL = 4
/*
0 no debug information
1 general information 
2 side information
3 edge information
4 led information
*/

var (
    serial_port    = flag.String("serial", "", "serial port")
    unix_socket    = flag.String("unixconnect", "", "connect to unix socket")
    listen_address = flag.String("listen", ":2500", "http service address")
    static_path    = flag.String("static", "static", "path to the static content")
)

func main() {
	var err error

    cube := NewCube()
    cube.RGBiteration()

    go func() {
        for {
            cube.side[0].edge[0].led[0].printRGB()
            time.Sleep(1000 * time.Millisecond)
        }
    }()

	http.Handle("/status", cube)
    http.Handle("/", http.FileServer(http.Dir(*static_path)))

    err = http.ListenAndServe(*listen_address, nil)
    if err != nil {
        log.Fatalf("ListenAndServe failed: %v", err)
    }
}



type Status2 struct {
    Cups  int    
    NoChangeSince [2]string
    NrOfCups [2]int
    History [2]map[string]float64
    HitOrder []int
    Error string
}