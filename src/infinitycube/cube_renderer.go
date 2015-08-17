package main

import (
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
	"time"
)

func (my *Cube) RGBiteration() {
	const speed = 5000 //effekt speed: RGB-cycle time in ms

	for i := 0; i < NR_OF_SIDES; i++ {
		my.side[i].renderer = my.side[i].setSide
	}

	go func() {
		if DEBUG_LVL >= 1 {
			fmt.Println("RGBiteration go routine started")
		}
		for {
			for i := 0; i < 3; i++ {
				switch i {
				case 0:
					for i := 0; i < NR_OF_SIDES; i++ {
						my.side[i].renderer(255, 10, 10)
					}
				case 1:
					for i := 0; i < NR_OF_SIDES; i++ {
						my.side[i].renderer(10, 255, 10)
					}
				case 2:
					for i := 0; i < NR_OF_SIDES; i++ {
						my.side[i].renderer(10, 10, 255)
					}
				}
				time.Sleep(speed * time.Millisecond)
			}
		}
	}()
}

func (my *Cube) simpleRunningLight(red, green, blue uint8) {
	for i := 0; i < NR_OF_SIDES; i++ {
		for o := 0; o < EDGES_PER_SIDE; o++ {
			my.side[i].edge[o].simpleRunningLight(red, green, blue)
		}
	}
}

func (my *Cube) growingRunningLight(red, green, blue uint8) {
	const speed = 333 //effekt speed: delay before start on next edge
	go func() {
		for {
			for i := 0; i < NR_OF_SIDES; i++ {
				for o := 0; o < EDGES_PER_SIDE; o++ {
					my.side[i].edge[o].simpleRunningLight(red, green, blue)
					time.Sleep(speed * time.Millisecond)
				}
			}
		}
	}()
}

func (my *Cube) fade() {
	const speed = 50
	bottomReached := false

	go func() {
		for {
			for i := 0; i < NR_OF_SIDES; i++ {
				for o := 0; o < EDGES_PER_SIDE; o++ {
					for p := 0; p < EDGE_LENGTH; p++ {
						if my.side[i].edge[o].led[p].Blue > 0 && bottomReached == false {
							my.side[i].edge[o].led[p].Blue -= 3
						} else {
							bottomReached = true
							my.side[i].edge[o].led[p].Blue += 3
						}
						if my.side[i].edge[o].led[p].Blue == 255 {
							bottomReached = false
						}
					}
				}
			}
			time.Sleep(speed * time.Millisecond)
		}
	}()
}

func (my *Cube) rainbowFade() {
	const speed = 50
	color := colorful.Hcl(0, 0, 0)
	var h float64

	go func() {
		for {
			for h = 0; h < 360; {
				for i := 0; i < NR_OF_SIDES; i++ {
					for o := 0; o < EDGES_PER_SIDE; o++ {
						for p := 0; p < EDGE_LENGTH; p++ {
							color = colorful.Hcl(h, 1, 0.3).Clamped()
							my.side[i].edge[o].led[p].Red = uint8(color.R * 255)
							my.side[i].edge[o].led[p].Green = uint8(color.G * 255)
							my.side[i].edge[o].led[p].Blue = uint8(color.B * 255)
							h += 0.01
						}
					}
				}
				time.Sleep(speed * time.Millisecond)
			}
		}
	}()
}

func (my *Cube) output() {
	for i := 0; i < NR_OF_SIDES; i++ {
		for o := 0; o < 1; o++ {
			for p := 0; p < EDGE_LENGTH; p++ {
				fmt.Print(my.side[i].edge[o].led[p].Red/10,
					my.side[i].edge[o].led[p].Green/10,
					my.side[i].edge[o].led[p].Blue/10, "|")
			}
			fmt.Println()
		}
	}
}
