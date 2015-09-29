package main

import (
    "github.com/gordonklaus/portaudio"
    "fmt"
)
const (
    SAMPLE_RATE = 44100
    FRAMES_PER_BUFFER = 512
    NUM_SECONDS = 5
    NUM_CHANNELS = 2
)

type SAMPLE float32

type paTestData struct{
    frameIndex int
    maxFrameIndex int
    recordedSamples []SAMPLE
    }

func getCrazy(){
    err := portaudio.Initialize()
    if err != nil {
        fmt.Println(err)
    }

    var stream portaudio.StreamParameters
    var data paTestData
    var totalFrames, numSamples int
    var max, val SAMPLE
    var average float64

    totalFrames = NUM_SECONDS * SAMPLE_RATE
    data.maxFrameIndex = totalFrames
    data.frameIndex = 0
    numSamples = totalFrames * NUM_CHANNELS
    data.recordedSamples = make([]SAMPLE, numSamples)

    for i := 0; i < numSamples; i++ {
        data.recordedSamples[i] = 0
    }

    stream.Input.Device, err = portaudio.DefaultInputDevice()
    if err != nil {
        fmt.Println(err)
    }
    stream.Output.Device, err = portaudio.DefaultOutputDevice()
    if err != nil {
        fmt.Println(err)
    }
    stream.SampleRate = SAMPLE_RATE
    stream.FramesPerBuffer = FRAMES_PER_BUFFER

    fmt.Println("Default Audio Device is:")
    fmt.Println(portaudio.DefaultInputDevice())
    fmt.Println("Sound stuff done")



    portaudio.Terminate()
}
