// playList
package main

import (
	//"fmt"
	"time"
)

type PlayList struct {
	Name       string
	Effects    []map[Effector]time.Duration
	slotTimes  []time.Duration
	activeSlot int
	totalTime  time.Duration
	lastChange time.Time
}

func NewPlayList(name string, EffectList []map[Effector]time.Duration) *PlayList {
	pl := &PlayList{
		Name:       name,
		Effects:    EffectList,
		slotTimes:  make([]time.Duration, len(EffectList)),
		activeSlot: 0,
		totalTime:  0 * time.Second,
		lastChange: time.Now().Add(-24 * time.Hour),
	}

	// find max effect duration for each slot an calculate total runtime of playList
	for i, slot := range pl.Effects {

		tMax := 0 * time.Second
		for _, effectDuration := range slot {

			if effectDuration > tMax {
				tMax = effectDuration
			}

		}
		pl.slotTimes[i] = tMax
		pl.totalTime += tMax
	}

	return pl
}

func (pl *PlayList) selectSlot() int {

	if time.Since(pl.lastChange) > pl.slotTimes[pl.activeSlot] {
		return pl.activeSlot + 1
	} else {
		return pl.activeSlot
	}

}

func (pl *PlayList) SlotEffects() []Effector {

	if pl.activeSlot < pl.selectSlot() {

		slotEffects := make([]Effector, len(pl.Effects[pl.activeSlot]))

		i := 0
		for effect, _ := range pl.Effects[pl.activeSlot] {
			slotEffects[i] = effect
			i++
		}

		pl.activeSlot = (pl.activeSlot + 1) % len(pl.slotTimes)
		pl.lastChange = time.Now()

		return slotEffects
	} else {
		return nil
	}

}
