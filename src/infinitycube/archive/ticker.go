type IntervalTicker struct {
	Consumer
	Last     time.Duration
	Interval time.Duration
}

func (i *IntervalTicker) Tick(d time.Duration, o interface{}) {
	fire := false
	if d-i.Last > i.Interval {
		fire = true
		i.Last = d
	}
	i.Consumer.Tick(d, fire)
}

type RandomTicker struct {
	Consumer
	Threshold float32
}

func (r *RandomTicker) Tick(d time.Duration, o interface{}) {
	v := rand.Float32()
	r.Consumer.Tick(d, v < r.Threshold)
}
