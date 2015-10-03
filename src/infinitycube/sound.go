package main

import (
    "github.com/gordonklaus/portaudio"
    "fmt"
)
const (
    SAMPLE_RATE = 44100
    FRAMES_PER_BUFFER = 512
    NUM_SECONDS = 5
    NUM_CHANNELS = 1
)

type SAMPLE int32

type paTestData struct{
    frameIndex int
    maxFrameIndex int
    recordedSamples []SAMPLE
    }

func getCrazy(){
    var err error
    var streamParameter portaudio.StreamParameters
    var data paTestData
    var totalFrames, numSamples int

    totalFrames = NUM_SECONDS * SAMPLE_RATE
    data.maxFrameIndex = totalFrames
    data.frameIndex = 0
    numSamples = totalFrames * NUM_CHANNELS
    data.recordedSamples = make([]SAMPLE, numSamples)
    if data.recordedSamples == nil {
        fmt.Println("recorded samples received nullpointer")
        return
    }
    for i := 0; i < numSamples; i++ { //not necessary?
        data.recordedSamples[i] = 0
    }

    portaudio.Initialize()
    streamParameter.Input.Device, err = portaudio.DefaultInputDevice()
    if err != nil {
        fmt.Println(err)
    }
    streamParameter.Input.Channels = NUM_CHANNELS //mono
    streamParameter.Input.Latency = streamParameter.Input.Device.DefaultLowInputLatency //not necessary?
    streamParameter.Output.Device = nil  //input only
    streamParameter.SampleRate = SAMPLE_RATE
    streamParameter.FramesPerBuffer = FRAMES_PER_BUFFER
    fmt.Println("Input Device is:", streamParameter.Input.Device)

    stream, err := portaudio.OpenStream(streamParameter, data.RecordCallback)
    if err != nil {
        fmt.Println(err)
    }
    err = stream.Start()
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("Now recording")
    StartWebServer(&data)
    portaudio.Terminate()
}

func (pa *paTestData) RecordCallback(buffer []SAMPLE){
    pa.recordedSamples = buffer
}
