package metrics

import (
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	otelMetrics "go.opentelemetry.io/otel/sdk/metric"
)

func ProvideMeter() (metric.Meter, error) {
	exporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}
	provider := otelMetrics.NewMeterProvider(otelMetrics.WithReader(exporter))
	meter := provider.Meter("rajds_performance_metric")
	return meter, nil
}
