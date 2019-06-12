package queue_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/gonethopper/queue"
)

func Benchmark_ArrayQueue(b *testing.B) {
	queue := queue.NewArrayQueue(1000)
	pushNum := 5
	num := b.N
	// num := 10000000
	// var wg sync.WaitGroup
	for i := 0; i < pushNum; i++ {
		// wg.Add(1)
		go func(l int) {
			for i := 0; i < l; i++ {
				queue.Push(i)
			}
			// wg.Done()
		}(num)
	}
	idx := 0
	for idx < num*pushNum {
		_, err := queue.Pop()
		if err == nil {
			idx += 1
		}
	}
}

func Benchmark_ListQueue(b *testing.B) {
	queue := queue.NewListQueue(1000)
	pushNum := 5
	num := b.N
	// num := 10000000
	// var wg sync.WaitGroup
	for i := 0; i < pushNum; i++ {
		// wg.Add(1)
		go func(l int) {
			for i := 0; i < l; i++ {
				queue.Push(i)
			}
			// wg.Done()
		}(num)
	}
	idx := 0
	for idx < num*pushNum {
		_, err := queue.Pop()
		if err == nil {
			idx += 1
		}
	}
}

func Benchmark_ChanQueue(b *testing.B) {
	queue := queue.NewChanQueue(1000)
	pushNum := 5
	num := b.N
	// num := 10000000
	// var wg sync.WaitGroup
	for i := 0; i < pushNum; i++ {
		// wg.Add(1)
		go func(l int) {
			for i := 0; i < l; i++ {
				queue.Push(i)
			}
			// wg.Done()
		}(num)
	}
	idx := 0
	for idx < num*pushNum {
		_, err := queue.Pop()
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
			q.Push(v)
		}

	}()

	for _, v := range list {
		v2, _ := q.Pop()
		if v != v2 {
			fmt.Printf("vail, %d != %d", v, v2)
			return false
		}
	}

	return true
}

func Test_ArrayQueue(t *testing.T) {
	q := queue.NewArrayQueue(1000)
	r := test_queue(q, 10000)
	if r == false {
		t.Error("array queue error")
	}
}
func Test_ListQueue(t *testing.T) {
	q := queue.NewListQueue(1000)
	r := test_queue(q, 10000)
	if r == false {
		t.Error("list queue error")
	}
}
func Test_ChanQueue(t *testing.T) {
	q := queue.NewChanQueue(1000)
	r := test_queue(q, 10000)
	if r == false {
		t.Error("chan queue error")
	}
}
