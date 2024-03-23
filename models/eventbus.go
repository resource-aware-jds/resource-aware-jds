package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventType uint32

const (
	SuccessTaskEventType EventType = 0
	FailTaskEventType    EventType = 0
)

type TaskEventBus struct {
	EventType EventType
	TaskID    primitive.ObjectID
}
