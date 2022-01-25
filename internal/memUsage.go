package internal

import (
	"encoding/json"

	"github.com/shirou/gopsutil/mem"
)

type DrcMemStats struct {
	Total     uint64 `json:"total"`
	Available uint64 `json:"available"`
	Used      uint64 `json:"used"`
}

func (d DrcMemStats) String() string {
	s, _ := json.Marshal(d)
	return string(s)
}

func GetMemoryUsage() (MemStats DrcMemStats) {
	tmpMem, _ := mem.VirtualMemory()
	return DrcMemStats{
		Total:     tmpMem.Total,
		Available: tmpMem.Available,
		Used:      tmpMem.Used,
	}
}
