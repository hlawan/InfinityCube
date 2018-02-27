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
	"time"
	//"reflect"
	//"github.com/fatih/structs"
)

/*
0 no debug information
1 general information
2 side information
3 edge information
4 led information
*/
const (
	debugLvl       = 1
	fpsTarget      = 100
	fpsDuration    = time.Second / fpsTarget
	EDGE_LENGTH    = 14 //in my setup there are always 14 leds in a row
	EDGES_PER_SIDE = 4  //well for me its a square...so 4
	NR_OF_SIDES    = 6  //regular cube => 6 sides
	LEDS           = EDGE_LENGTH * EDGES_PER_SIDE * NR_OF_SIDES
)

/*
serverPtr 			address of Raspberry running the fadeCandyServer
serial_port
listen_address 		port the webserver is using, leads to webinterface
static_path			folder containg website content
*/
var (
	serverPtr      = flag.String("fcserver", "localhost:7890", "Fadecandy server and port to connect to")
	serial_port    = flag.String("serial", "/dev/zero", "serial port")
	listen_address = flag.String("listen", ":2500", "http service address")
	static_path    = flag.String("static", "static", "path to the static content")
)

func main() {
	flag.Parse()

	/*
		create a Cube with the total number of leds and a connection to the
		dispaying fadeCandyServer whichs is the controlling leds/simulation
	*/
	c, err := NewCube(*serverPtr, LEDS)
	CheckErr(err)

	// connect to microphone
	rawSoundData := StartSoundTracking()
	// start volume analysis, furier analysis, clap detection....
	audio := StartAudioProcessing(rawSoundData)

	/*
		The effectHandler gets a Display to work on -> the cube c. The Handler can
		start/stop effects and shows the comibined result on its display. Also the
		Handler tells the webServer which effects are available and which effects
		are active.
	*/
	eH := NewEffectHandler(c, fpsTarget, audio)

	// the webserver displays the webInterface, where effects can be add/configured
	StartWebServer(audio, eH.availableEffects, eH.effectRequest)

	// from now on the effectHandler renders the effects
	for {
		eH.render()
	}
}
