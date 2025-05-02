#ifndef WATER_FACTORY
#define WATER_FACTORY

#include <barrier>
#include <semaphore>

class WaterFactory {
private:
    std::barrier<> barrier;
    std::counting_semaphore<> hSem;
    std::counting_semaphore<> oSem;

public:
    WaterFactory();

    void oxygen(void (*bond)());
    void hydrogen(void (*bond)());

    void add_hydrogen(void (*bond)());
    void add_oxygen(void (*bond)());
    void shutdown();
};

#endif
