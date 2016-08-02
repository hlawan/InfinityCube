package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Status gathers all the information that is passed to or received from
// the webserver.
type Status struct {
	io.ReadWriter
	AvailableEffects []string
	ActiveEffects    []string
	effectRequest    chan string
	EffectParameter  map[string]string
	SoundSignal      []SAMPLE
	SpectralDensity  []float64
	Freqs            []float64
	CurrentVolume    float64
	AverageVolume    float64
	MaxPeak          float64
	PeakAverageRatio float64
}

func NewStatus(data *ProcessedAudio, fx []string, ch chan string, h io.ReadWriter) (s *Status) {
	//data.Lock()
	//defer data.Unlock()
	s = &Status{
		ReadWriter:    h,
		effectRequest: ch}
	s.AvailableEffects = make([]string, len(fx))
	s.AvailableEffects = fx
	s.ActiveEffects = make([]string, 0)
	s.SoundSignal = data.recordedSamples
	fmt.Println(len(data.spektralDensity), len(data.freqs))
	s.SpectralDensity = make([]float64, len(data.spektralDensity))
	s.Freqs = make([]float64, len(data.freqs))
	return s
}

func (s *Status) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	effRequest := req.FormValue("t")
	if effRequest != "" {
		fmt.Println("Effect Request: ", effRequest)
		go s.requestEffect("Add" + effRequest)
	}

	parRequest := req.FormValue("act")
	if parRequest != "" {
		fmt.Println("Parameter Request: ", parRequest)
		go s.requestParameter(parRequest)
	}

	delRequest := req.FormValue("r")
	if delRequest != "" {
		fmt.Println("Delete Request: ", delRequest)
		go s.requestEffect("del" + delRequest)
	}

	w.Header().Add("Content-Type", "text/json")
	json.NewEncoder(w).Encode(s)
}

func (s *Status) requestParameter(eff string) {
	if strings.HasPrefix(eff, "set") {
		s.effectRequest <- (eff)
		// ****** to be conituned ********

	} else {
		s.effectRequest <- ("par" + eff)
		fmt.Println("<- sent parameterRequest for", eff)

		// receive current Parameter
		var newPar map[string]string
		newPar = make(map[string]string)
		next := true
		var key, value string
		for Par := range s.effectRequest {
			if Par == "done" {
				//fmt.Println("Received done (requestParameter() webserver)")
				break
			}
			if next {
				key = Par
				//fmt.Println("-> received key = ", key)
				next = false
			} else {
				value = Par
				//fmt.Println("-> received value = ", value)
				newPar[key] = value
				next = true
			}
		}
		//fmt.Println("Effect Parameter: ", newPar)
		s.EffectParameter = newPar
	}
}

func (s *Status) requestEffect(eff string) {
	//fmt.Println("trying to send effect request: ", eff)
	s.effectRequest <- eff
	fmt.Println("Sent Effect Request from webserver: ", eff)
	newActEff := make([]string, 0)
	for actEff := range s.effectRequest {
		if actEff == "done" {
			break
		}
		newActEff = append(newActEff, actEff)
	}
	fmt.Println("Active Effects: ", newActEff)
	s.ActiveEffects = newActEff
}

func (s *Status) UpdateStatus(data *ProcessedAudio) {
	for {
		//data.Lock()
		s.SoundSignal = data.recordedSamples
		for i := 0; i < len(data.freqs); i++ {
			s.SpectralDensity[i] = data.spektralDensity[i]
			s.Freqs[i] = data.freqs[i]
		}
		s.CurrentVolume = data.currentVolume
		s.AverageVolume = data.averageVolume
		s.MaxPeak = data.maxPeak
		s.PeakAverageRatio = data.peakAverageRatio
		//data.Unlock()
		time.Sleep(23 * time.Millisecond)
	}
}

func StartWebServer(data *ProcessedAudio, fx []string, ch chan string) (s *Status) {
	var h io.ReadWriter
	var err error
	h, err = os.OpenFile(*serial_port, os.O_RDWR, 0)
	CheckErr(err)
	s = NewStatus(data, fx, ch, h)
	go s.UpdateStatus(data)

	http.Handle("/toggle", s)
	http.Handle("/status", s)
	http.Handle("/", http.FileServer(http.Dir(*static_path)))
	go func() {
		err = http.ListenAndServe(*listen_address, nil)
		CheckErr(err)
	}()
	return s
}
