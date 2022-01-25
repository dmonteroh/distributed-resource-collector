package internal

import (
	"encoding/json"
	"strings"

	"github.com/shirou/gopsutil/disk"
)

type drcDiskStat struct {
	Device       string  `json:"device"`
	SerialNumber string  `json:"serialNumber"`
	Path         string  `json:"path"`
	Label        string  `json:"label"`
	Fstype       string  `json:"fstype"`
	Total        uint64  `json:"total"`
	Used         uint64  `json:"used"`
	UsedPercent  float64 `json:"usedPercent"`
}

func (d drcDiskStat) String() string {
	s, _ := json.Marshal(d)
	return string(s)
}

func GetDiskUsage() []*drcDiskStat {
	parts, err := disk.Partitions(false)
	CheckError(err)

	var drcUsage []*drcDiskStat

	for _, part := range parts {
		u, err := disk.Usage(part.Mountpoint)
		CheckError(err)

		if !strings.Contains(u.Path, "/snap/") {
			tmpUsage := drcDiskStat{
				Device:       part.Device,
				SerialNumber: disk.GetDiskSerialNumber(part.Device),
				Path:         u.Path,
				Label:        disk.GetLabel(part.Device),
				Fstype:       part.Fstype,
				Total:        u.Total,
				Used:         u.Used,
				UsedPercent:  u.UsedPercent,
			}
			//fmt.Println(tmpUsage)
			drcUsage = append(drcUsage, &tmpUsage)
		}
	}
	return drcUsage
}
