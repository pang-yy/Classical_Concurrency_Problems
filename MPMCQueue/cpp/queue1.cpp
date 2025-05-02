#include <condition_variable>
#include <mutex>
#include <optional>
#include <queue>

/*
 * Coarse-grained concurrent queue implemented using monitor.
 */
template<typename T>
class mpmc_queue {
private:
    std::queue<T> queue;
    mutable std::mutex mut;
    std::condition_variable cond;

public:
    mpmc_queue() : queue{}, mut{}, cond{} {}
    mpmc_queue(const mpmc_queue& other) {
        std::unique_lock<std::mutex> lock{other.mut};
        queue = other.queue;
    }
    mpmc_queue& operator=(const mpmc_queue&) = delete;

    void enqueue(T t) {
        std::unique_lock<std::mutex> lock{mut};
        queue.push(std::move(t));
        cond.notify_one();
    }

    std::optional<T> try_dequeue() {
        std::unique_lock<std::mutex> lock{mut};
        if (queue.empty()) {
            return std::nullopt;
        }
        T t = std::move(queue.front());
        queue.pop();
        return std::optional<T>{t};
    }

    T dequeue() {
        std::unique_lock<std::mutex> lock{mut};

        // Prevent spurious wakeup
        cond.wait(lock, [this](){ return !queue.empty(); });
        // or this way
        /*
        while (queue.empty()) {
            cond.wait(lock);
        }
        */

        T t = std::move(queue.front());
        queue.pop();
        return t;
    }

    bool empty() {
        std::unique_lock<std::mutex> lock{mut};
        return queue.empty();
    }
};
