package queue_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/gonethopper/queue"
	"github.com/gonethopper/queue/array_queue"
	"github.com/gonethopper/queue/chan_queue"
	"github.com/gonethopper/queue/pool_queue"
)

func Benchmark_Queue1(b *testing.B) {
	queue := array_queue.NewQueue(1000)
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
		_, err := queue.Get()
		if err == nil {
			idx += 1
		}
	}
}

func Benchmark_Queue2(b *testing.B) {
	queue := pool_queue.NewQueue(1000)
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
		_, err := queue.Get()
		if err == nil {
			idx += 1
		}
	}
}

func Benchmark_Queue3(b *testing.B) {
	queue := chan_queue.NewQueue(1000)
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
		_, err := queue.Get()
		if err == nil {
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

func test_queue(q queue.Queue, size int) bool {
	list := genRandomList(size)
	fmt.Printf("test_queue, %d\n", size)

	go func() {
		for _, v := range list {
			// fmt.Printf("put, %d\n", v)
			q.Put(v)
		}

	}()

	for _, v := range list {
		v2, _ := q.Get()
		if v != v2 {
			fmt.Printf("vail, %d != %d", v, v2)
			return false
		}
	}

	return true
}

func Test_PoolQueue(t *testing.T) {
	q := pool_queue.NewQueue(1000)
	r := test_queue(q, 10000)
	if r == false {
		t.Error("pool queue error")
	}
}

func Test_ArrayQueue(t *testing.T) {
	q := array_queue.NewQueue(1000)
	r := test_queue(q, 10000)
	if r == false {
		t.Error("array queue error")
	}
}

func Test_ChanQueue(t *testing.T) {
	q := chan_queue.NewQueue(1000)
	r := test_queue(q, 10000)
	if r == false {
		t.Error("chan queue error")
	}
}
