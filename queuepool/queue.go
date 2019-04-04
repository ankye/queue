package queuepool

import (
	"sync"
)

var node_pool *sync.Pool = &sync.Pool{New: func() interface{} { return new(Node) }}
var p_pool *sync.Pool = &sync.Pool{New: func() interface{} { return make([]*Node, 0, NODE_POOL_CELL_COUNT*2) }}

const NODE_POOL_CELL_COUNT = 8
const NODE_POOL_MAX_COUNT = 4096

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

//Queue struct
type Queue struct {
	head      *Node
	end       *Node
	node_pool *nodePool
	c         int
	lock      sync.Mutex
}

//NewQueue create queue instance
func NewQueue(capaciity int) *Queue {
	q := &Queue{
		head:      nil,
		end:       nil,
		node_pool: newNodePool(),
		c:         0,
	}
	return q
}

//Init queue init
func (q *Queue) Init(capaciity int) bool {
	return true
}

//Put push one object to queue
func (q *Queue) Put(data interface{}) bool {
	q.lock.Lock()
	defer q.lock.Unlock()
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
	q.c++

	return true
}

//Get get one object from queue
func (q *Queue) Get() (interface{}, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.head == nil {
		return nil, false
	}

	n := q.head
	data := n.data
	q.head = n.next
	if q.head == nil {
		q.end = nil
	}
	q.c--
	q.node_pool.Put(n)

	return data, true
}

//Length objects count in queue
func (q *Queue) Length() int {
	return q.c
}

//Empty queue is empty return true else return false
func (q *Queue) Empty() bool {
	if q.head == nil {
		return true
	}
	return false
}
