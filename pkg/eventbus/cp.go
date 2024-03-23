package eventbus

import (
	"context"
	"errors"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/datastructure"
	"sync"
)

type TaskEventBus datastructure.Observable[models.TaskEventBus]

type baseEventBus[T any] struct {
	mutex        sync.Mutex
	observerList map[datastructure.Observer[T]]bool
}

func (s *baseEventBus[T]) AddObserver(o datastructure.Observer[T]) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.observerList[o] = true
}

func (s *baseEventBus[T]) RemoveObserver(o datastructure.Observer[T]) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.observerList, o)
}

func (s *baseEventBus[T]) NotifyObserver(ctx context.Context, data T) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	var wg sync.WaitGroup

	var errListMutex sync.Mutex
	var errResponse error
	for observer := range s.observerList {
		wg.Add(1)
		go func(obs datastructure.Observer[T]) {
			err := obs.OnEvent(ctx, data)
			errListMutex.Lock()
			errResponse = errors.Join(errResponse, err)
			errListMutex.Unlock()
			wg.Done()
		}(observer)
	}
	wg.Wait()
	return errResponse
}

func ProvideTaskEventBus() TaskEventBus {
	return &baseEventBus[models.TaskEventBus]{
		observerList: make(map[datastructure.Observer[models.TaskEventBus]]bool),
	}
}
