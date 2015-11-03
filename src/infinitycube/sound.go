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
    d.Stream, err = portaudio.OpenDefaultStream(1, 0, 44100, FRAMES_PER_BUFFER, d.RecordCallback)
    CheckErr(err)
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
    portaudio.Initialize()
    s := NewSoundSingnal()
    CheckErr(s.Start())
    fmt.Println("Now recording")
    return s
}


func (pa *SoundSingnal) RecordCallback(buffer []SAMPLE) {
    pa.bufferChannel <- buffer
}

func (audio *processedAudio) processAudio(data *SoundSingnal){
  go func() {
    for {
      audio.Lock()
      audio.recordedSamples = <- data.bufferChannel
      for i := 0; i < len(audio.recordedSamples) - 1; i++ {
        audio.buffer64[i] = float64(audio.recordedSamples[i])
      }
      audio.anlayseSpectrum()
      audio.Unlock()
    }
  }()
}

func (audio *processedAudio) anlayseSpectrum(){
    pwelchOptions := spectral.PwelchOptions{NFFT: FRAMES_PER_BUFFER}
    audio.spektralDensity, audio.freqs = spectral.Pwelch(audio.buffer64, SAMPLE_RATE, &pwelchOptions)
}



func CheckErr(err error) {
    if err != nil {
        fmt.Println(err)
    }
}
