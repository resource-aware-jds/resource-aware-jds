package metrics

import "fmt"

func GenerateControlPlaneMetric(metricName string) string {
	return fmt.Sprintf("rajds_cp_%s", metricName)
}

func GenerateWorkerNodeMetric(metricName string) string {
	return fmt.Sprintf("rajds_cp_%s", metricName)
}
