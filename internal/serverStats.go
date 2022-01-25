package internal

import "encoding/json"

type DcrStats struct {
	CPUStats   DrcCPUStats      `json:"cpuStats"`
	MemStats   DrcMemStats      `json:"memStats"`
	DiskStats  []DrcDiskStats   `json:"diskStats"`
	ProcStats  DrcProcStats     `json:"procStats"`
	DockerSats []DrcDockerStats `json:"dockerStats"`
}

func (d DcrStats) String() string {
	s, _ := json.Marshal(d)
	return string(s)
}

func GetServerStats() (dcrStats DcrStats) {

	return DcrStats{
		CPUStats:   GetCPUUsage(),
		MemStats:   GetMemoryUsage(),
		DiskStats:  GetDiskUsage(),
		ProcStats:  GetProcStats(),
		DockerSats: GetDockerStats(),
	}
}
