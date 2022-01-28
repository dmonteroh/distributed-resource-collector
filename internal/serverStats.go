package internal

import (
	"time"
)

func GetServerStats() (dcrStats DcrStats) {
	tmpTime := time.Now()
	timestamp := DcrTimestamp{
		TimeLocal:   tmpTime,
		TimeSeconds: tmpTime.Unix(),
		TimeNano:    tmpTime.UnixNano(),
	}

	return DcrStats{
		Timestamp:  timestamp,
		CPUStats:   GetCPUUsage(),
		MemStats:   GetMemoryUsage(),
		DiskStats:  GetDiskUsage(),
		ProcStats:  GetProcStats(),
		DockerSats: GetDockerStats(),
	}
}
