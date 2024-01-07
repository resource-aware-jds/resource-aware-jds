package di

import (
	"github.com/google/wire"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
)

var ControlPlaneCertWireSet = wire.NewSet(
	cert.ProvideCACertificate,
	cert.ProvideTransportCertificate,
	cert.ProvideWorkerNodeCACertificate,
)

var WorkerNodeCertWireSet = wire.NewSet(
	cert.ProvideWorkerNodeCACertificate,
	cert.ProvideWorkerNodeTransportCertificate,
)
