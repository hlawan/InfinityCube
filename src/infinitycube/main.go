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
	// "log"
	// "net/http"
	"time"
	//"os"
)

const (
		DEBUG_LVL = 1
    fps_target = 60
    fps_duration = time.Second / fps_target
		EDGE_LENGTH = 14 //in my setup there are always 14 leds in a row
		EDGES_PER_SIDE = 4 //well for me its a square...so 4
		NR_OF_SIDES = 6 //regular cube => 6 sides
		LEDS = EDGE_LENGTH * EDGES_PER_SIDE * NR_OF_SIDES
		H_MAX = 360 // maximum Hue value (Hsv)
		H_MIN = 0 // minimum Hue value (Hsv)
)

/*
0 no debug information
1 general information
2 side information
3 edge information
4 led information
*/

var (
	serial_port = flag.String("serial", "", "serial port")
	//serial_port    = flag.String("serial", "/tmp/so", "serial port")
	//unix_socket    = flag.String("unixconnect", "", "connect to unix socket")
	listen_address = flag.String("listen", ":2500", "http service address")
	static_path    = flag.String("static", "static", "path to the static content")
)

func main() {
	//var err error
	//g := NewGenerator()
  //r := &RandomTicker{Threshold: .05}
  //i := &IntervalTicker{Interval: 1 * time.Second / 2 / EDGE_LENGTH}
  myHsvFader := NewHsvFader(0, LEDS, 15)
  bf := &DirtyBlurFilter{}
  c, err := NewCubeX()
  if err != nil {
      fmt.Print(err)
      return
  }

  //r.Consumer = g
  //i.Consumer = g
  myHsvFader.Consumer = bf
  bf.Consumer = c

  var elapsedTime, sleepingTime [200]time.Duration
  var elapsed, slept time.Duration
  var z time.Time
  i := 0
  starttime := time.Now()
  for {
    a := time.Now()

    //i.Tick(a.Sub(starttime), true)
    myHsvFader.Tick(starttime, nil)



    b := time.Now()
    elapsed = b.Sub(a)
    time.Sleep(fps_duration - elapsed)

    //-----------------------------------------------
    //only needed for FPS calculation
    if true {
      z = time.Now()

      sleepingTime[i] = fps_duration - elapsed
      elapsedTime[i] = z.Sub(a)
      if (i > 198) {
        totalTime := 0 * time.Second
        currentFps := 0 * time.Second
        for p:= 1; p < 200; p++ {
          slept += sleepingTime[p]
          totalTime += elapsedTime[p]
        }
        slept /= 199
        totalTime /= 199
        currentFps = (1 * time.Second / totalTime)
        sleepPercent := (100 * time.Millisecond / totalTime) * slept

        if (DEBUG_LVL > 0) {
          fmt.Println("-->loop time:", totalTime,
                      "-->FPS:",currentFps.Nanoseconds(),
                      "-->I slept for:", slept,
                      " (", sleepPercent.Seconds() * 1000, "%)" )
        }
        i = 0
      }
      i++
    }
    //End of FPS calculation
    //-----------------------------------------------
  }

  //  MakeWorld()


	// http.Handle("/status", cube)
	// http.Handle("/", http.FileServer(http.Dir(*static_path)))
	//
	// err = http.ListenAndServe(*listen_address, nil)
	// if err != nil {
	// 	log.Fatalf("ListenAndServe failed: %v", err)
	// }
}
