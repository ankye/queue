package array_queue

import (
	"sync"
)

// Queue for
type ArrayQueue struct {
	queue []interface{}
	lock  sync.Mutex
}

func NewQueue(capaciity int) *ArrayQueue {
	q := &ArrayQueue{
		queue: make([]interface{}, 0, 100),
	}

	return q
}

func (q *ArrayQueue) Init(capaciity int) bool {
	return true
}

func (q *ArrayQueue) Get() (val interface{}, ok bool) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if len(q.queue) == 0 {
		return nil, false
	}

	//n := len(q.queue)
	x := q.queue[0]
	q.queue = q.queue[1:]
	return x, true
}

func (q *ArrayQueue) Put(x interface{}) bool {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.queue = append(q.queue, x)
	return true
}

func (q *ArrayQueue) Length() int {
	return len(q.queue)
}
