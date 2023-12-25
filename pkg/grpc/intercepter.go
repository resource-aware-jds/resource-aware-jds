package grpc

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func grpcUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	logrusFields := logrus.WithFields(logrus.Fields{
		"protocol": "grpc",
		"type":     "CALL_DUMP",
		"method":   info.FullMethod,
	})

	logrusFields.Info("Incoming request")
	result, err := handler(ctx, req)
	if err != nil {
		logrusFields.Error("Handler response error: ", err)
	} else {
		logrusFields.Info("Handler response success")
	}

	return result, err
}
