package main

import (
	"reflect"
	"time"
	"fmt"
	"encoding/json"
	"strings"
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
	myDisplay		 Display
}

func NewEffect(disp Display, colorOp, blackOp float64) *Effect {
	ef := &Effect{
		Leds: 				make([]Led, disp.NrOfLeds()),
		Offset:				0,
		Length:				disp.NrOfLeds(),
		ColorOpacity: colorOp,
		BlackOpacity: blackOp,
		myDisplay:		disp}

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
	effectProperties map[int][]byte
	availableEffects []string
	effectRequest		 chan string
	myDisplay        Display
	lastUpdate       time.Time
	updateRate       int
	loopTime         time.Duration
	audio						 *ProcessedAudio
}

func NewEffectHandler(newDisplay Display, newUpdateRate int, newAudio *ProcessedAudio) (eH *EffectHandler) {
	eH = &EffectHandler{
		myDisplay:        newDisplay,
		lastUpdate:       time.Now(),
		updateRate:       newUpdateRate,
		loopTime:         10 * time.Millisecond,
		effectProperties: make(map[int][]byte),
		effectRequest: 		make(chan string),
		audio:						newAudio}

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

func (eH *EffectHandler) handleRequests() (err error){
  eHValue := reflect.ValueOf(eH)
	for{
		//receive effect request from webserver
		req := <- eH.effectRequest
		//add requested effect to activeEffects
		req = "Add" + req
  	m := eHValue.MethodByName(req)
  	if !m.IsValid() {
      return fmt.Errorf("Method not found \"%s\"", req)
  	}
  	in := make([]reflect.Value, 0)
  	m.Call(in)
		//send updated list of active effects to webserver
		for _, effect := range eH.activeEffects {
			eType := reflect.TypeOf(effect)
			eH.effectRequest <- eType.Elem().Name()
		}
		eH.effectRequest <- "done"
	}
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

func (eH *EffectHandler) listEffectProperties() {
	var err error
	for i, ele := range eH.activeEffects {
		eH.effectProperties[i], err = json.Marshal(ele)
		CheckErr(err)
	}
}

// List of addable Effects

func (eH *EffectHandler) AddCellularAutomata() {
	ca := NewCellularAutomata(eH.myDisplay)
	eH.activeEffects = append(eH.activeEffects, ca)
	eH.listEffectProperties() //where is a nice place for me?
}

func (eH *EffectHandler) AddWhiteSpectrum() {
	ws := NewWhiteSpectrum(eH.myDisplay, eH.audio)
	eH.activeEffects = append(eH.activeEffects, ws)
	eH.listEffectProperties() //where is a nice place for me?
}

func (eH *EffectHandler) AddWhiteEdgeSpectrum() {
	wes := NewWhiteEdgeSpectrum(eH.myDisplay, eH.audio)
	eH.activeEffects = append(eH.activeEffects, wes)
	eH.listEffectProperties() //where is a nice place for me?
}

func (eH *EffectHandler) AddEdgeVolume() {
	ev := NewEdgeVolume(eH.myDisplay, eH.audio)
	eH.activeEffects = append(eH.activeEffects, ev)
	eH.listEffectProperties() //where is a nice place for me?
}


func (eH *EffectHandler) AddRunningLight(){
	fmt.Println("look at me I'm so pretty :-*")
}
