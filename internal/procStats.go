package internal

import (
	"encoding/json"

	"github.com/shirou/gopsutil/load"
)

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

func GetProcStats() (procStats DrcProcStats) {
	tmpProcs, _ := load.Misc()

	return DrcProcStats{
		TotalProcs:   tmpProcs.ProcsTotal,
		CreatedProcs: tmpProcs.ProcsCreated,
		RunningProcs: tmpProcs.ProcsRunning,
		BlockedProcs: tmpProcs.ProcsBlocked,
	}
}
