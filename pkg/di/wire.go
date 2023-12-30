package di

import (
	"github.com/google/wire"
	certDI "github.com/resource-aware-jds/resource-aware-jds/pkg/cert/di"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/dockerclient"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/mongo"
)

var PKGWireSet = wire.NewSet(
	mongo.ProvideMongoConnection,
	grpc.ProvideGRPCServer,
	certDI.CertWireSet,
	dockerclient.ProvideDockerClient,
)
