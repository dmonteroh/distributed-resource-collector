package internal

import (
	"strings"

	"github.com/shirou/gopsutil/disk"
)

type drcDiskStat struct {
	Path        string  `json:"path"`
	Fstype      string  `json:"fstype"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}

func GetDiskUsage() []*drcDiskStat {
	parts, err := disk.Partitions(false)
	check(err)

	var drcUsage []*drcDiskStat

	for _, part := range parts {
		u, err := disk.Usage(part.Mountpoint)
		check(err)

		if !strings.Contains(u.Path, "/snap/") {
			tmpUsage := drcDiskStat{
				Path:        u.Path,
				Fstype:      u.Fstype,
				Total:       u.Total,
				Used:        u.Used,
				UsedPercent: u.UsedPercent,
			}
			//fmt.Println(tmpUsage)
			drcUsage = append(drcUsage, &tmpUsage)
		}
	}
	return drcUsage
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
