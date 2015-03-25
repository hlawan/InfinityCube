package main 

import (
	"time"
)

//ticks @ "x" ticks per Second
func hzTicker(ticksPerSec int) (c chan time.Duration) {
    c = make(chan time.Duration)

    go func() {
        for {
            start := time.Now()
            for {
                c <- time.Since(start)
            }
            time.Sleep(time.Second/time.Duration(ticksPerSec))
        }
    }()
    return
}