package internal

import (
	"github.com/shirou/gopsutil/mem"
)

type dcrMemStat struct {
	Total     uint64 `json:"total"`
	Available uint64 `json:"available"`
	Used      uint64 `json:"used"`
}

func GetMemoryUsage() (MemStats dcrMemStat) {
	tmpMem, _ := mem.VirtualMemory()
	return dcrMemStat{
		Total:     tmpMem.Total,
		Available: tmpMem.Available,
		Used:      tmpMem.Used,
	}
}
