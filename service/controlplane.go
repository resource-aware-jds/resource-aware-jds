package service

import "github.com/resource-aware-jds/resource-aware-jds/repository"

type ControlPlane struct {
	controlPlaneRepository repository.IControlPlane
}

type IControlPlane interface {
}

func ProvideControlPlane(controlPlaneRepository repository.IControlPlane) IControlPlane {
	return &ControlPlane{
		controlPlaneRepository: controlPlaneRepository,
	}
}
