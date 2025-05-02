#include <algorithm>
#include <memory>
#include <mutex>
#include <optional>

/*
 * Fine-grained lock based concurrent queue.
 */
template<typename T>
class mpmc_queue {
private:
    struct node {
        T data;
        std::unique_ptr<node> next;

        node() : data{}, next{nullptr} {}
        node(T data_) : data{data_} {}
    };
    std::mutex front_mutex;
    std::unique_ptr<node> front;

    std::mutex back_mutex;
    node* back; // Cannot be unique_ptr because initially front and back will point to same node

    node* get_back() {
        std::unique_lock<std::mutex> lock{back_mutex};
        return back;
    }

    std::unique_ptr<node> dequeue_front() {
        std::unique_lock<std::mutex> lock{front_mutex};
        if (front.get() == get_back()) {
            return nullptr;
        }
        std::unique_ptr<node> old_front = std::move(front);
        front = std::move(old_front->next);
        return old_front;
    }
public:
    // Pre-allocate a dummy node with no data,
    // so that there's always at least one node in the queue 
    // to separate the node being accessed at the front 
    // from that being accessed at the back.
    // Basically in push/pop operations, reduce/prevent both front and back being modified in a function.
    mpmc_queue() : front_mutex{}, front{new node{}}, back_mutex{}, back{front.get()} {}
    mpmc_queue(const mpmc_queue& other) = delete;
    mpmc_queue& operator=(const mpmc_queue& other) = delete;

    void enqueue(T t) {
        std::unique_ptr<node> p{new node{}}; // New dummy node
        node* const new_back = p.get();

        std::unique_lock<std::mutex> lock{back_mutex};
        back->data = std::move(t);
        back->next = std::move(p);
        back = new_back;
    }

    std::optional<T> try_dequeue() {
        std::unique_ptr<node> old_front = dequeue_front();
        if (!old_front) {
            return std::nullopt;
        }
        return std::optional<T>{std::move(old_front->data)};
    }

    bool empty() {
        std::unique_lock<std::mutex> lock{front_mutex};
        return front.get() == get_back();
    }
};
