package internal

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

type drcCPUStat struct {
	ModelName    string    `json:"modelName"`
	VendorID     string    `json:"vendorId"`
	AverageUsage float64   `json:"averageUsage"`
	CoreUsage    []float64 `json:"coreUsage"`
}

func (d drcCPUStat) String() string {
	s, _ := json.Marshal(d)
	return string(s)
}

func GetCPUUsage() (CPUStats drcCPUStat) {
	tmpCPU, _ := cpu.Percent(time.Second/10, true)
	totalPercent := 0.0
	for _, percent := range tmpCPU {
		totalPercent += percent
	}

	vendorList := []string{}
	modelList := []string{}
	cpuInfo, _ := cpu.Info()
	for _, cpu := range cpuInfo {
		vendorList = append(vendorList, cpu.VendorID)
		modelList = append(modelList, cpu.ModelName)
	}

	modelList = UniqueString(modelList)
	model := ""
	if len(modelList) > 1 {
		model = strings.Join(modelList, " / ")
	} else {
		model = modelList[0]
	}
	vendorList = UniqueString(vendorList)
	vendor := ""
	if len(vendorList) > 1 {
		vendor = strings.Join(vendorList, " / ")
	} else {
		vendor = vendorList[0]
	}

	return drcCPUStat{
		ModelName:    model,
		VendorID:     vendor,
		AverageUsage: totalPercent / float64(len(tmpCPU)),
		CoreUsage:    tmpCPU,
	}
}
