package internal

import (
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

func GetCPUUsage() (CPUStats DrcCPUStats) {
	tmpCPU, _ := cpu.Percent(time.Second/5, true)
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

	return DrcCPUStats{
		ModelName:    model,
		VendorID:     vendor,
		AverageUsage: totalPercent / float64(len(tmpCPU)),
		CoreUsage:    tmpCPU,
	}
}
