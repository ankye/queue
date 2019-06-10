package queue

//Queue interface define
type Queue interface {
	Init(capacity int) error
	//Put 阻塞写队列
	Put(data interface{}) error
	//AsyncPut 异步非阻塞写队列
	AsyncPut(data interface{}) error
	//Get 阻塞读队列
	Get() (interface{}, error)
	//AsyncGet 异步读队列
	AsyncGet() (interface{}, error)
	//Capacity 队列大小
	Capacity() (int, error)
	//队列占用数量
	Length() (int, error)
	//Close 关闭队列
	Close() error
	//IsClosed 队列是否已经关闭 关闭返回true,否则返回false
	IsClosed() bool
}
