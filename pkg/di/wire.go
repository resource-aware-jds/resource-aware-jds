package di

import (
	"github.com/google/wire"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/distribution"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/dockerclient"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/mongo"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/pool"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/taskBuffer"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/taskqueue"
)

var PKGWireSet = wire.NewSet(
	mongo.ProvideMongoConnection,
	grpc.ProvideGRPCServer,
	grpc.ProvideGRPCSocketServer,
	grpc.ProvideRAJDSGrpcClient,
	dockerclient.ProvideDockerClient,
	taskqueue.ProvideTaskQueue,
	taskBuffer.ProvideTaskBuffer,
	pool.ProvideWorkerNode,
	distribution.ProvideRoundRobinDistributor,
)
