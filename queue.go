package queue

//IQueue interface define
type IQueue interface {
	Init(capaciity int) bool
	Put(data interface{}) bool
	Get() (interface{}, bool)
	Length() int
}
