package datastructure

import (
	"fmt"
	"sync/atomic"
)

type RoundRobin[T any] interface {
	Next() T
}

type roundRobin[T any] struct {
	data []T
	next uint32
}

func ProvideRoundRobin[T any](data ...T) (RoundRobin[T], error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("round robin require at least 1 element")
	}

	return &roundRobin[T]{
		data: data,
	}, nil
}

func (r *roundRobin[T]) Next() T {
	n := atomic.AddUint32(&r.next, 1)
	return r.data[(int(n)-1)%len(r.data)]
}
