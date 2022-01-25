package main

import (
	"fmt"

	"github.com/dmonteroh/distributed-resource-collector/internal"
)

func main() {
	fmt.Println("Disk Usage:")
	diskUsage := internal.GetDiskUsage()
	for _, u := range diskUsage {
		fmt.Println(u.String())
	}

	fmt.Println("CPU Stats:")
	cpuStats := internal.GetCPUUsage()
	fmt.Println(cpuStats)

	fmt.Println("Mem Stats:")
	memStats := internal.GetMemoryUsage()
	fmt.Println(memStats)
}
