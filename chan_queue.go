package queue

import (
	"sync/atomic"
	"time"
)

type ChanQueue struct {
	innerChan  chan interface{}
	capacity   int
	size       int32
	timer      *time.Timer
	closedChan chan struct{}
}

func NewChanQueue(capacity int) Queue {
	return &ChanQueue{
		innerChan:  make(chan interface{}, capacity),
		capacity:   capacity,
		size:       0,
		timer:      time.NewTimer(time.Second),
		closedChan: make(chan struct{}),
	}
}

func (q *ChanQueue) Init(capacity int) error {
	return nil
}

func (q *ChanQueue) Get() (val interface{}, err error) {

	v, ok := <-q.innerChan
	if ok {
		atomic.AddInt32(&q.size, -1)
		return v, nil
	}
	return nil, ErrQueueIsClosed

}
func (q *ChanQueue) AsyncGet() (val interface{}, err error) {

	select {
	case v, ok := <-q.innerChan:
		if ok {
			atomic.AddInt32(&q.size, -1)
			return v, nil
		}
		return nil, ErrQueueIsClosed
	default:
		return nil, ErrQueueEmpty
	}

}
func (q *ChanQueue) Put(x interface{}) error {

	if q.IsClosed() {
		return ErrQueueIsClosed
	}
	q.innerChan <- x
	atomic.AddInt32(&q.size, 1)
	return nil
}

func (q *ChanQueue) AsyncPut(x interface{}) error {

	if q.IsClosed() {
		return ErrQueueIsClosed
	}
	select {
	case q.innerChan <- x:
		atomic.AddInt32(&q.size, 1)
		return nil
	default:
		return ErrQueueFull
	}
}

func (q *ChanQueue) Length() int {
	return int(q.size)
}

func (q *ChanQueue) Capacity() int {
	return q.capacity
}

//Close 不需要关闭innerChan,交给GC回收,多写的时候直接关闭innerChan会出问题
func (q *ChanQueue) Close() error {
	if q.IsClosed() {
		return ErrQueueIsClosed
	}
	close(q.closedChan)

	return nil
}

func (q *ChanQueue) IsClosed() bool {
	select {
	case <-q.closedChan:
		return true
	default:
	}
	return false

}

func (q *ChanQueue) GetChan(timeout time.Duration) (<-chan interface{}, <-chan error) {
	timeoutChan := make(chan error, 1)
	resultChan := make(chan interface{}, 1)
	go func() {
		if timeout < 0 {
			item := <-q.innerChan
			atomic.AddInt32(&q.size, -1)
			resultChan <- item
		} else {
			select {
			case item := <-q.innerChan:
				atomic.AddInt32(&q.size, -1)
				resultChan <- item
			case <-time.After(timeout):
				timeoutChan <- ErrQueueTimeout
			}
		}
	}()
	return resultChan, timeoutChan
}
