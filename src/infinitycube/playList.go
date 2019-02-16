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
	fadeLength time.Duration
	fading     bool
}

func NewPlayList(name string, EffectList []map[Effector]time.Duration) *PlayList {
	pl := &PlayList{
		Name:       name,
		Effects:    EffectList,
		slotTimes:  make([]time.Duration, len(EffectList)),
		activeSlot: -1,
		totalTime:  0 * time.Second,
		lastChange: time.Now(),
		fadeLength: 2 * time.Second,
		fading:     false,
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

func (pl *PlayList) fadingPhase() bool {

	moreThaOneEffect := (len(pl.Effects) > 1)
	notFirstRun := (pl.activeSlot >= 0)

	if moreThaOneEffect {
		if notFirstRun {
			// if inside fading window
			if time.Since(pl.lastChange) > (pl.slotTimes[pl.activeSlot]-pl.fadeLength) &&
				time.Since(pl.lastChange) < pl.slotTimes[pl.activeSlot] {
				return true
			}
		}
	}

	// else
	return false
}

func (pl *PlayList) dimmingFaktors() (float64, float64) {

	tmax := pl.slotTimes[pl.activeSlot].Nanoseconds()
	t := time.Since(pl.lastChange).Nanoseconds()

	// start [0 ... 1] end
	fadeProgress := float64(tmax-t) / float64(pl.fadeLength.Nanoseconds())

	// increasing values for rising, decreasing for falling
	falling := fadeProgress
	rising := 1 - fadeProgress

	return rising, falling
}

func (pl *PlayList) SlotEffects() []Effector {

	if pl.fadingPhase() {
		rising, falling := pl.dimmingFaktors()

		next := (pl.activeSlot + 1) % len(pl.Effects)

		slotEffects := make([]Effector, len(pl.Effects[pl.activeSlot])+len(pl.Effects[next]))

		// Add effects from current slot with decreasing faktor
		//iterate of map (key, value) key = type Effect
		i := 0
		for effect, _ := range pl.Effects[pl.activeSlot] {
			slotEffects[i] = effect
			slotEffects[i].SetDimming(falling)

			i++
		}
		// Add effects from next slot with incresing dimming faktor
		for effect, _ := range pl.Effects[next] {
			slotEffects[i] = effect
			slotEffects[i].SetDimming(rising)

			i++
		}

		return slotEffects

	} else if pl.slotChange() {

		pl.fading = false
		pl.activeSlot = (pl.activeSlot + 1) % len(pl.Effects)

		slotEffects := make([]Effector, len(pl.Effects[pl.activeSlot]))

		// default value
		currentDimming := 1.0

		// Add effects from current slot
		//iterate of map (key, value) key = type Effect
		i := 0
		for effect, _ := range pl.Effects[pl.activeSlot] {
			slotEffects[i] = effect
			slotEffects[i].SetDimming(currentDimming)
			i++
		}

		pl.lastChange = time.Now()

		return slotEffects
	} else {
		return nil
	}

}
