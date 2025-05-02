package main

import (
	"fmt"
	"math/rand/v2"
	"time"
	
	//sol "h2o/solution1"
	sol "h2o/solution2"
)

func main() {
	oxygenBond := func() {
		fmt.Println("Bonding oxygen")
		time.Sleep(5 * time.Millisecond)
	}
	hydrogenBond := func() {
		fmt.Println("Bonding hydrogen")
		time.Sleep(5 * time.Millisecond)
	}
	wf := sol.NewWaterFactory()
	o := 0
	h := 0
	for i := 0; i < 1000; i += 1 {
		if rand.IntN(3) == 2 {
			o += 1
			wf.AddOxygen(oxygenBond)
		} else {
			h += 1
			wf.AddHydrogen(hydrogenBond)
		}
	}
	fmt.Printf("h: %d, o: %d\n", h, o)
	wf.Shutdown()
}
