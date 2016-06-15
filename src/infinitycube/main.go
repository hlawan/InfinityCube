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
	"fmt"
	"time"
	"reflect"
	"github.com/fatih/structs"
)

const (
	debugLvl       = 1
	fpsTarget      = 100
	fpsDuration    = time.Second / fpsTarget
	EDGE_LENGTH    = 14 //in my setup there are always 14 leds in a row
	EDGES_PER_SIDE = 4  //well for me its a square...so 4
	NR_OF_SIDES    = 6  //regular cube => 6 sides
	LEDS           = EDGE_LENGTH * EDGES_PER_SIDE * NR_OF_SIDES
	H_MAX          = 360 // maximum Hue value (Hsv)
	H_MIN          = 0   // minimum Hue value (Hsv)
)

/*
0 no debug information
1 general information
2 side information
3 edge information
4 led information
*/

var (
	serverPtr      = flag.String("fcserver", "192.168.1.222:7890", "Fadecandy server and port to connect to")
	serial_port    = flag.String("serial", "/dev/zero", "serial port")
	listen_address = flag.String("listen", ":2500", "http service address")
	static_path    = flag.String("static", "static", "path to the static content")
)

func main() {
	flag.Parse()

	c, err := NewCube(*serverPtr, LEDS)
	if err != nil {
		fmt.Print(err)
		return
	}

	eH := NewEffectHandler(c, fpsTarget)

	//rawSoundData := StartSoundTracking()
	//audio := StartAudioProcessing(rawSoundData)
	//webInterface := StartWebServer(audio)
  eH.addCellularAutomata(1.0, 0.0, .3, 152)
	fmt.Println(reflect.TypeOf(eH.Effects[0]))
	fmt.Println(structs.Names(eH.Effects[0]))
	for{eH.render()}
}

/*
	//initializing generators, cubes, filters...
	myHsvFader := NewHsvFader(0, LEDS, 10, .20, 0)
	brl := NewBinaryRunningLight(LEDS, 1, .5, 0)
	rl0 := NewRunningLight(violett, 2*EDGE_LENGTH, 0.001, .5, 0)
	rl1 := NewRunningLight(blue, 3*EDGE_LENGTH, 0.002, .5, 0)
	grl0 := NewGausRunningLight(redish, 2*EDGE_LENGTH, 5, .5, 0)
	grl1 := NewGausRunningLight(red, 1*EDGE_LENGTH, 2, .5, 0)
	grl2 := NewGausRunningLight(violett, 4*EDGE_LENGTH, 11, .5, 0)
	eq := NewEqualizer(0, LEDS, 1, 0, audio)
	//cA := NewCellularAutomata(1, 0, 152, .5)
	i0 := &IntervalTicker{Interval: 10 * time.Microsecond / 2 / EDGE_LENGTH}



	//combining all parts as liked

	//r.Consumer = g
	//i0.Consumer = brl
	myHsvFader.Consumer = c
	brl.Consumer = c
	rl0.Consumer = c
	rl1.Consumer = c
	grl0.Consumer = c
	grl1.Consumer = c
	grl2.Consumer = c
	//bf.Consumer = c
	eq.Consumer = c
	//cA.Consumer = c

	i0.Consumer = rl1

	//main loop

	var elapsedTime, sleepingTime [200]time.Duration
	var elapsed, slept time.Duration
	var z time.Time
	var selector int
	o := 0
	starttime := time.Now()
	for {
		a := time.Now()

		if webInterface.clapSelect {
			selector = audio.clapCount % 7
		} else {
			selector = webInterface.selectedEffect
			//selector = 6;
		}

		switch selector { //audio.clapCount % 7
		case 0:
			myHsvFader.Tick(starttime, nil)
		case 1:
			brl.Tick(a.Sub(starttime), true)
		case 2:
			rl0.Tick(a.Sub(starttime), true)
			rl1.Tick(a.Sub(starttime), true)
		case 3:
			grl0.Tick(a.Sub(starttime), true)
			grl1.Tick(a.Sub(starttime), true)
			grl2.Tick(a.Sub(starttime), true)
		case 4:
			eq.EdgeVolume()
		case 5:
			eq.WhiteSpectrum()
		case 6:
			//eq.WhiteEdgeSpectrum()
			//cA.Update()
		}

		c.render()

		b := time.Now()
		elapsed = b.Sub(a)
		time.Sleep(fpsDuration - elapsed)

		//only needed for FPS calculation
		if false {
			z = time.Now()
			sleepingTime[o] = fpsDuration - elapsed
			elapsedTime[o] = z.Sub(a)
			if o > 198 {
				totalTime := 0 * time.Second
				currentFps := 0 * time.Second
				for p := 1; p < 200; p++ {
					slept += sleepingTime[p]
					totalTime += elapsedTime[p]
				}
				slept /= 199
				totalTime /= 199
				currentFps = (1 * time.Second / totalTime)
				sleepPercent := (100 * time.Millisecond / totalTime) * slept

				if debugLvl > 0 {
					fmt.Println("-->loop time:", totalTime,
						"-->FPS:", currentFps.Nanoseconds(),
						"-->I slept for:", slept,
						" (", sleepPercent.Seconds()*1000, "%)")
				}
				o = 0
			}
			o++
		}
	}
	*/
