// playList
package main

import (
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
		activeSlot: -1,
		totalTime:  0 * time.Second,
		lastChange: time.Now(),
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

func (pl *PlayList) slotChange() bool {

	// first run
	if pl.activeSlot == -1 {
		return true
	}

	// slot finished
	if time.Since(pl.lastChange) > pl.slotTimes[pl.activeSlot] {
		return true
	} else {
		return false
	}

}

func (pl *PlayList) SlotEffects() []Effector {

	if pl.slotChange() {
		pl.activeSlot = (pl.activeSlot + 1) % len(pl.Effects)

		slotEffects := make([]Effector, len(pl.Effects[pl.activeSlot]))

		i := 0
		for effect, _ := range pl.Effects[pl.activeSlot] {
			slotEffects[i] = effect
			i++
		}

		pl.lastChange = time.Now()
		return slotEffects
	} else {
		return nil
	}

}
