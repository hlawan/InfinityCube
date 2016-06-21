package main

import (
	"reflect"
	"time"
	"fmt"
	"encoding/json"
	"strings"
	//"github.com/fatih/structs"
)

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
}

func NewEffectHandler(newDisplay Display, newUpdateRate int) (eH *EffectHandler) {
	eH = &EffectHandler{
		myDisplay:        newDisplay,
		lastUpdate:       time.Now(),
		updateRate:       newUpdateRate,
		loopTime:         10 * time.Millisecond,
		effectProperties: make(map[int][]byte),
		effectRequest: 		make(chan string)}

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
		req := <- eH.effectRequest
		req = "Add" + req
  	m := eHValue.MethodByName(req)
  	if !m.IsValid() {
      return fmt.Errorf("Method not found \"%s\"", req)
  	}
  	in := make([]reflect.Value, 0)
  	m.Call(in)
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

func (eH *EffectHandler) AddCellularAutomata() {
	cA := NewCellularAutomata(eH.myDisplay)
	eH.activeEffects = append(eH.activeEffects, cA)
	eH.listEffectProperties() //where is a nice place for me?
}

func (eH *EffectHandler) AddRunningLight(){
	fmt.Println("look at me I'm so pretty :-*")
}
