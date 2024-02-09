package datastructure

import (
	"context"
	"go.opentelemetry.io/otel/metric"
	"slices"
)

type queueOption struct {
	size int

	// Metrics
	meter     metric.Meter
	meterName string
}

type QueueOptionFunc func(option *queueOption)

func WithQueueSize(size int) QueueOptionFunc {
	return func(q *queueOption) {
		q.size = size
	}
}

func WithQueueMetrics(meter metric.Meter, meterName string) QueueOptionFunc {
	return func(q *queueOption) {
		q.meter = meter
		q.meterName = meterName
	}
}

// Queue is a FIFO (First in Last out) queue data structure
type Queue[Data any] struct {
	data             []Data
	queueSizeCounter metric.Int64UpDownCounter
}

func ProvideQueue[Data any](ops ...QueueOptionFunc) Queue[Data] {
	var option queueOption

	for _, eachOps := range ops {
		eachOps(&option)
	}

	data := make([]Data, 0, option.size)
	result := Queue[Data]{
		data: data,
	}

	if option.meter != nil {
		counter, err := option.meter.Int64UpDownCounter(option.meterName)
		if err != nil {
			panic(err)
		}
		result.queueSizeCounter = counter
	}
	return result
}

func (q *Queue[Data]) Pop() (*Data, bool) {
	if q.Empty() {
		return nil, false
	}

	result := q.data[0]
	q.data = q.data[1:]
	q.queueSizeCounter.Add(context.Background(), -1)
	return &result, true
}

func (q *Queue[Data]) Push(d Data) {
	q.data = append(q.data, d)
	q.queueSizeCounter.Add(context.Background(), 1)
}

func (q *Queue[Data]) PopWithFilter(filter func(Data) bool) (*Data, bool) {
	if q.Empty() {
		return nil, false
	}
	idx := slices.IndexFunc(q.data, filter)
	if idx == -1 {
		return nil, false
	}
	result := q.data[idx]
	q.data = append(q.data[:idx], q.data[idx+1:]...)
	q.queueSizeCounter.Add(context.Background(), -1)
	return &result, true
}

func (q *Queue[Data]) ReadQueue() []Data {
	return q.data
}

func (q *Queue[Data]) RemoveWithCondition(removeCondition func(data Data) bool) {
	sizeBeforeRemove := len(q.data)
	q.data = Filter(q.data, removeCondition)
	sizeAfterRemove := len(q.data)
	q.queueSizeCounter.Add(context.Background(), int64(sizeBeforeRemove-sizeAfterRemove))
}

func (q *Queue[Data]) Empty() bool {
	return len(q.data) == 0
}
