package queue

type Queue[T any] struct {
	ch chan T
}
