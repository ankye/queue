package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/ankye/queue"
	"github.com/ankye/queue/queuearray"
	"github.com/ankye/queue/queuepool"
)

const TIME = 9999999

var wg sync.WaitGroup

func push(q queue.IQueue) {
	for i := 0; i < TIME; i++ {
		q.Put(i)
	}
	wg.Done()
}
func pop(q queue.IQueue, pushNum int) {
	count := 0
	start1 := time.Now()
	sum := 0
	for i := 0; i < 100*TIME; i++ {
		v, ok := q.Get()
		if ok && v != nil {
			sum += v.(int)
			count++
		}
		if count == pushNum*TIME {
			break
		}
	}
	fmt.Println(sum)
	start2 := time.Since(start1)
	fmt.Println(start2)
	wg.Done()
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
	queue1 := queuearray.NewQueue(1000)
	doTest(queue1)

	queue2 := queuepool.NewQueue(1000)
	doTest(queue2)

}
