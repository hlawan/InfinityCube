package main

import (
    "github.com/gordonklaus/portaudio"
    "github.com/mjibson/go-dsp/spectral"
    "fmt"
    //"time"
    "sync"
)
const (
    SAMPLE_RATE = 44100
    FRAMES_PER_BUFFER = 256
    NUM_CHANNELS = 1
    DEBUG = false
)

type SAMPLE float32

type SoundSingnal struct {
    sync.Mutex
    *portaudio.Stream
    buffer []SAMPLE
    bufferChannel chan []SAMPLE
    // buffer64 []float64
    // spektralDensity []float64
    // freqs []float64
    }

type processedAudio struct {
    sync.Mutex
    recordedSamples []SAMPLE
    buffer64 []float64 //same as recordedSamples just 64 bit
    spektralDensity []float64
    freqs []float64
}


func NewSoundSingnal() (*SoundSingnal) {
    var err error
    d := &SoundSingnal{}
    d.buffer = make([]SAMPLE, FRAMES_PER_BUFFER)
    d.bufferChannel = make(chan []SAMPLE)

    // h, err := portaudio.DefaultHostApi()
    // CheckErr(err)
    // p := portaudio.LowLatencyParameters(h.DefaultInputDevice, nil)
    // d.Stream, err = portaudio.OpenStream(p, d.RecordCallback)
    d.Stream, err = portaudio.OpenDefaultStream(1, 0, 44100, FRAMES_PER_BUFFER, d.RecordCallback)
    CheckErr(err)
    //fmt.Println("Input Device is:", p.Input.Device)
    // d.buffer64 = make([]float64, FRAMES_PER_BUFFER)
    // d.spektralDensity = make([]float64, FRAMES_PER_BUFFER / 2 + 1)
    // d.freqs = make([]float64, FRAMES_PER_BUFFER / 2 + 1)
    return d
}

func NewProcessedAudio() (*processedAudio) {
    d := &processedAudio{}
    d.buffer64 = make([]float64, FRAMES_PER_BUFFER)
    d.recordedSamples = make([]SAMPLE, FRAMES_PER_BUFFER)
    d.spektralDensity = make([]float64, FRAMES_PER_BUFFER / 2 + 1)
    d.freqs = make([]float64, FRAMES_PER_BUFFER / 2 + 1)
    return d
}

func StartSoundTracking() (*SoundSingnal){
    //var streamParameter portaudio.StreamParameters

    portaudio.Initialize()
    s := NewSoundSingnal()

    // streamParameter.Input.Device, err = portaudio.DefaultInputDevice()
    // CheckErr(err)
    // streamParameter.Input.Channels = NUM_CHANNELS //mono
    // streamParameter.Input.Latency = streamParameter.Input.Device.DefaultLowInputLatency //not necessary?
    // streamParameter.SampleRate = SAMPLE_RATE
    // streamParameter.FramesPerBuffer = FRAMES_PER_BUFFER
    //
    //
    // stream, err := portaudio.OpenStream(streamParameter, data.RecordCallback)
    // CheckErr(err)
    // err = stream.Start()
    CheckErr(s.Start())
    fmt.Println("Now recording")
    //time.Sleep(500 * time.Millisecond) //I dont know why this would be neccessary but it somehow is...
    return s
}


func (pa *SoundSingnal) RecordCallback(buffer []SAMPLE) {
    //pa.Lock()
    //pwelchOptions := spectral.PwelchOptions{NFFT: FRAMES_PER_BUFFER}
    if(true) {println("RecordCallback - buffer is", buffer)}
    if(DEBUG) {println("RecordCallback - wants to write to channel")}
    pa.bufferChannel <- buffer
    if(DEBUG) {println("RecordCallback - wrote to channel:", pa.buffer)}
    //pa.Unlock()
    //for i := 0; i < len(buffer) - 1; i++ {
    //    pa.buffer64[i] = float64(buffer[i])
    //}
    //pa.spektralDensity, pa.freqs = spectral.Pwelch(pa.buffer64, SAMPLE_RATE, &pwelchOptions)
}

func (audio *processedAudio) processAudio(data *SoundSingnal){
    go func() {
        for {
            audio.Lock()
            if(DEBUG) {println("processAudio - waiting for data from channel")}
            audio.recordedSamples = <- data.bufferChannel
            if(true) {println("processAudio - received data:", audio.recordedSamples)}
            pwelchOptions := spectral.PwelchOptions{NFFT: FRAMES_PER_BUFFER}
            for i := 0; i < len(audio.recordedSamples) - 1; i++ {
                audio.buffer64[i] = float64(audio.recordedSamples[i])
            }
            audio.spektralDensity, audio.freqs = spectral.Pwelch(audio.buffer64, SAMPLE_RATE, &pwelchOptions)
            audio.Unlock()
        }
    }()
}


func CheckErr(err error) {
    if err != nil {
        fmt.Println(err)
    }
}
