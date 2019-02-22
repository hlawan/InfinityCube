package main

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/gordonklaus/portaudio"
	"github.com/mjibson/go-dsp/spectral"
)

const (
	SAMPLE_RATE       = 20000
	FRAMES_PER_BUFFER = 512
	NUM_CHANNELS      = 1
	DEBUG             = false
)

type SAMPLE float32

type SoundSingnal struct {
	sync.Mutex
	*portaudio.Stream
	buffer        []SAMPLE
	bufferChannel chan []SAMPLE
}

type ProcessedAudio struct {
	sync.Mutex
	sampleRate              int
	recordedSamples         []SAMPLE
	buffer64                []float64 //same as recordedSamples just 64 bit
	spektralDensity         []float64
	freqs                   []float64
	currentVolume           float64
	averageVolume           float64
	maxPeak                 float64
	peakAverageRatio        float64
	averagePeakAverageValue float64
	loudestFreq             float64
	loudestFreqAmpl         float64
	timeOfLastClap          time.Time
	clapDetected            bool
	clapCount               int
}

func NewSoundSingnal() *SoundSingnal {
	var err error
	d := &SoundSingnal{}
	d.buffer = make([]SAMPLE, FRAMES_PER_BUFFER)
	d.bufferChannel = make(chan []SAMPLE)
	d.Stream, err = portaudio.OpenDefaultStream(1, 0, SAMPLE_RATE, FRAMES_PER_BUFFER, d.RecordCallback)
	CheckErr(err)
	return d
}

func NewProcessedAudio() *ProcessedAudio {
	d := &ProcessedAudio{}
	d.buffer64 = make([]float64, FRAMES_PER_BUFFER)
	d.recordedSamples = make([]SAMPLE, FRAMES_PER_BUFFER)
	d.spektralDensity = make([]float64, FRAMES_PER_BUFFER/2+1)
	d.freqs = make([]float64, FRAMES_PER_BUFFER/2+1)
	d.sampleRate = SAMPLE_RATE
	return d
}

func StartSoundTracking() *SoundSingnal {

	portaudio.Initialize()
	s := NewSoundSingnal()

	// check if audio input stream is available on the system
	if s.Stream != nil {
		err := s.Start()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Audio input found")
		}
	} else {
		if debugLvl >= 1 {
			fmt.Println("No audio input available. Audio-reactive Effects will not work!")
		}
	}

	return s
}

func (pa *SoundSingnal) RecordCallback(buffer []SAMPLE) {
	pa.bufferChannel <- buffer
}

func StartAudioProcessing(data *SoundSingnal) *ProcessedAudio {
	audio := NewProcessedAudio()
	go func() {
		for {
			audio.Lock() //pretty much blocks audio the whole time
			audio.recordedSamples = <-data.bufferChannel
			for i := 0; i < len(audio.recordedSamples)-1; i++ {
				audio.buffer64[i] = float64(audio.recordedSamples[i])
			}
			audio.anlayseSpectrum()
			audio.getVolume()
			audio.detectClap()
			audio.Unlock()
		}
	}()
	return audio
}

func (audio *ProcessedAudio) anlayseSpectrum() {
	pwelchOptions := spectral.PwelchOptions{NFFT: FRAMES_PER_BUFFER}
	audio.spektralDensity, audio.freqs = spectral.Pwelch(audio.buffer64, SAMPLE_RATE, &pwelchOptions)
	amplitude := 0.0
	index := 0
	for i := 0; i < len(audio.spektralDensity); i++ {
		audio.spektralDensity[i] *= 2000 //make the values nicer....
		if audio.spektralDensity[i] > amplitude {
			amplitude = audio.spektralDensity[i]
			index = i
		}
	}
	audio.loudestFreq = audio.freqs[index]
	audio.loudestFreqAmpl = amplitude
}

const (
	VolumeWeight       = .5                                                                 //the higher, the faster the averageVolume changes
	scaledVolumeWeight = VolumeWeight * (float64(FRAMES_PER_BUFFER) / float64(SAMPLE_RATE)) //makes averageCalculation independent from BufferSize
)

func (audio *ProcessedAudio) getVolume() {
	newAverageVolume := 0.0
	maxPeak := audio.buffer64[0]
	for i := 0; i < FRAMES_PER_BUFFER; i++ {
		newAverageVolume += math.Abs(audio.buffer64[i])
		if math.Abs(audio.buffer64[i]) > maxPeak {
			maxPeak = math.Abs(audio.buffer64[i])
		}
	}
	newAverageVolume /= FRAMES_PER_BUFFER
	audio.currentVolume = newAverageVolume
	audio.averageVolume = audio.averageVolume*(1-scaledVolumeWeight) + newAverageVolume*scaledVolumeWeight
	audio.maxPeak = maxPeak
	if audio.averageVolume > 0 {
		audio.peakAverageRatio = maxPeak / audio.averageVolume
	}
	audio.averagePeakAverageValue = audio.averagePeakAverageValue*(1-scaledVolumeWeight) + audio.peakAverageRatio*scaledVolumeWeight
}

func (audio *ProcessedAudio) detectClap() {
	if audio.peakAverageRatio > 15 && time.Since(audio.timeOfLastClap) > 1000*time.Millisecond {
		audio.timeOfLastClap = time.Now()
		audio.clapDetected = true
		audio.clapCount += 1
		//fmt.Println("Clap detected average PeakAverageRatio is:", audio.averagePeakAverageValue, "\tclapCount:", audio.clapCount)
	} else {
		audio.clapDetected = false
	}
}

func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
