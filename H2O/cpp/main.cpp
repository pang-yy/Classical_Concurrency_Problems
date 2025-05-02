#include <chrono>
#include <cstdlib>
#include <ctime>
#include <iostream>
#include <ostream>
#include <thread>

#include "solution.h"

inline void oxygen_bond() {
    std::cout << "Bonding oxygen" << std::endl;
    std::this_thread::sleep_for(std::chrono::milliseconds{5});
}

inline void hydrogen_bond() {
    std::cout << "Bonding hydrogen" << std::endl;
    std::this_thread::sleep_for(std::chrono::milliseconds{5});
}

int main() {
    WaterFactory wf{};
    int o = 0;
    int h = 0;
    std::srand(std::time(0));
    for (int i = 0; i < 1000; i += 1) {
        if (std::rand() % 3 == 2) {
            o += 1;
            wf.add_oxygen(oxygen_bond);
        } else {
            h += 1;
            wf.add_hydrogen(hydrogen_bond);
        }
    }
    std::cout << "h: " << h << ", o: " << o << std::endl;
    wf.shutdown();
}
