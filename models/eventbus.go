package models

import (
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventType uint32

const (
	SuccessTaskEventType EventType = 0
	FailTaskEventType    EventType = 0
)

type TaskEventBus struct {
	EventType  EventType
	RawRequest *proto.ReportSuccessTaskRequest
	TaskID     primitive.ObjectID
}
