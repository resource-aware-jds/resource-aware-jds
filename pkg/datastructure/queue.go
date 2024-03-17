package datastructure

import (
	"context"
	"go.opentelemetry.io/otel/metric"
	"slices"
)

type queueOption struct {
	size int

	// Metrics
	meter         metric.Meter
	metricName    string
	metricOptions []metric.Int64ObservableCounterOption
}

type QueueOptionFunc func(option *queueOption)

func WithQueueSize(size int) QueueOptionFunc {
	return func(q *queueOption) {
		q.size = size
	}
}

func WithQueueMetrics(meter metric.Meter, metricName string, opts ...metric.Int64ObservableCounterOption) QueueOptionFunc {
	return func(q *queueOption) {
		q.meter = meter
		q.metricName = metricName
		q.metricOptions = opts
	}
}

// Queue is a FIFO (First in Last out) queue data structure
type Queue[Data any] struct {
	data             []Data
	queueSizeCounter metric.Int64ObservableCounter
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
		option.metricOptions = append(
			option.metricOptions,
			metric.WithInt64Callback(func(ctx context.Context, observer metric.Int64Observer) error {
				observer.Observe(int64(len(data)))
				return nil
			}),
		)

		counter, err := option.meter.Int64ObservableCounter(
			option.metricName,
			option.metricOptions...,
		)
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
	if idx == -1 {
		return nil, false
	}
	result := q.data[idx]
	q.data = append(q.data[:idx], q.data[idx+1:]...)
	return &result, true
}

func (q *Queue[Data]) ReadQueue() []Data {
	return q.data
}

func (q *Queue[Data]) RemoveWithCondition(removeCondition func(data Data) bool) {
	q.data = Filter(q.data, removeCondition)
}

func (q *Queue[Data]) Empty() bool {
	return len(q.data) == 0
}
