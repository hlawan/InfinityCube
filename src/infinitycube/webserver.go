package main

import (
    "encoding/json"
    "net/http"
)

type Status struct { 
	LedR [NR_OF_SIDES][EDGES_PER_SIDE][EDGE_LENGTH]int
	LedG [NR_OF_SIDES][EDGES_PER_SIDE][EDGE_LENGTH]int
	LedB [NR_OF_SIDES][EDGES_PER_SIDE][EDGE_LENGTH]int
}

func NewStatus(c *Cube) (s *Status) {
    s = &Status{}  
    for i := 0; i < NR_OF_SIDES; i++ {
    	for o := 0; o < EDGES_PER_SIDE; o++ {
    		for p := 0; p < EDGE_LENGTH; p++ {
    			s.LedR[i][o][p] = int(c.side[i].edge[o].led[p].Red)
				s.LedG[i][o][p] = int(c.side[i].edge[o].led[p].Green)
				s.LedB[i][o][p] = int(c.side[i].edge[o].led[p].Blue)
			}
		}
	}
    return
}

func (c *Cube) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s := NewStatus(c)
    w.Header().Add("Content-Type", "text/json")
    json.NewEncoder(w).Encode(s)
}