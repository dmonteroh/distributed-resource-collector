package internal

import (
	"encoding/json"
	"time"
)

// -- CPU
type DrcCPUStats struct {
	ModelName    string    `json:"modelName"`
	VendorID     string    `json:"vendorId"`
	AverageUsage float64   `json:"averageUsage"`
	CoreUsage    []float64 `json:"coreUsage"`
}

func (d DrcCPUStats) String() string {
	s, _ := json.Marshal(d)
	return string(s)
}

// -- DISK
type DrcDiskStats struct {
	Device string `json:"device"`
	//SerialNumber string  `json:"serialNumber"`
	Path        string  `json:"path"`
	Label       string  `json:"label"`
	Fstype      string  `json:"fstype"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}

func (d DrcDiskStats) String() string {
	s, _ := json.Marshal(d)
	return string(s)
}

// -- MEMORY / RAM
type DrcMemStats struct {
	Total     uint64 `json:"total"`
	Available uint64 `json:"available"`
	Used      uint64 `json:"used"`
}

func (d DrcMemStats) String() string {
	s, _ := json.Marshal(d)
	return string(s)
}

// -- PROCESSES
type DrcProcStats struct {
	TotalProcs   int `json:"totalProcs"`
	CreatedProcs int `json:"createdProcs"`
	RunningProcs int `json:"runningProcs"`
	BlockedProcs int `json:"blockedProcs"`
}

func (d DrcProcStats) String() string {
	s, _ := json.Marshal(d)
	return string(s)
}

// -- DOCKER
type DrcDockerStats struct {
	ContainerID string `json:"containerID"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	Status      string `json:"status"`
	State       string `json:"State"`
}

func (d DrcDockerStats) String() string {
	s, _ := json.Marshal(d)
	return string(s)
}

// -- TIMESTAMP
type DrcTimestamp struct {
	TimeLocal   time.Time `json:"timeLocal"`
	TimeSeconds int64     `json:"timeSeconds"`
	TimeNano    int64     `json:"timeNano"`
}

func (d DrcTimestamp) String() string {
	s, _ := json.Marshal(d)
	return string(s)
}

// -- HOST INFO
type DrcHost struct {
	Hostname             string `json:"hostname"`
	Uptime               int64  `json:"uptime"`
	BootTime             int64  `json:"boottime"`
	Platform             string `json:"platform"`
	VirtualizationSystem string `json:"virtualizationSystem"`
	VirtualizationRole   string `json:"virtualizationRole"`
	HostID               string `json:"hostid"`
}

func (d DrcHost) String() string {
	s, _ := json.Marshal(d)
	return string(s)
}

// -- RESPONSE OBJECT
type DrcStats struct {
	Timestamp  DrcTimestamp     `json:"timestamp"`
	DrcHost    DrcHost          `json:"host"`
	CPUStats   DrcCPUStats      `json:"cpuStats"`
	MemStats   DrcMemStats      `json:"memStats"`
	DiskStats  []DrcDiskStats   `json:"diskStats"`
	ProcStats  DrcProcStats     `json:"procStats"`
	DockerSats []DrcDockerStats `json:"dockerStats"`
}

func (d DrcStats) String() string {
	s, _ := json.Marshal(d)
	return string(s)
}

func (d DrcStats) Marshal() []byte {
	s, err := json.Marshal(d)
	if err != nil {
		return s
	} else {
		return nil
	}
}
