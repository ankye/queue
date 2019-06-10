package queue

import "errors"

var ErrQueueFull = errors.New("Queue is full")
var ErrQueueEmpty = errors.New("Queue is empty")
var ErrQueueIsClosed = errors.New("Queue is Closed")
var ErrQueueTimeout = errors.New("Timeout waiting on queue")
