//go:build exclude
package solution2

// Idea: Since only one oxygen atom is needed, oxygen atoms can become "leaders"
// Each leader will:
//   1. Receive arrival requests  from 2 hydrogen atoms.
//   2. Tell the 2 hydrogen atoms to start bonding.
//   3. The leader itself will also bond.
//   3. Wait for the 2 hydrogen atoms to finish bonding and inform leader to continue.
//   4. Step down from being leader.

type WaterFactory struct {
	oxygenMutex chan struct{}
	precomH chan chan struct{}
}

func NewWaterFactory() *WaterFactory {
	return &WaterFactory{
		oxygenMutex: make(chan struct{}, 1),
		precomH: make(chan chan struct{}),
	}
}


func (wf *WaterFactory) Hydrogen(bond func()) {
	// Pre-commit (step 1)
	c := make(chan struct{})
	wf.precomH <- c

	// Wait for signal from daemon for step 2
	<-c

	// Start bonding
	bond()

	// Post-commit (step 3)
	c <- struct{}{}
}

func (wf *WaterFactory) Oxygen(bond func()) {
	// Become leader
	wf.oxygenMutex <- struct{}{}

	// Step 1
	h1 := <-wf.precomH
	h2 := <-wf.precomH

	// Step 2
	h1 <- struct{}{}
	h2 <- struct{}{}
	bond()

	// Step 3
	<-h1
	<-h2

	// Step 4
	<-wf.oxygenMutex
}

///////////////////////////////////////////////////////////////
//             To clean up or wait for goroutine             //
///////////////////////////////////////////////////////////////

func (wf *WaterFactory) AddHydrogen(bond func()) {
	go wf.Hydrogen(bond)
}

func (wf *WaterFactory) AddOxygen(bond func()) {
	go wf.Oxygen(bond)
}

func (wf *WaterFactory) Shutdown() {
}
