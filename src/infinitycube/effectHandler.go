package main

import (
	"reflect"
	"time"
	//"fmt"
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
		effectProperties: make(map[int][]byte)}

	eH.listAvailableEffects()
	return
}

func (eH *EffectHandler) listAvailableEffects() {
	eHType := reflect.TypeOf(eH)
	for i := 0; i < eHType.NumMethod(); i++ {
		method := eHType.Method(i).Name
		if strings.Contains(method, "add") {
			method = strings.TrimPrefix(method, "add")
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

func (eH *EffectHandler) addCellularAutomata(colorOpacity, blackOpacity, secsPerGen float64, rule int) {
	cA := NewCellularAutomata(eH.myDisplay, colorOpacity, blackOpacity, rule, secsPerGen)
	eH.activeEffects = append(eH.activeEffects, cA)
	eH.listEffectProperties() //where is a nice place for me?
}
