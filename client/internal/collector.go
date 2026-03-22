package internal

import (
	"context"
	"time"
)

type CPUInfo struct {
	Name         []string
	UsagePercent []float64
	CPUAvg       float64
}

type SystemInfo struct {
	CPUInfo   CPUInfo
	Memory    MemoryUsage
	DiskSpace DiskSpace
}

func Collect(ctx context.Context, interval time.Duration, out chan<- SystemInfo) (SystemInfo, error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var prevCPUInfo []CPUTime

	for {
		select {
		case <-ctx.Done():
			return SystemInfo{}, ctx.Err()

		case <-ticker.C:
			currentCPUInfo, err := readCPUInfo()
			if err != nil {
				return SystemInfo{}, err
			}

			if prevCPUInfo == nil {
				prevCPUInfo = currentCPUInfo
				continue
			}

			sysInfo := SystemInfo{}

			sysInfo.CPUInfo.Name = make([]string, len(currentCPUInfo))
			for i, cpu := range currentCPUInfo {
				sysInfo.CPUInfo.Name[i] = cpu.Name
			}
			sysInfo.CPUInfo.UsagePercent = getCPUUsage(prevCPUInfo, currentCPUInfo)
			sysInfo.CPUInfo.CPUAvg = getCPUAvg(sysInfo.CPUInfo.UsagePercent)
			prevCPUInfo = currentCPUInfo

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
