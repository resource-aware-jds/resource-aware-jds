package util

import (
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"strings"
)

func ExtractCpuUsage(containerUsage models.ContainerResourceUsage) (float64, error) {
	trimmedStr := strings.TrimSuffix(containerUsage.CpuUsage, "%")
	percentageFloat, err := strconv.ParseFloat(trimmedStr, 64)

	if err != nil {
		logrus.Errorf("There was an error converting the string to a float:  %v\n", err)
		// TODO add error handler
		return 0, err
	}
	return percentageFloat, nil
}

func ExtractMemoryUsage(input string) models.MemorySize {
	regex := regexp.MustCompile(`(\d+(\.\d+)?)([a-zA-Z]+)`)
	match := regex.FindStringSubmatch(input)

	if match != nil {
		number, _ := strconv.ParseFloat(match[1], 64)
		unit := match[3]

		result := models.MemorySize{
			Size: number,
			Unit: unit,
		}

		return result
	}

	return models.MemorySize{}
}
