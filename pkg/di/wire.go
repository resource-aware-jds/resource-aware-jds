package di

import (
	"github.com/google/wire"
	certDI "github.com/resource-aware-jds/resource-aware-jds/pkg/cert/di"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/dockerclient"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/mongo"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/taskqueue"
)

var PKGWireSet = wire.NewSet(
	mongo.ProvideMongoConnection,
	grpc.ProvideGRPCServer,
	grpc.ProvideGRPCSocketServer,
	certDI.CertWireSet,
	dockerclient.ProvideDockerClient,
	taskqueue.ProvideTaskQueue,
)
