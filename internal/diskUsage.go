package internal

import (
	"encoding/json"

	"github.com/shirou/gopsutil/disk"
)

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

func GetDiskUsage() []DrcDiskStats {
	parts, err := disk.Partitions(false)
	CheckError(err)

	var drcUsage []DrcDiskStats

	for _, part := range parts {
		u, err := disk.Usage(part.Mountpoint)
		CheckError(err)

		if !CustomContains(u.Path, "/snap/", "/etc/") {
			tmpUsage := DrcDiskStats{
				Device: part.Device,
				//SerialNumber: disk.GetDiskSerialNumber(part.Device),
				Path:        u.Path,
				Label:       disk.GetLabel(part.Device),
				Fstype:      part.Fstype,
				Total:       u.Total,
				Used:        u.Used,
				UsedPercent: u.UsedPercent,
			}
			//fmt.Println(tmpUsage)
			drcUsage = append(drcUsage, tmpUsage)
		}
	}
	return drcUsage
}
