package datastructure

import "context"

type Observer[T any] interface {
	OnEvent(context.Context, T) error
}

type Observable[T any] interface {
	AddObserver(Observer[T])
	RemoveObserver(Observer[T])
	NotifyObserver(context.Context, T) error
}
