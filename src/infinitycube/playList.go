// playList
package main

import (
	"fmt"
	"time"
)

type PlayList struct {
	Name    string
	Effects []map[Effector]time.Duration
}

func NewPlayList(name string, EffectList []map[Effector]time.Duration) *PlayList {
	pl := &PlayList{
		Name:    name,
		Effects: EffectList,
	}

	// for _, effects := range EffectList {

	// 	TimigMap := map[Effector]time.Duration{}

	// 	for effect, playTime := range effects {
	// 		TimigMap[effect] = playTime
	// 	}

	// 	pl.Effects = append(pl.Effects, TimigMap)
	// }

	fmt.Println("playList")
	fmt.Println(pl)

	return pl
}
