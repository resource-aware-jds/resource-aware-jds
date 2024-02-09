package util

import (
	"fmt"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"strings"
)

func ConvertToGb(memory models.MemoryWithUnit) models.MemoryWithUnit {
	// TODO change unit to enum
	lowerCaseUnit := strings.ToLower(memory.Unit)
	if lowerCaseUnit == "mb" || lowerCaseUnit == "mib" {
		return models.MemoryWithUnit{
			Size: memory.Size / 1024,
			Unit: "GiB",
		}
	}
	return memory
}

func SumInGb(first models.MemoryWithUnit, second models.MemoryWithUnit) models.MemoryWithUnit {
	firstGb := ConvertToGb(first)
	secondGb := ConvertToGb(second)
	fmt.Println(firstGb)
	fmt.Println(secondGb)
	return models.MemoryWithUnit{
		Size: firstGb.Size + secondGb.Size,
		Unit: "GiB",
	}
}
