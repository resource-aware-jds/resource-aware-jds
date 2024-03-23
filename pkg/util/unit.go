package util

import (
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"strings"
)

func ConvertToGb(memory models.MemorySize) models.MemorySize {
	// TODO change unit to enum
	lowerCaseUnit := strings.ToLower(memory.Unit)
	if lowerCaseUnit == "mb" || lowerCaseUnit == "mib" {
		return models.MemorySize{
			Size: memory.Size / 1024,
			Unit: "GiB",
		}
	}
	return memory
}

func ConvertToMib(memory models.MemorySize) models.MemorySize {
	// TODO change unit to enum
	lowerCaseUnit := strings.ToLower(memory.Unit)
	if lowerCaseUnit == "gb" || lowerCaseUnit == "gib" {
		return models.MemorySize{
			Size: memory.Size * 1024,
			Unit: "MiB",
		}
	}
	return memory
}

func SumInGb(first models.MemorySize, second models.MemorySize) models.MemorySize {
	firstGb := ConvertToGb(first)
	secondGb := ConvertToGb(second)
	return models.MemorySize{
		Size: firstGb.Size + secondGb.Size,
		Unit: "GiB",
	}
}

func SubtractInGb(first models.MemorySize, second models.MemorySize) models.MemorySize {
	firstGb := ConvertToGb(first)
	secondGb := ConvertToGb(second)
	return models.MemorySize{
		Size: firstGb.Size - secondGb.Size,
		Unit: "GiB",
	}
}

func DivideBy(first models.MemorySize, value float64) models.MemorySize {
	return models.MemorySize{
		Size: first.Size / value,
		Unit: first.Unit,
	}
}
