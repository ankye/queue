package queue

import "time"

//Queue interface define
type Queue interface {
	// Push 阻塞写队列
	Push(data interface{}) error
	// AsyncPush 异步非阻塞写队列
	AsyncPush(data interface{}) error
	// Pop 阻塞读队列
	Pop() (interface{}, error)
	// AsyncPop 异步读队列
	AsyncPop() (interface{}, error)
	// Capacity 队列大小
	Capacity() int
	// Length 队列占用数量
	Length() int
	// Close 关闭队列
	Close() error
	// IsClosed 队列是否已经关闭 关闭返回true,否则返回false
	IsClosed() bool
}

const TIMEOUT = time.Second * 15
