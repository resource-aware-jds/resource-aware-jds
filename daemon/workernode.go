package daemon

import (
	"context"
	"fmt"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
)

type workerNode struct {
	ctx                    context.Context
	cancelFunc             func()
	controlPlaneGRPCClient proto.ControlPlaneClient
	workerNodeCertificate  cert.TLSCertificate
}

type WorkerNode interface {
	Start()
}

func ProvideWorkerNodeDaemon(controlPlaneGRPCClient proto.ControlPlaneClient, workerNodeCertificate cert.ClientCATLSCertificate) WorkerNode {
	ctx := context.Background()
	ctxWithCancel, cancelFunc := context.WithCancel(ctx)
	return &workerNode{
		ctx:                    ctxWithCancel,
		cancelFunc:             cancelFunc,
		controlPlaneGRPCClient: controlPlaneGRPCClient,
		workerNodeCertificate:  workerNodeCertificate,
	}
}

func (w *workerNode) Start() {
	err := w.checkInNodeToControlPlane()
	if err != nil {
		panic(fmt.Sprintf("Failed to check in worker node to control plane (%s)", err.Error()))
	}
}

func (w *workerNode) checkInNodeToControlPlane() error {
	certificate, err := w.workerNodeCertificate.GetCertificateInPEM()
	if err != nil {
		return err
	}

	_, err = w.controlPlaneGRPCClient.WorkerCheckIn(w.ctx, &proto.WorkerCheckInRequest{
		Certificate: certificate,
	})
	return err
}
