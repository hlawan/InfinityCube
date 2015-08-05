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
	//"time"
	"fmt"
	"net"
	//"os"
	"reflect"
	"time"
)

const DEBUG_LVL = 1

/*
0 no debug information
1 general information
2 side information
3 edge information
4 led information
*/

var (
	//serial_port    = flag.String("serial", "", "serial port")
	//serial_port    = flag.String("serial", "/tmp/so", "serial port")
	//unix_socket    = flag.String("unixconnect", "", "connect to unix socket")
	listen_address = flag.String("listen", ":2500", "http service address")
	static_path    = flag.String("static", "static", "path to the static content")
)

func main() {
	var err error

	cube := NewCube()
	cube.RGBiteration()
	cube.simpleRunningLight(255, 0, 255)
	//cube.side[2].edge[1].simpleRunningLight(0,255,0)

	fooType := reflect.TypeOf(Cube{})
	for i := 0; i < fooType.NumMethod(); i++ {
		method := fooType.Method(i)
		fmt.Println(method.Name)
	}
	fmt.Println("Before socket stuff...")
	
	 c,err := net.Dial("unix", "/tmp/so")
         for {
       		c.Write([]byte("hi\n"))
         	time.Sleep(1e9)
	 }

	http.Handle("/status", cube)
	http.Handle("/", http.FileServer(http.Dir(*static_path)))

	err = http.ListenAndServe(*listen_address, nil)
	if err != nil {
		log.Fatalf("ListenAndServe failed: %v", err)
	}
}
