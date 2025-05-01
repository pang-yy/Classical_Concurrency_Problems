//go:build exclude
package solution1

// Idea: Create a "daemon" goroutine to coordinate everything.
// The daemon will:
//   1. Receive arrival requests from 2 hydrogen atoms and 1 oxygen atoms.
//   2. Tell all 3 atoms to proceed and begin bonding.
//   3. Wait for all 3 atoms to finish bonding and inform daemon to continue.
//   4. Go to step 1.

type WaterFactory struct {
	precomH chan chan struct{}
	precomO chan chan struct{}
}

func NewWaterFactory() *WaterFactory {
	wf := &WaterFactory{
		precomH: make(chan chan struct{}),
		precomO: make(chan chan struct{}),
	}

	go initDaemon(wf)

	return wf
}

func (wf *WaterFactory) Hydrogen(bond func()) {
	//fmt.Println("New Hydrogen")

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
	//fmt.Println("New Oxygen")
	
	// Pre-commit (step 1)
	c := make(chan struct{})
	wf.precomO <- c

	// Wait for signal from daemon for step 2
	<-c

	// Start bonding
	bond()

	// Post-commit (step 3)
	c <- struct{}{}
}

func initDaemon(wf *WaterFactory) {
	for {
		// Step 1
		h1 := <-wf.precomH
		h2 := <-wf.precomH
		o  := <-wf.precomO

		// Step 2
		h1 <- struct{}{}
		h2 <- struct{}{}
		o  <- struct{}{}

		// Step 3
		<-h1
		<-h2
		<-o
	}
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
