package pool_queue

import (
	"runtime"
	"sync"
	"time"

	. "github.com/gonethopper/queue/error"
)

var node_pool *sync.Pool = &sync.Pool{New: func() interface{} { return new(Node) }}
var p_pool *sync.Pool = &sync.Pool{New: func() interface{} { return make([]*Node, 0, NODE_POOL_CELL_COUNT*2) }}

const NODE_POOL_CELL_COUNT = 8
const NODE_POOL_MAX_COUNT = 4096
const TIMEOUT = time.Second * 15

type nodePool struct {
	p []*Node
	c int
}

func newNodePool() *nodePool {
	t := &nodePool{
		p: p_pool.Get().([]*Node),
	}
	return t
}

func (t *nodePool) Get() *Node {
	var x *Node = nil
	last := len(t.p) - 1
	if last >= 0 {
		x = t.p[last]
		t.p = t.p[:last]
	} else {
		x = node_pool.Get().(*Node)
	}
	return x
}

func (t *nodePool) Put(x *Node) {
	x.reset()
	if len(t.p) < cap(t.p) {
		t.p = append(t.p, x)
		t.c = 0
	} else {
		t.c++
		if t.c < NODE_POOL_CELL_COUNT || cap(t.p) >= NODE_POOL_MAX_COUNT {
			node_pool.Put(x)
			return
		}
		newCap := (2*NODE_POOL_CELL_COUNT + cap(t.p))
		s := make([]*Node, len(t.p), newCap)
		copy(s, t.p)
		o := t.p[:0]
		t.p = s
		p_pool.Put(o)
		t.c = 0
		t.p = append(t.p, x)
	}
}

type Node struct {
	data interface{}
	next *Node
}

func (n *Node) reset() {
	n.data = nil
	n.next = nil
}

//PoolQueue struct
type PoolQueue struct {
	head       *Node
	end        *Node
	node_pool  *nodePool
	size       int
	lock       sync.Mutex
	closedChan chan struct{}
	capacity   int
}

//NewQueue create queue instance
func NewQueue(capacity int) *PoolQueue {
	q := &PoolQueue{
		head:       nil,
		end:        nil,
		node_pool:  newNodePool(),
		size:       0,
		capacity:   capacity,
		closedChan: make(chan struct{}),
	}
	return q
}

//Init queue init
func (q *PoolQueue) Init(capaciity int) error {
	return nil
}

//Put push one object to queue
func (q *PoolQueue) Put(x interface{}) error {
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

//Get get one object from queue

func (q *PoolQueue) Get() (interface{}, error) {

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

//Length objects count in queue
func (q *PoolQueue) Length() (int, error) {
	return q.size, nil
}

//Empty queue is empty return true else return false
func (q *PoolQueue) Empty() bool {
	if q.head == nil {
		return true
	}
	return false
}

func (q *PoolQueue) AsyncPut(data interface{}) error {
	if q.IsClosed() {
		return ErrQueueIsClosed
	}

	q.lock.Lock()
	defer q.lock.Unlock()
	if q.size >= q.capacity {
		return ErrQueueFull
	}
	n := q.node_pool.Get()
	n.data = data
	n.next = nil
	if q.end == nil {
		q.head = n
		q.end = n
	} else {
		q.end.next = n
		q.end = n
	}
	q.size++

	return nil
}

//AsyncGet 异步读队列
func (q *PoolQueue) AsyncGet() (interface{}, error) {
	if q.IsClosed() {
		return nil, ErrQueueIsClosed
	}
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.head == nil {
		return nil, ErrQueueEmpty
	}

	n := q.head
	data := n.data
	q.head = n.next
	if q.head == nil {
		q.end = nil
	}
	q.size--
	q.node_pool.Put(n)

	return data, nil
}

//Capacity 队列大小
func (q *PoolQueue) Capacity() (int, error) {
	return q.capacity, nil
}

func (q *PoolQueue) Close() error {
	if q.IsClosed() {
		return ErrQueueIsClosed
	}
	close(q.closedChan)

	return nil
}

func (q *PoolQueue) IsClosed() bool {
	select {
	case <-q.closedChan:
		return true
	default:
	}
	return false

}
