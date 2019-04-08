package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/ankye/queue"
	"github.com/ankye/queue/array_queue"
	"github.com/ankye/queue/chan_queue"
)

const TIME = 1000000

var wg sync.WaitGroup

func push(q queue.IQueue) {
	defer wg.Done()
	for i := 0; i < TIME; i++ {
		if err := q.Put(i); err != nil {
			fmt.Println(err)
		}
	}

}
func pop(q queue.IQueue, pushNum int) {
	defer wg.Done()
	count := 0
	start1 := time.Now()
	sum := 0
	for i := 0; i < 10000*TIME; i++ {
		v, err := q.AsyncGet()
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
func doTest(q queue.IQueue) {
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

	// queue2 := pool_queue.NewQueue(capacity)
	// doTest(queue2)

	queue3 := chan_queue.NewQueue(capacity)
	doTest(queue3)

}
