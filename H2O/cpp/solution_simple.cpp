#include "solution_simple.h"

WaterFactory::WaterFactory() : barrier{3},  hSem{2}, oSem{1} {}

void WaterFactory::oxygen(void (*bond)()) {
    oSem.acquire();
    barrier.arrive_and_wait();
    bond();
    oSem.release();
}

void WaterFactory::hydrogen(void (*bond)()) {
    hSem.acquire();
    barrier.arrive_and_wait();
    bond();
    hSem.release();
}

void WaterFactory::add_oxygen(void (*bond)()) {
    // keep track of threads spawned
    // spawn thread to run oxygen
}

void WaterFactory::add_hydrogen(void (*bond)()) {
    // keep track of threads spawned
    // spawn thread to run hydrogen
}

void WaterFactory::shutdown() {

}
