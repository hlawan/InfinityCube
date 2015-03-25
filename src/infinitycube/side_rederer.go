package main

import (
    "fmt"
)

func (my *side) setSide(red, green, blue uint8) {
	if(DEBUG_LVL >=2) {fmt.Println("setSide called with Arguments", red, green, blue, "btw i've got", len(my.edge), "edges")}
	for i := 0; i < EDGES_PER_SIDE; i++ {
		my.edge[i].renderer = my.edge[i].setEdge //y we no work?!
		my.edge[i].renderer(red, green, blue)
	}

}

//--------------new stuff----------------//

