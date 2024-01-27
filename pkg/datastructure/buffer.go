package datastructure

import (
	"github.com/sirupsen/logrus"
	"reflect"
)

type Buffer[T comparable, U any] map[T]U

func (t Buffer[T, U]) Store(id T, object U) {
	logrus.Info("Buffer ", reflect.TypeOf(object), " with id: ", id)
	t[id] = object
}

func (t Buffer[T, U]) Pop(id T) *U {
	object := t.Get(id)
	if object == nil {
		return nil
	}
	delete(t, id)
	return object
}

func (t Buffer[T, U]) Get(id T) *U {
	object, ok := t[id]
	if !ok {
		return nil
	}
	return &object
}

func (t Buffer[T, U]) IsObjectInBuffer(id T) bool {
	_, ok := t[id]
	return ok
}

func (t Buffer[T, U]) GetKeys() []T {
	keys := make([]T, len(t))
	for k := range t {
		keys = append(keys, k)
	}
	return keys
}
