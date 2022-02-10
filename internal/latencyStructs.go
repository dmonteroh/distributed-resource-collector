package internal

import "encoding/json"

// -- Latency Targets
type LatencyTargets struct {
	Hostname string          `json:"hostname"`
	Targets  []LatencyTarget `json:"targets"`
}

type LatencyTarget struct {
	Hostname     string `json:"hostname"`
	Hostport     string `json:"hostPort"`
	HostUser     string `json:"hostUser"`
	HostPassword string `json:"hostPassword"`
}

func LatencyJsonToStrcut(v string) (targets LatencyTargets, err error) {
	err = json.Unmarshal([]byte(v), &targets)
	return targets, err
}

// Latency Results
type LatencyResults struct {
	Hostname string          `json:"hostname"`
	Results  []LatencyResult `json:"results"`
}

func (r LatencyResults) String() string {
	s, _ := json.Marshal(r)
	return string(s)
}

type LatencyResult struct {
	Hostname string `json:"hostname"`
	Latency  int64  `json:"latency"`
}

func (r LatencyResult) String() string {
	s, _ := json.Marshal(r)
	return string(s)
}
