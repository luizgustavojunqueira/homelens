package internal

import (
	"context"
	"time"
)

type SystemInfo struct {
	CPUUsage    CPUUsage      `json:"cpu_usage"`
	Memory      MemoryUsage   `json:"memory"`
	DiskSpace   DiskSpace     `json:"disk_space"`
	DiskIOUsage []DiskIOUsage `json:"disk_io_usage"`
}

func Collect(ctx context.Context, interval time.Duration, out chan<- SystemInfo) (SystemInfo, error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var prevCPUTime []CPUTime
	var prevDiskIO []DiskIO

	for {
		select {
		case <-ctx.Done():
			return SystemInfo{}, ctx.Err()

		case <-ticker.C:
			currentCPUTime, err := readCPUTime()
			if err != nil {
				return SystemInfo{}, err
			}

			currentDiskIO, err := readDiskIO()
			if err != nil {
				return SystemInfo{}, err
			}

			if prevCPUTime == nil {
				prevCPUTime = currentCPUTime
				continue
			}

			if prevDiskIO == nil {
				prevDiskIO = currentDiskIO
				continue
			}

			sysInfo := SystemInfo{}

			sysInfo.CPUUsage = getCPUUsage(prevCPUTime, currentCPUTime)
			sysInfo.DiskIOUsage = calcDiskIOUsage(prevDiskIO, currentDiskIO, interval)

			sysInfo.Memory, err = readMemoryUsage()
			if err != nil {
				return SystemInfo{}, err
			}

			sysInfo.DiskSpace, err = readDiskSpace("/")
			if err != nil {
				return SystemInfo{}, err
			}

			prevCPUTime = currentCPUTime
			prevDiskIO = currentDiskIO

			out <- sysInfo
		}
	}
}
