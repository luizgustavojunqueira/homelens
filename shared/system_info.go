// Package shared
package shared

type SystemInfo struct {
	CPUUsage    CPUUsage      `json:"cpu_usage"`
	Memory      MemoryUsage   `json:"memory"`
	DiskSpace   DiskSpace     `json:"disk_space"`
	DiskIOUsage []DiskIOUsage `json:"disk_io_usage"`
	NetUsage    []NetUsage    `json:"net_usage"`
	Temperature []TempInfo    `json:"temperature"`
}

type CPUUsage struct {
	CPUInfo []CPUInfo `json:"cpu_info"`
	CPUAvg  float64   `json:"cpu_avg"`
}

type CPUInfo struct {
	Name         string  `json:"name"`
	UsagePercent float64 `json:"usage_percent"`
}

type MemoryUsage struct {
	Total     uint64 `json:"total"`
	Available uint64 `json:"available"`
	Used      uint64 `json:"used"`
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

type NetUsage struct {
	Name  string  `json:"name"`
	RxBps float64 `json:"rx_bps"`
	TxBps float64 `json:"tx_bps"`
}

type TempInfo struct {
	Zone string  `json:"zone"`
	Temp float64 `json:"temp_c"`
}
