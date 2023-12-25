package di

import (
	"github.com/google/wire"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
)

var CertWireSet = wire.NewSet(
	cert.ProvideCACertificate,
	cert.ProvideTransportCertificate,
)
