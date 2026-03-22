package internal

import (
	"context"
	"time"
)

type SystemInfo struct {
	CPUInfo   CPUInfo     `json:"cpu_info"`
	Memory    MemoryUsage `json:"memory"`
	DiskSpace DiskSpace   `json:"disk_space"`
}

func Collect(ctx context.Context, interval time.Duration, out chan<- SystemInfo) (SystemInfo, error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var prevCPUTime []CPUTime

	for {
		select {
		case <-ctx.Done():
			return SystemInfo{}, ctx.Err()

		case <-ticker.C:
			currentCPUTime, err := readCPUTime()
			if err != nil {
				return SystemInfo{}, err
			}

			if prevCPUTime == nil {
				prevCPUTime = currentCPUTime
				continue
			}

			sysInfo := SystemInfo{}

			sysInfo.CPUInfo.Name = make([]string, len(currentCPUTime))
			for i, cpu := range currentCPUTime {
				sysInfo.CPUInfo.Name[i] = cpu.Name
			}
			sysInfo.CPUInfo.UsagePercent = getCPUUsage(prevCPUTime, currentCPUTime)
			sysInfo.CPUInfo.CPUAvg = getCPUAvg(sysInfo.CPUInfo.UsagePercent)
			prevCPUTime = currentCPUTime

			sysInfo.Memory, err = readMemoryUsage()
			if err != nil {
				return SystemInfo{}, err
			}

			sysInfo.DiskSpace, err = readDiskSpace("/")
			if err != nil {
				return SystemInfo{}, err
			}

			out <- sysInfo
		}
	}
}
