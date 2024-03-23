package datastructure

import "context"

//go:generate mockgen -source=./observer.go -destination=./mock_datastructure/mock_observer.go -package=mock_datastructure

type Observer[T any] interface {
	OnEvent(context.Context, T) error
}

type Observable[T any] interface {
	AddObserver(Observer[T])
	RemoveObserver(Observer[T])
	NotifyObserver(context.Context, T) error
}
