package di

import (
	"github.com/google/wire"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/distribution"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/dockerclient"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/http"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/metrics"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/mongo"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/pool"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/taskqueue"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/workerlogic"
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
	workerlogic.ProvideDelayWorkerDistributor,
	http.ProvideHttpServer,
	grpc.ProvideRAJDSGRPCResolver,
	metrics.ProvideMeter,
)
