package array_queue

import (
	"errors"
	"sync"
)

var ErrQueueFull = errors.New("Queue is full")
var ErrQueueEmpty = errors.New("Queue is empty")
var ErrQueueIsClosed = errors.New("Queue is Closed")

// Queue for
type ArrayQueue struct {
	queue      []interface{}
	capacity   int
	lock       sync.Mutex
	closedChan chan struct{}
}

func NewQueue(capacity int) *ArrayQueue {
	q := &ArrayQueue{
		queue:      make([]interface{}, 0, capacity),
		capacity:   capacity,
		closedChan: make(chan struct{}),
	}

	return q
}

func (q *ArrayQueue) Init(capacity int) error {
	return nil
}

func (q *ArrayQueue) Get() (interface{}, error) {

	if q.IsClosed() {
		return nil, ErrQueueIsClosed
	}

	if len(q.queue) == 0 {
		return nil, ErrQueueEmpty
	}
	q.lock.Lock()
	defer q.lock.Unlock()

	x := q.queue[0]
	q.queue = q.queue[1:]
	return x, nil
}

//AsyncGet 异步读队列
func (q *ArrayQueue) AsyncGet() (interface{}, error) {
	if q.IsClosed() {
		return nil, ErrQueueIsClosed
	}
	if len(q.queue) == 0 {
		return nil, ErrQueueEmpty
	}
	q.lock.Lock()
	defer q.lock.Unlock()
	x := q.queue[0]
	q.queue = q.queue[1:]
	return x, nil
}

func (q *ArrayQueue) Put(x interface{}) error {

	if len, err := q.Length(); err == nil {
		if len >= q.capacity {
			return ErrQueueFull
		}
	} else {
		return err
	}
	q.lock.Lock()
	defer q.lock.Unlock()
	q.queue = append(q.queue, x)
	return nil
}

func (q *ArrayQueue) AsyncPut(x interface{}) error {
	if len, err := q.Length(); err == nil {
		if len >= q.capacity {
			return ErrQueueFull
		}
	} else {
		return err
	}
	q.lock.Lock()
	defer q.lock.Unlock()
	q.queue = append(q.queue, x)
	return nil
}

func (q *ArrayQueue) Length() (int, error) {
	return len(q.queue), nil
}

//Capacity 队列大小
func (q *ArrayQueue) Capacity() (int, error) {
	return q.capacity, nil
}

func (q *ArrayQueue) Close() error {
	if q.IsClosed() {
		return ErrQueueIsClosed
	}
	close(q.closedChan)
	return nil
}

func (q *ArrayQueue) IsClosed() bool {
	select {
	case <-q.closedChan:
		return true
	default:
	}
	return false
}
