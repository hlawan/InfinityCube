package main

import (
	"encoding/json"
	"net/http"
    "os"
    "io"
    "fmt"
    "time"
)

// /*Status gathers all the information that is passed on to the webserver.

type Status struct {
  io.ReadWriter
  SoundSignal []SAMPLE
	SpectralDensity []float64
	Freqs []float64
}

func NewStatus(data *processedAudio, h io.ReadWriter) (s *Status) {
	data.Lock()
	defer data.Unlock()
	s = &Status{ReadWriter: h}
  s.SoundSignal = data.recordedSamples
	fmt.Println(len(data.spektralDensity), len(data.freqs))
	s.SpectralDensity = make([]float64, len(data.spektralDensity))
	s.Freqs = make([]float64, len(data.freqs))
  return s
}

func (s *Status) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/json")
	json.NewEncoder(w).Encode(s)
}

func (s *Status) UpdateStatus(data *processedAudio) {
  for {
		//data.Lock()
  	s.SoundSignal = data.recordedSamples
		for i := 0; i < len(data.freqs); i++ {
			s.SpectralDensity[i] = data.spektralDensity[i]
			s.Freqs[i] = data.freqs[i]
		}
		//data.Unlock()
  	time.Sleep(23 * time.Millisecond)
  }
}

func StartWebServer(data *processedAudio) {
  var h io.ReadWriter
  var err error
	h, err = os.OpenFile(*serial_port, os.O_RDWR, 0)
	CheckErr(err)
  s := NewStatus(data, h)
  go s.UpdateStatus(data)

  http.Handle("/status", s)
  http.Handle("/", http.FileServer(http.Dir(*static_path)))
	go func(){
  	err = http.ListenAndServe(*listen_address, nil)
		CheckErr(err)
	}()
	return
}
