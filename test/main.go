package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/gonethopper/queue"
	"github.com/gonethopper/queue/array_queue"
	"github.com/gonethopper/queue/chan_queue"
)

const TIME = 10000000

var wg sync.WaitGroup

func push(q queue.Queue) {
	defer wg.Done()
	for i := 0; i < TIME; i++ {
		if err := q.Put(i); err != nil {
			fmt.Println(err)
		}
	}

}
func pop(q queue.Queue, pushNum int) {
	defer wg.Done()
	count := 0
	start1 := time.Now()
	sum := 0
	for {
		v, err := q.Get()
		if err == nil && v != nil {
			sum += v.(int)
			count++
		}
		if count == pushNum*TIME {
			break
		}
	}
	fmt.Println(sum)
	fmt.Println(count)
	start2 := time.Since(start1)
	fmt.Println(start2)

}
func doTest(q queue.Queue) {
	pushNum := 5
	wg.Add(pushNum + 1)

	for i := 0; i < pushNum; i++ {
		go push(q)
	}
	go pop(q, pushNum)

	wg.Wait()

}
func main() {
	capacity := 1000
	queue1 := array_queue.NewQueue(capacity)
	doTest(queue1)

	queue2 := chan_queue.NewQueue(capacity)
	doTest(queue2)

}
