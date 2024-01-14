package di

import (
	"github.com/google/wire"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/distribution"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/dockerclient"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/mongo"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/pool"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/taskqueue"
)

var PKGWireSet = wire.NewSet(
	mongo.ProvideMongoConnection,
	grpc.ProvideGRPCServer,
	grpc.ProvideWorkerNodeReceiverGRPCServer,
	grpc.ProvideRAJDSGrpcClient,
	dockerclient.ProvideDockerClient,
	taskqueue.ProvideTaskQueue,
	pool.ProvideWorkerNode,
	distribution.ProvideRoundRobinDistributor,
)
