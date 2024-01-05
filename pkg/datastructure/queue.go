package datastructure

import "slices"

type Queue[Data any] struct {
	data []Data
}

func ProvideQueue[Data any](size int) Queue[Data] {
	return Queue[Data]{
		data: make([]Data, size),
	}
}

func (q *Queue[Data]) Pop() (*Data, bool) {
	if q.Empty() {
		return nil, false
	}

	result := q.data[0]
	q.data = q.data[1:]
	return &result, true
}

func (q *Queue[Data]) Push(d Data) {
	q.data = append(q.data, d)
}

func (q *Queue[Data]) PopWithFilter(filter func(Data) bool) (*Data, bool) {
	if q.Empty() {
		return nil, false
	}

	idx := slices.IndexFunc(q.data, filter)
	result := q.data[idx]
	q.data = append(q.data[:idx], q.data[idx+1:]...)
	return &result, true
}

func (q *Queue[Data]) Empty() bool {
	return len(q.data) == 0
}
