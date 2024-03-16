package distribution

import (
	"context"
	"fmt"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type baseDistributor struct {
	logger               *logrus.Entry
	metricCounter        metric.Int64Counter
	failureMetricCounter metric.Int64Counter
}

func newBaseDistributor(name DistributorName, meter metric.Meter) baseDistributor {
	counter, _ := meter.Int64Counter(
		fmt.Sprintf("rajds_cp_%s_distributor_total_distribute_task", name),
		metric.WithUnit("Task"),
		metric.WithDescription(fmt.Sprintf("The total task(s) that has been distributed using %s distributor", name)),
	)

	failureCounter, _ := meter.Int64Counter(
		fmt.Sprintf("rajds_cp_%s_distributor_total_distribute_task", name),
		metric.WithUnit("Task"),
		metric.WithDescription(fmt.Sprintf("The total task(s) that failed to be distributed using %s distributor", name)),
	)
	logger := logrus.WithFields(logrus.Fields{
		"role":             "Control Plane",
		"component":        "distributor",
		"distributor_name": name,
	})

	return baseDistributor{
		metricCounter:        counter,
		logger:               logger,
		failureMetricCounter: failureCounter,
	}
}

func (b *baseDistributor) distributeToNode(ctx context.Context, node NodeMapper, task models.Task, successTask *[]models.Task, errorTask *[]DistributeError) {
	logger := b.logger.WithFields(node.Logger.Data).WithField("taskID", task.ID.Hex())
	logger.Info("Sending task to the worker node")
	_, err := node.GRPCConnection.SendTask(ctx, &proto.RecievedTask{
		ID:             task.ID.Hex(),
		TaskAttributes: task.TaskAttributes,
		DockerImage:    task.ImageUrl,
	})

	if err != nil {
		logger.Warnf("Fail to distribute task to worker node (%s)", err.Error())
		task.DistributionFailure(node.NodeEntry.NodeID, err)
		*errorTask = append(*errorTask, DistributeError{
			NodeEntry: node.NodeEntry,
			Task:      task,
			Error:     err,
		})
		b.failureMetricCounter.Add(
			ctx,
			1,
			metric.WithAttributes(attribute.String("nodeID", node.NodeEntry.NodeID)),
		)
		return
	}
	// Add log to success task
	task.DistributionSuccess(node.NodeEntry.NodeID)
	b.metricCounter.Add(
		ctx,
		1,
		metric.WithAttributes(attribute.String("nodeID", node.NodeEntry.NodeID)),
	)
	*successTask = append(*successTask, task)
	logger.Info("Worker Node accepted the task")
}
