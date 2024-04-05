package distribution

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/metrics"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"time"
)

type baseDistributor struct {
	name                 models.DistributorName
	logger               *logrus.Entry
	metricCounter        metric.Int64Counter
	failureMetricCounter metric.Int64Counter
}

func newBaseDistributor(name models.DistributorName, meter metric.Meter) baseDistributor {
	counter, _ := meter.Int64Counter(
		metrics.GenerateControlPlaneMetric("distributor_total_success_task"),
		metric.WithUnit("Task"),
		metric.WithDescription("The total task(s) that has been distributed"),
	)

	failureCounter, _ := meter.Int64Counter(
		metrics.GenerateControlPlaneMetric("distributor_total_failure_task"),
		metric.WithUnit("Task"),
		metric.WithDescription("The total task(s) that failed to be distributed"),
	)
	logger := logrus.WithFields(logrus.Fields{
		"role":             "Control Plane",
		"component":        "distributor",
		"distributor_name": name,
	})

	return baseDistributor{
		name:                 name,
		metricCounter:        counter,
		logger:               logger,
		failureMetricCounter: failureCounter,
	}
}

func (b *baseDistributor) distributeToNode(ctx context.Context, node NodeMapper, task models.Task, successTask *[]models.Task, errorTask *[]models.DistributeError) {
	logger := b.logger.WithFields(node.Logger.Data).WithField("taskID", task.ID.Hex())
	logger.Info("Sending task to the worker node")
	ctxWithTimeout, cancelFunc := context.WithTimeout(ctx, 5*time.Second)
	_, err := node.GRPCConnection.SendTask(ctxWithTimeout, &proto.RecievedTask{
		ID:             task.ID.Hex(),
		TaskAttributes: task.TaskAttributes,
		DockerImage:    task.ImageUrl,
	})
	cancelFunc()

	metricAttributes := metric.WithAttributes(
		attribute.String("nodeID", node.NodeEntry.NodeID),
		attribute.String("distributor", string(b.name)),
	)

	if err != nil {
		logger.Warnf("Fail to distribute task to worker node (%s)", err.Error())
		task.DistributionFailure(node.NodeEntry.NodeID, err)
		*errorTask = append(*errorTask, models.DistributeError{
			NodeEntry: node.NodeEntry,
			Task:      task,
			Error:     err,
		})
		b.failureMetricCounter.Add(
			ctx,
			1,
			metricAttributes,
		)
		return
	}
	// Add log to success task
	task.DistributionSuccess(node.NodeEntry.NodeID)
	b.metricCounter.Add(
		ctx,
		1,
		metricAttributes,
	)
	*successTask = append(*successTask, task)
	logger.Info("Worker Node accepted the task")
}
