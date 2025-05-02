#include <mutex>
#include <optional>
#include <queue>
#include <semaphore>

/*
 * Coarse-grained concurrent queue implemented using semaphore.
 */
template<typename T>
class mpmc_queue {
private:
    std::queue<T> queue;
    mutable std::mutex mut;
    std::counting_semaphore<> sem;

public:
    mpmc_queue() : queue{}, mut{}, sem{0} {}
    mpmc_queue(const mpmc_queue& other) {
        std::unique_lock<std::mutex> lock{other.mut};
        queue = other.queue;
    }
    mpmc_queue& operator=(const mpmc_queue&) = delete;

    void enqueue(T t) {
        std::unique_lock<std::mutex> lock{mut};
        queue.push(std::move(t));
        sem.release();
    }

    std::optional<T> try_dequeue() {
        if (!sem.try_acquire()) {
            return std::nullopt;
        }
        std::unique_lock<std::mutex> lock{mut};
        T t = std::move(queue.front());
        queue.pop();
        return std::optional<T>{t};
    }

    T dequeue() {
        sem.acquire();
        std::unique_lock<std::mutex> lock{mut};

        T t = std::move(queue.front());
        queue.pop();
        return t;
    }

    bool empty() {
        std::unique_lock<std::mutex> lock{mut};
        return queue.empty();
    }
};
