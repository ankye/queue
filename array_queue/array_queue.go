package array_queue

import (
	"runtime"
	"sync"
	"time"

	. "github.com/gonethopper/queue/error"
)

const TIMEOUT = time.Second * 15

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

	var i int
	for start := time.Now(); ; {

		if i>>3 == 1 {
			i = 1
			if time.Since(start) > TIMEOUT {
				return nil, ErrQueueTimeout
			}
			runtime.Gosched()
		}
		i++
		if v, err := q.AsyncGet(); err == nil {
			return v, nil
		} else if err == ErrQueueIsClosed {
			return nil, err
		}
	}
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
	var i int
	for start := time.Now(); ; {
		if i>>3 == 1 {
			i = 1
			if time.Since(start) > TIMEOUT {
				return ErrQueueTimeout
			}
			runtime.Gosched()
		}
		i++
		if err := q.AsyncPut(x); err == nil {
			return nil
		} else if err == ErrQueueIsClosed {
			return err
		}
	}
}

func (q *ArrayQueue) AsyncPut(x interface{}) error {
	if q.IsClosed() {
		return ErrQueueIsClosed
	}
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
