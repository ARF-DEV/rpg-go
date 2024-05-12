package game

type Queue[T any] struct {
	values []T
}

func (q *Queue[T]) Put(val T) {
	q.values = append(q.values, val)
}

func (q *Queue[T]) Pop() T {
	val := q.values[0]
	q.values = q.values[1:]
	return val
}

func (q *Queue[T]) Len() int {
	return len(q.values)
}
