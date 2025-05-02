package solution2

import (
	"fmt"
	"sync"
)

// Idea: Since only one oxygen atom is needed, oxygen atoms can become "leaders"
// Each leader will:
//   1. Receive arrival requests  from 2 hydrogen atoms.
//   2. Tell the 2 hydrogen atoms to start bonding.
//   3. The leader itself will also bond.
//   3. Wait for the 2 hydrogen atoms to finish bonding and inform leader to continue.
//   4. Step down from being leader.
// Requirements:
//   1. The ratio of hydrogen:oxygen is not guaranteed to be exactly 2:1.
//   2. All goroutines/struct need to be cleaned up.

type WaterFactory struct {
	oxygenMutex chan struct{}
	precomH chan chan struct{}

	hCount int
	oCount int
	cond *sync.Cond
	quit chan struct{}
	wg *sync.WaitGroup
}

func NewWaterFactory() *WaterFactory {
	wf := &WaterFactory{
		oxygenMutex: make(chan struct{}, 1),
		precomH: make(chan chan struct{}),

		hCount: 0,
		oCount: 0,
		quit: make(chan struct{}, 10),
		wg: &sync.WaitGroup{},
	}
	wf.cond = sync.NewCond(&sync.Mutex{})

	return wf
}

func (wf *WaterFactory) Hydrogen(bond func()) {
	defer wf.wg.Done()
	c := make(chan struct{})
	for {
		select {
		// Pre-commit (step 1)
		case wf.precomH <- c:
			for {
				select {
				// Wait for signal from daemon for step 2
				case <-c:
					// Start bonding
					bond()

					wf.cond.L.Lock()
					wf.hCount -= 1
					wf.cond.L.Unlock()
					
					// Post-commit (step 3)
					c <- struct{}{}
					return
				case <-wf.quit:
					fmt.Println("Unbound hydrogen")
					return
				}
			}
		case <-wf.quit:
			fmt.Println("Unbound hydrogen")
			return
		}
	}
}

func (wf *WaterFactory) Oxygen(bond func()) {
	defer wf.wg.Done()
	for {
		select {
		case wf.oxygenMutex <- struct{}{}:
			hs := make([]chan struct{}, 2)
			i := 0
			for {
				select {
				case hs[i] = <-wf.precomH:
					i += 1
					if (i == 2) {
						hs[0] <- struct{}{}
						hs[1] <- struct{}{}
						bond()

						<-hs[0]
						<-hs[1]

						wf.cond.L.Lock()
						wf.oCount -= 1
						wf.cond.L.Unlock()
						wf.cond.Signal()

						// Step 4
						<-wf.oxygenMutex
						return
					}
				case <-wf.quit:
					fmt.Println("Unbond oxygen")
					return
				}
			}
		case <-wf.quit:
			fmt.Println("Unbond oxygen")
			return
		}
	}
}

///////////////////////////////////////////////////////////////
//             To clean up or wait for goroutine             //
///////////////////////////////////////////////////////////////

func (wf *WaterFactory) AddHydrogen(bond func()) {
	wf.wg.Add(1)

	wf.cond.L.Lock()
	wf.hCount += 1
	wf.cond.L.Unlock()

	go wf.Hydrogen(bond)
}

func (wf *WaterFactory) AddOxygen(bond func()) {
	wf.wg.Add(1)

	wf.cond.L.Lock()
	wf.oCount += 1
	wf.cond.L.Unlock()

	go wf.Oxygen(bond)
}

func (wf *WaterFactory) Shutdown() {
	wf.cond.L.Lock()
	defer wf.cond.L.Unlock()
	for wf.hCount >= 2 && wf.oCount >= 1 {
		wf.cond.Wait()
	}

	// Signals to leftover Hydrogen and Oxygen goroutines quit.
	close(wf.quit)
	wf.wg.Wait()
}
