package datastructure

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/metric"
	"reflect"
)

type Buffer[T comparable, U any] map[T]U

type bufferOption struct {
	// Metrics
	meter      metric.Meter
	metricName string
	metricOpts []metric.Int64ObservableCounterOption
}

type BufferOptionFunc func(option *bufferOption)

func WithBufferMetrics(meter metric.Meter, metricName string, opts ...metric.Int64ObservableCounterOption) BufferOptionFunc {
	return func(q *bufferOption) {
		q.meter = meter
		q.metricName = metricName
		q.metricOpts = opts
	}
}

func ProvideBuffer[T comparable, U any](ops ...BufferOptionFunc) Buffer[T, U] {
	data := make(map[T]U)

	var option bufferOption

	for _, eachOps := range ops {
		eachOps(&option)
	}
	if option.meter != nil {
		option.metricOpts = append(
			option.metricOpts,
			metric.WithInt64Callback(func(ctx context.Context, observer metric.Int64Observer) error {
				observer.Observe(int64(len(data)))
				return nil
			}),
		)

		_, err := option.meter.Int64ObservableCounter(
			option.metricName,
			option.metricOpts...,
		)
		if err != nil {
			panic(err)
		}
	}

	return data
}

func (t *Buffer[T, U]) Store(id T, object U) {
	logrus.Info("Buffer ", reflect.TypeOf(object), " with id: ", id)
	(*t)[id] = object
}

func (t *Buffer[T, U]) Pop(id T) *U {
	object := t.Get(id)
	if object == nil {
		return nil
	}
	delete(*t, id)
	return object
}

func (t *Buffer[T, U]) Get(id T) *U {
	object, ok := (*t)[id]
	if !ok {
		return nil
	}
	return &object
}

func (t *Buffer[T, U]) IsObjectInBuffer(id T) bool {
	_, ok := (*t)[id]
	return ok
}

func (t *Buffer[T, U]) GetKeys() []T {
	keys := make([]T, 0, len(*t))
	for k := range *t {
		keys = append(keys, k)
	}
	return keys
}

func (t *Buffer[T, U]) GetValues() []U {
	values := make([]U, 0)
	for _, v := range *t {
		values = append(values, v)
	}
	return values
}
