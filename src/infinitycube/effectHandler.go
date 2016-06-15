package main

import (
	"time"
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
	Effects    []Effector
	myDisplay  Display
	lastUpdate time.Time
	updateRate int
	loopTime   time.Duration
}

func NewEffectHandler(newDisplay Display, newUpdateRate int) (eH *EffectHandler) {
	eH = &EffectHandler{
		myDisplay:  newDisplay,
		lastUpdate: time.Now(),
		updateRate: newUpdateRate,
		loopTime:   10 * time.Millisecond}
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
	for _, effect := range eH.Effects {
		effect.Update()
	}
}

func (eH *EffectHandler) addCellularAutomata(colorOpacity, blackOpacity, secsPerGen float64, rule int) {
	cA := NewCellularAutomata(eH.myDisplay, colorOpacity, blackOpacity, rule, secsPerGen)
	eH.Effects = append(eH.Effects, cA)
}
