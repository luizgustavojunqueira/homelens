// Package shared
package shared

type SystemInfo struct {
	CPU         []CPU             `json:"cpu"`
	Memory      Memory            `json:"memory"`
	Disk        Disk              `json:"disk"`
	Network     []Network         `json:"network"`
	Temperature []Temperature     `json:"temperature,omitempty"`
	Containers  []DockerContainer `json:"containers,omitempty"`
	AgentIP     string            `json:"agent_ip"`
	Processes   []Process         `json:"processes"`
}

type CPU struct {
	Name         string  `json:"name"`
	UsagePercent float64 `json:"usage_percent"`
}

type Memory struct {
	Total     uint64 `json:"total"`
	Available uint64 `json:"available"`
	Used      uint64 `json:"used"`
}

type Disk struct {
	DiskSpace   DiskSpace     `json:"disk_space"`
	DiskIOUsage []DiskIOUsage `json:"disk_io_usage"`
}

type DiskSpace struct {
	Path         string  `json:"path"`
	Total        uint64  `json:"total"`
	Available    uint64  `json:"available"`
	Used         uint64  `json:"used"`
	UsagePercent float64 `json:"usage_percent"`
}

type DiskIOUsage struct {
	Name      string  `json:"name"`
	ReadMBps  float64 `json:"read_mbps"`
	WriteMBps float64 `json:"write_mbps"`
	IOPercent float64 `json:"io_percent"`
}

type Network struct {
	Name  string  `json:"name"`
	RxBps float64 `json:"rx_bps"`
	TxBps float64 `json:"tx_bps"`
}

type Temperature struct {
	Zone string  `json:"zone"`
	Temp float64 `json:"temp_c"`
}

type DockerPort struct {
	PrivatePort int    `json:"private_port"`
	PublicPort  int    `json:"public_port"`
	Type        string `json:"type"`
}

type DockerContainer struct {
	Name   string       `json:"name"`
	State  string       `json:"state"`
	Image  string       `json:"image"`
	Status string       `json:"status"`
	Ports  []DockerPort `json:"ports"`
}

type Process struct {
	PID     int     `json:"pid"`
	User    string  `json:"user"`
	CPU     float64 `json:"cpu"`
	Memory  float64 `json:"memory"`
	RSS     uint64  `json:"rss"`
	Uptime  int     `json:"uptime"`
	Name    string  `json:"name"`
	Cmdline string  `json:"command"`
}
