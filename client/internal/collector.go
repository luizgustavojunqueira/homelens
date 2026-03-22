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
	NetUsage    []NetUsage    `json:"net_usage"`
	Temperature []TempInfo    `json:"temperature"`
}

func Collect(ctx context.Context, interval time.Duration, out chan<- SystemInfo) (SystemInfo, error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var prevCPUTime []CPUTime
	var prevDiskIO []DiskIO
	var prevNetInfo []NetInfo

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

			currentNetInfo, err := readNetInfo()
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

			if prevNetInfo == nil {
				prevNetInfo = currentNetInfo
				continue
			}

			sysInfo := SystemInfo{}

			sysInfo.CPUUsage = getCPUUsage(prevCPUTime, currentCPUTime)
			sysInfo.DiskIOUsage = calcDiskIOUsage(prevDiskIO, currentDiskIO, interval)
			sysInfo.NetUsage = calcNetUsage(prevNetInfo, currentNetInfo, interval)

			sysInfo.Memory, err = readMemoryUsage()
			if err != nil {
				return SystemInfo{}, err
			}

			sysInfo.DiskSpace, err = readDiskSpace("/")
			if err != nil {
				return SystemInfo{}, err
			}

			sysInfo.Temperature, err = readTempInfo()
			if err != nil {
				return SystemInfo{}, err
			}

			prevCPUTime = currentCPUTime
			prevDiskIO = currentDiskIO
			prevNetInfo = currentNetInfo

			out <- sysInfo
		}
	}
}
