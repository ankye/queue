package chan_queue

type ChanQueue struct {
	innerChan chan interface{}
}

func NewQueue(capaciity int) *ChanQueue {
	return &ChanQueue{
		innerChan: make(chan interface{}, capaciity),
	}
}

func (q *ChanQueue) Init(capaciity int) bool {
	return true
}

func (q *ChanQueue) Get() (val interface{}, ok bool) {
	v := <-q.innerChan
	return v, true

}
func (q *ChanQueue) Put(x interface{}) bool {
	q.innerChan <- x
	return true
}
func (q *ChanQueue) Length() int {
	return 0
}
