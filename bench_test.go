package queue_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/ankye/queue"
	"github.com/ankye/queue/chan_quque"
	"github.com/ankye/queue/queuearray"
	"github.com/ankye/queue/queuepool"
)

func Benchmark_Queue1(b *testing.B) {
	queue := queuearray.NewQueue(1000)
	pushNum := 5
	num := b.N
	// num := 10000000
	// var wg sync.WaitGroup
	for i := 0; i < pushNum; i++ {
		// wg.Add(1)
		go func(l int) {
			for i := 0; i < l; i++ {
				queue.Put(i)
			}
			// wg.Done()
		}(num)
	}
	idx := 0
	for idx < num*pushNum {
		_, ok := queue.Get()
		if ok {
			idx += 1
		}
	}
}

func Benchmark_Queue2(b *testing.B) {
	queue := queuepool.NewQueue(1000)
	pushNum := 5
	num := b.N
	// num := 10000000
	// var wg sync.WaitGroup
	for i := 0; i < pushNum; i++ {
		// wg.Add(1)
		go func(l int) {
			for i := 0; i < l; i++ {
				queue.Put(i)
			}
			// wg.Done()
		}(num)
	}
	idx := 0
	for idx < num*pushNum {
		_, ok := queue.Get()
		if ok {
			idx += 1
		}
	}
}

func Benchmark_Queue3(b *testing.B) {
	queue := chan_quque.NewQueue(1000)
	pushNum := 5
	num := b.N
	// num := 10000000
	// var wg sync.WaitGroup
	for i := 0; i < pushNum; i++ {
		// wg.Add(1)
		go func(l int) {
			for i := 0; i < l; i++ {
				queue.Put(i)
			}
			// wg.Done()
		}(num)
	}
	idx := 0
	for idx < num*pushNum {
		_, ok := queue.Get()
		if ok {
			idx += 1
		}
	}
}
func genRandomList(size int) []int {
	list := make([]int, size)
	for i, _ := range list {
		v := rand.Intn(10)
		list[i] = v
	}
	return list
}

func test_queue(q queue.IQueue, size int) bool {
	list := genRandomList(size)
	fmt.Printf("test_queue, %d\n", size)
	for _, v := range list {
		// fmt.Printf("put, %d\n", v)
		q.Put(v)
	}
	for _, v := range list {
		v2, _ := q.Get()
		if v != v2 {
			fmt.Printf("vail, %d != %d", v, v2)
			return false
		}
	}
	return true
}

func Test_Queue1(t *testing.T) {
	q := queuepool.NewQueue(1000)
	r := test_queue(q, 10000)
	if r == false {
		t.Error("queuepool error")
	}
}

func Test_Queue2(t *testing.T) {
	q := queuearray.NewQueue(1000)
	r := test_queue(q, 10000)
	if r == false {
		t.Error("queuepool error")
	}
}

func Test_Queue3(t *testing.T) {
	q := chan_quque.NewQueue(20000)
	r := test_queue(q, 10000)
	if r == false {
		t.Error("queuepool error")
	}
}
