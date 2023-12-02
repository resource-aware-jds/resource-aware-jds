package di

import (
	"github.com/google/wire"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/mongo"
)

var PKGWireSet = wire.NewSet(
	mongo.ProvideMongoConnection,
	grpc.ProvideGRPCServer,
)
