#include "queue1.cpp"

#include <cstdio>
#include <thread>

int main() {
    mpmc_queue<int> queue;

    auto producer1 = std::thread([&queue]() {
        for(int i = 1; i <= 150000; i++)
            queue.enqueue(i);
    });

    auto producer2 = std::thread([&queue]() {
        for(int i = 150001; i <= 300000; i++)
            queue.enqueue(i);
    });

    int sum1 = 0;
    int sum2 = 0;
    int sum3 = 0;

    auto consumer_fn = [&queue](int& sum) {
        for (int i = 0; i < 100000; i++) {
            while (true) {
                std::optional<int> n = queue.try_dequeue();
                if (n) {
                    sum += n.value();
                    break;
                }
            }
        }
    };

    auto consumer1 = std::thread(consumer_fn, std::ref(sum1));
    auto consumer2 = std::thread(consumer_fn, std::ref(sum2));
    auto consumer3 = std::thread(consumer_fn, std::ref(sum3));

    producer1.join();
    producer2.join();
    consumer1.join();
    consumer2.join();
    consumer3.join();

    printf("Sum of 1 to 300000 modulo integer limit: %d\n", sum1 + sum2 + sum3);
}
