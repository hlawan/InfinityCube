package main

import (
	//"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	//"github.com/fatih/structs"
)

// An Effect has to provide information about the current led-pattern,
// the blackOpacity and the colorOpacity of its leds and which Display
// it belongs to.
type Effect struct {
	Leds         []Led
	Offset       int
	Length       int
	ColorOpacity float64
	BlackOpacity float64
	myDisplay    Display
}

func NewEffect(disp Display, colorOp, blackOp float64) *Effect {
	ef := &Effect{
		Leds:         make([]Led, disp.NrOfLeds()),
		Offset:       0,
		Length:       disp.NrOfLeds(),
		ColorOpacity: colorOp,
		BlackOpacity: blackOp,
		myDisplay:    disp}

	return ef
}

// An Effector is able to generate (light-)Patterns. The Update() method gets
// called by its EffectHandler and has to call the AddPattern() method of its
// Display.
type Effector interface {
	Update()
}

// A Display is able to receive, merge and show the (light-)Patterns generated
// by various Effectors.
type Display interface {
	NrOfLeds() int
	AddPattern([]Led, float64, float64)
	Show()
}

// An EffectHandler creates, deletes, configures and updates Effectors. After
// updating all active Effectors the EffectHandler calls the show() method of
// its Display. The updateRate defines the frames per second.
type EffectHandler struct {
	activeEffects    []Effector
	effectProperties []map[string]string
	availableEffects []string
	effectRequest    chan string
	myDisplay        Display
	lastUpdate       time.Time
	updateRate       int
	loopTime         time.Duration
	audio            *ProcessedAudio
}

func NewEffectHandler(newDisplay Display, newUpdateRate int, newAudio *ProcessedAudio) (eH *EffectHandler) {
	eH = &EffectHandler{
		myDisplay:  newDisplay,
		lastUpdate: time.Now(),
		updateRate: newUpdateRate,
		loopTime:   10 * time.Millisecond,
		//effectProperties: make([]map[string][]string),
		effectRequest: make(chan string),
		audio:         newAudio}

	eH.listAvailableEffects()
	go eH.handleRequests()
	return
}

func (eH *EffectHandler) render() {
	loopStart := time.Now()

	eH.updateAll()
	eH.myDisplay.Show()

	if time.Since(loopStart) < (eH.loopTime) {
		time.Sleep(eH.loopTime - time.Since(loopStart))
	}
}

func (eH *EffectHandler) updateAll() {
	for _, effect := range eH.activeEffects {
		effect.Update()
	}
}

func (eH *EffectHandler) handleRequests() (err error) {
	eHValue := reflect.ValueOf(eH)
	var idx int
	for {
		//receive effect request from webserver
		fmt.Println("Waiting for Request (Backend)...")
		req := <-eH.effectRequest

		if strings.HasPrefix(req, "par") {
			// it is a parameter request
			req = strings.TrimPrefix(req, "par")
			idx, err = strconv.Atoi(req)
			CheckErr(err)
			eH.sendEffectProperties(idx)
		} else {
			if strings.HasPrefix(req, "del") {
				// it is a delete request
				req = strings.TrimPrefix(req, "del")
				idx, err = strconv.Atoi(req)
				CheckErr(err)
				eH.removeEffect(idx)
			} else {
				//request to add effect to activeEffects
				req = "Add" + req
				m := eHValue.MethodByName(req)
				if !m.IsValid() {
					return fmt.Errorf("Method not found \"%s\"", req)
				}
				in := make([]reflect.Value, 0)
				m.Call(in)
			}

			//send updated list of active effects to webserver
			for _, effect := range eH.activeEffects {
				eType := reflect.TypeOf(effect)
				eH.effectRequest <- eType.Elem().Name()
			}
			eH.effectRequest <- "done"
		}
	}
}

func (eH *EffectHandler) removeEffect(ele int) {
	var newList []Effector
	for i, effect := range eH.activeEffects {
		if i != ele {
			newList = append(newList, effect)
		}
	}
	eH.activeEffects = newList
}

func (eH *EffectHandler) listAvailableEffects() {
	eHType := reflect.TypeOf(eH)
	for i := 0; i < eHType.NumMethod(); i++ {
		method := eHType.Method(i).Name
		if strings.Contains(method, "Add") {
			method = strings.TrimPrefix(method, "Add")
			eH.availableEffects = append(eH.availableEffects, method)
		}
	}
}

func (eH *EffectHandler) sendEffectProperties(nr int) {
  eH.listEffectProperties()
	fmt.Println("Nr of Effects with Properties = ", len(eH.effectProperties),
							" Trying to send props of nr: ",nr+1)
	for par, val := range eH.effectProperties[nr] {
		eH.effectRequest <- par
		fmt.Println(" sent: ", par, " as par")
		eH.effectRequest <- val
		fmt.Println(" sent: ", val, " as val")
	}
	eH.effectRequest <- "done"
	fmt.Println(" sent: done (EffectProperties)")
}

func (eH *EffectHandler) listEffectProperties() {
	var propList []map[string]string
	//var props = make(map[string][]string))

	for i, ele := range eH.activeEffects {
		propList = append(propList, make(map[string]string))

		s := reflect.ValueOf(ele).Elem()
		typ := s.Type()

		for p := 0; p < s.NumField(); p++ {
			prop := s.Field(p)
			id := typ.Field(p).Name
			if strings.HasSuffix(id, "Par") {
				id = strings.TrimSuffix(id, "Par")
				if prop.Kind() == reflect.Int {
					propList[i][id] = strconv.Itoa(int(prop.Int()))
				} else if prop.Kind() == reflect.Float64 {
					propList[i][id] = FloatToString(prop.Float())
				} else {
					propList[i][id] = " "
				}
			}
		}

		eH.effectProperties = propList;
	}

	// fmt.Println("*** Listed EffectProperties")
	// for i, props := range propList {
	// 	fmt.Println("* Effect Nr ", i)
	// 	fmt.Println(props)
	// }

}

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

// List of addable Effects

func (eH *EffectHandler) AddCellularAutomata() {
	ca := NewCellularAutomata(eH.myDisplay)
	eH.activeEffects = append(eH.activeEffects, ca)
}

func (eH *EffectHandler) AddWhiteSpectrum() {
	ws := NewWhiteSpectrum(eH.myDisplay, eH.audio)
	eH.activeEffects = append(eH.activeEffects, ws)
}

func (eH *EffectHandler) AddWhiteEdgeSpectrum() {
	wes := NewWhiteEdgeSpectrum(eH.myDisplay, eH.audio)
	eH.activeEffects = append(eH.activeEffects, wes)
}

func (eH *EffectHandler) AddEdgeVolume() {
	ev := NewEdgeVolume(eH.myDisplay, eH.audio)
	eH.activeEffects = append(eH.activeEffects, ev)
}

func (eH *EffectHandler) AddHsvFade() {
	hsv := NewHsvFade(eH.myDisplay, eH.updateRate)
	eH.activeEffects = append(eH.activeEffects, hsv)
}

func (eH *EffectHandler) AddRunningLight() {
	rl := NewRunningLight(eH.myDisplay)
	eH.activeEffects = append(eH.activeEffects, rl)
}

func (eH *EffectHandler) AddGausRunningLight() {
	grl := NewGausRunningLight(eH.myDisplay, eH.updateRate)
	eH.activeEffects = append(eH.activeEffects, grl)
}

func (eH *EffectHandler) AddPrettyPrint() {
	fmt.Println("look at me I'm so pretty :-*")
}
