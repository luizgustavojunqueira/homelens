package internal

import (
	"context"
	"time"
)

type SystemInfo struct {
	CPUUsages []float64
	CPUAvg    float64
}

func Collect(ctx context.Context, interval time.Duration, out chan<- SystemInfo) (SystemInfo, error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var prevCPUInfo []CPUInfo

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
			sysInfo.CPUUsages = getCPUUsage(prevCPUInfo, currentCPUInfo)
			sysInfo.CPUAvg = getCPUAvg(sysInfo.CPUUsages)
			prevCPUInfo = currentCPUInfo

			out <- sysInfo
		}
	}
}
