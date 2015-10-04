package main

import (
    "github.com/gordonklaus/portaudio"
    "github.com/mjibson/go-dsp/spectral"
    "fmt"
    "time"
)
const (
    SAMPLE_RATE = 44100
    FRAMES_PER_BUFFER = 512
    NUM_CHANNELS = 1
)

type SAMPLE float32

type paTestData struct {
    frameIndex int
    maxFrameIndex int
    buffer64 []float64
    recordedSamples []SAMPLE
    spektralDensity []float64
    freqs []float64
    }

func NewPaTestData() (*paTestData) {
    d := &paTestData{}
    d.buffer64 = make([]float64, FRAMES_PER_BUFFER)
    d.recordedSamples = make([]SAMPLE, FRAMES_PER_BUFFER)
    d.spektralDensity = make([]float64, FRAMES_PER_BUFFER / 2 + 1)
    d.freqs = make([]float64, FRAMES_PER_BUFFER / 2 + 1)
    d.maxFrameIndex = FRAMES_PER_BUFFER
    d.frameIndex = 0
    return d
}

func StartSoundTracking() (*paTestData){
    var err error
    var streamParameter portaudio.StreamParameters

    data := NewPaTestData()
    data.frameIndex = 0 //prevent "not used" error
    portaudio.Initialize()

    streamParameter.Input.Device, err = portaudio.DefaultInputDevice()
    CheckErr(err)
    streamParameter.Input.Channels = NUM_CHANNELS //mono
    streamParameter.Input.Latency = streamParameter.Input.Device.DefaultLowInputLatency //not necessary?
    streamParameter.Output.Device = nil  //input only
    streamParameter.SampleRate = SAMPLE_RATE
    streamParameter.FramesPerBuffer = FRAMES_PER_BUFFER

    fmt.Println("Input Device is:", streamParameter.Input.Device)

    stream, err := portaudio.OpenStream(streamParameter, data.RecordCallback)
    CheckErr(err)
    err = stream.Start()
    CheckErr(err)
    fmt.Println("Now recording")
    time.Sleep(500 * time.Millisecond) //I dont know why this would be neccessary but it somehow is...
    return data
}


func (pa *paTestData) RecordCallback(buffer []SAMPLE) {
    pwelchOptions := spectral.PwelchOptions{NFFT: FRAMES_PER_BUFFER}
    pa.recordedSamples = buffer
    for i := 0; i < len(buffer) - 1; i++ {
        pa.buffer64[i] = float64(buffer[i])
    }
    pa.spektralDensity, pa.freqs = spectral.Pwelch(pa.buffer64, SAMPLE_RATE, &pwelchOptions)
}

func CheckErr(err error) {
    if err != nil {
        fmt.Println(err)
    }
}
