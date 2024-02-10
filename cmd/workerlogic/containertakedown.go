package workerlogic

type ContainerTakeDownState struct {
}

type ContainerTakeDown interface {
	Calculate(state ContainerTakeDownState) []string
}

func ProvideOverResourceUsageContainerTakeDown() ContainerTakeDown {
	return &OverResourceUsageContainerTakeDown{}
}

type OverResourceUsageContainerTakeDown struct{}

func (o OverResourceUsageContainerTakeDown) Calculate(state ContainerTakeDownState) []string {
	//TODO implement me
	panic("implement me")
}
