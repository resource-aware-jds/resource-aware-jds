package buffer

import (
	"github.com/sirupsen/logrus"
	"reflect"
)

type Buffer[T comparable, U any] struct {
	commonMap map[T]*U
}

func (t *Buffer[T, U]) Store(id T, object *U) {
	logrus.Info("Buffer ", reflect.TypeOf(object), " with id: ", id)
	t.commonMap[id] = object
}

func (t *Buffer[T, U]) Pop(id T) *U {
	object, ok := t.commonMap[id]
	logrus.Info("Remove ", reflect.TypeOf(object), " from Buffer with id: ", id)
	if !ok {
		return nil
	}
	delete(t.commonMap, id)
	return object
}

func (t *Buffer[T, U]) IsObjectInBuffer(id T) bool {
	_, ok := t.commonMap[id]
	return ok
}

func (t *Buffer[T, U]) GetKeys() []T {
	keys := make([]T, 0, len(t.commonMap))
	for k := range t.commonMap {
		keys = append(keys, k)
	}
	return keys
}

func (t *Buffer[T, U]) InitializeMap() {
	t.commonMap = make(map[T]*U)
}
