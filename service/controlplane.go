package service

import (
	"context"
	"errors"
	"github.com/resource-aware-jds/resource-aware-jds/repository"
)

var (
	ErrNodeAlreadyRegistered = errors.New("there is the other node registered with the inputted key")
)

type ControlPlane struct {
	controlPlaneRepository repository.IControlPlane
}

type IControlPlane interface {
	RegisterWorker(ctx context.Context, ip string, port int, nodePublicKey string) (nodeID string, err error)
}

func ProvideControlPlane(controlPlaneRepository repository.IControlPlane) IControlPlane {
	return &ControlPlane{
		controlPlaneRepository: controlPlaneRepository,
	}
}

func (s *ControlPlane) RegisterWorker(ctx context.Context, ip string, port int, nodePublicKey string) (nodeID string, err error) {
	isExists, err := s.controlPlaneRepository.IsNodeAlreadyRegistered(ctx, nodePublicKey)
	if err != nil {
		return "", err
	}
	if isExists {
		return "", ErrNodeAlreadyRegistered
	}

	return "", nil
}
