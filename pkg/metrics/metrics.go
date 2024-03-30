package metrics

import (
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	otelMetrics "go.opentelemetry.io/otel/sdk/metric"
)

func ProvideMeter(cert cert.TransportCertificate) (metric.Meter, error) {
	exporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}
	provider := otelMetrics.NewMeterProvider(otelMetrics.WithReader(exporter))
	meter := provider.Meter(
		"rajds_performance_metric",
		metric.WithInstrumentationAttributes(attribute.KeyValue{
			Key:   "nodeID",
			Value: attribute.StringValue(cert.GetNodeID()),
		}),
	)
	return meter, nil
}
