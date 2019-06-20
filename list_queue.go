package queue

import (
	"container/list"
	"runtime"
	"sync"
	"time"
)

// ListQueue for
type ListQueue struct {
	queue      *list.List
	capacity   int32
	lock       sync.Mutex
	closedChan chan struct{}
}

// NewListQueue create ListQueue instance
func NewListQueue(capacity int32) Queue {
	q := &ListQueue{
		queue:      list.New(),
		capacity:   capacity,
		closedChan: make(chan struct{}),
	}
	return q
}

func (q *ListQueue) Pop() (interface{}, error) {

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
		if v, err := q.AsyncPop(); err == nil {
			return v, nil
		} else if err == ErrQueueIsClosed {
			return nil, err
		}
	}
}

// AsyncPop 异步读队列
func (q *ListQueue) AsyncPop() (interface{}, error) {
	if q.IsClosed() {
		return nil, ErrQueueIsClosed
	}
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.queue.Len() == 0 {
		return nil, ErrQueueEmpty
	}

	x := q.queue.Front() // 取出链表第一个元素
	q.queue.Remove(x)
	return x.Value, nil
}

func (q *ListQueue) Push(x interface{}) error {
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
		if err := q.AsyncPush(x); err == nil {
			return nil
		} else if err == ErrQueueIsClosed {
			return err
		}
	}
}

func (q *ListQueue) AsyncPush(x interface{}) error {
	if q.IsClosed() {
		return ErrQueueIsClosed
	}

	q.lock.Lock()
	defer q.lock.Unlock()

	if q.Length() >= q.capacity {
		return ErrQueueFull
	}

	q.queue.PushBack(x)
	return nil
}

func (q *ListQueue) Length() int32 {
	return int32(q.queue.Len())
}

//Capacity 队列大小
func (q *ListQueue) Capacity() int32 {
	return q.capacity
}

func (q *ListQueue) Close() error {
	if q.IsClosed() {
		return ErrQueueIsClosed
	}
	close(q.closedChan)
	return nil
}

func (q *ListQueue) IsClosed() bool {
	select {
	case <-q.closedChan:
		return true
	default:
	}
	return false
}
