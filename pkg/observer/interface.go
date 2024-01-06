package observer

type Publisher[Event any] interface {
	Subscribe(subscriber Subscriber[Event])
	Unsubscribe(subscriber Subscriber[Event])
	NotifySubscribers(e Event)
}

type Subscriber[Event any] interface {
	OnEvent(e Event)
}
