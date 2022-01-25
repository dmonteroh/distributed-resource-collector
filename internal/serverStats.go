package internal

import (
	"encoding/json"
	"time"
)

type DcrStats struct {
	Timestamp  DcrTimestamp     `json:"timestamp"`
	CPUStats   DrcCPUStats      `json:"cpuStats"`
	MemStats   DrcMemStats      `json:"memStats"`
	DiskStats  []DrcDiskStats   `json:"diskStats"`
	ProcStats  DrcProcStats     `json:"procStats"`
	DockerSats []DrcDockerStats `json:"dockerStats"`
}

type DcrTimestamp struct {
	TimeLocal   time.Time `json:"timeLocal"`
	TimeSeconds int64     `json:"timeSeconds"`
	TimeNano    int64     `json:"timeNano"`
}

func (d DcrStats) String() string {
	s, _ := json.Marshal(d)
	return string(s)
}

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
