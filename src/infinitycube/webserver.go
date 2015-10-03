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
// *  The color infomation have to be transformed into int arrays because the
// *  coffeescript didn't like my Cube struct.*/
type Status struct {
    io.ReadWriter
    SoundSignal []SAMPLE
}
//
// /*NewStatus parses the color information from the Cube struct to the int array
// *  and collects all the other information.*/
func NewStatus(data *paTestData, h io.ReadWriter) (s *Status) {
	s = &Status{ReadWriter: h}
    s.SoundSignal = data.recordedSamples
    return s
}
// 	//idea: automatic gathering of all known methodes of Cube type in a string (reflect?)
// 	s.CubeRenderer = []string{"RGBiteration", "simpleRunningLight", "and so on"}
// 	return
// }
//
func (s *Status) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/json")
	json.NewEncoder(w).Encode(s)
}

func (s *Status) UpdateStatus(data *paTestData){
    for{
        s.SoundSignal = data.recordedSamples
        time.Sleep(100 * time.Millisecond)
    }
}

func StartWebServer(data *paTestData) {
    var h io.ReadWriter
    var err error
    s := NewStatus(data, h)
    go s.UpdateStatus(data)
    h, err = os.OpenFile(*serial_port, os.O_RDWR, 0)
    if err != nil {
        fmt.Print("there seems to be an error: ", err)
        return
    }
    http.Handle("/status", s)
    http.Handle("/", http.FileServer(http.Dir(*static_path)))

    err = http.ListenAndServe(*listen_address, nil)
    if err != nil {
        fmt.Print(err)
    }
}
