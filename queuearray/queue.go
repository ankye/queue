package queuearray

import (
	"sync"
)

// Queue for
type Queue struct {
	queue []interface{}
	Len   chan chan int
	lock  sync.Mutex
}

func NewQueue(capaciity int) *Queue {
	q := &Queue{
		queue: make([]interface{}, 0, 100),
		Len:   make(chan chan int)}

	return q
}

func (q *Queue) Init(capaciity int) bool {
	return true
}

func (q *Queue) Get() (val interface{}, ok bool) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if len(q.queue) == 0 {
		return nil, false
	}

	n := len(q.queue)
	x := q.queue[n-1]
	q.queue = q.queue[0 : n-1]

	return x, true
}

func (q *Queue) Put(x interface{}) bool {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.queue = append(q.queue, x)
	return true
}

func (q *Queue) Length() int {
	return len(q.queue)
}
