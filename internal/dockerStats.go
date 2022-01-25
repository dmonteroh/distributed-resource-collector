package internal

import (
	"encoding/json"

	"github.com/shirou/gopsutil/docker"
)

type DrcDockerStats struct {
	ContainerID string `json:"containerID"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	Status      string `json:"status"`
	Running     bool   `json:"running"`
}

func (d DrcDockerStats) String() string {
	s, _ := json.Marshal(d)
	return string(s)
}

func GetDockerStats() (dockerStats []DrcDockerStats) {
	tmpDocker, _ := docker.GetDockerStat()
	for _, docker := range tmpDocker {
		tmp := DrcDockerStats{
			ContainerID: docker.ContainerID,
			Name:        docker.Name,
			Image:       docker.Image,
			Status:      docker.Status,
			Running:     docker.Running,
		}
		dockerStats = append(dockerStats, tmp)
	}

	return dockerStats
}
