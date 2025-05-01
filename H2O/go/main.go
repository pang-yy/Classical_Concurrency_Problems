package main

import (
	"fmt"
	"math/rand"
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
	rand.Seed(time.Now().Unix())
	oxy := 10
	//hydro := oxy * 2
	for i := 0; i < oxy; i += 1 {
	/*
		if rand.Intn(3) == 2 {
			wf.AddOxygen(oxygenBond)
		} else {
			wf.AddHydrogen(hydrogenBond)
		}
	*/
		wf.AddHydrogen(hydrogenBond)
		wf.AddOxygen(oxygenBond)
		wf.AddHydrogen(hydrogenBond)
	}
	wf.Shutdown()
}
